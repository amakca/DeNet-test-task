package v1

import (
	"denet-test-task/internal/api/v1/apierrs"
	"denet-test-task/internal/services/contracts"
	"denet-test-task/internal/services/users"
	"denet-test-task/pkg/validator"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authRoutes struct {
	authService contracts.Auth
}

type authInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newAuthRoutes(router chi.Router, authService contracts.Auth) {
	routes := &authRoutes{
		authService: authService,
	}

	router.Post("/sign-up", routes.handleSignup)
	router.Post("/sign-in", routes.handleLogin)
}

func (a *authRoutes) handleSignup(w http.ResponseWriter, req *http.Request) {
	var input authInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.NewCustomValidator().Validate(input); err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := a.authService.CreateUser(req.Context(), contracts.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == users.ErrUserAlreadyExists {
			apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, err.Error())
			return
		}
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, "internal server error")
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response{Id: id})
}

func (a *authRoutes) handleLogin(w http.ResponseWriter, req *http.Request) {
	var input authInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.NewCustomValidator().Validate(input); err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := a.authService.GenerateToken(req.Context(), contracts.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == users.ErrUserNotFound {
			apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid username or password")
			return
		}
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, "internal server error")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response{Token: token})
}
