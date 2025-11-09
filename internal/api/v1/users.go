package v1

import (
	"denet-test-task/internal/api/v1/apierrs"
	"denet-test-task/internal/services/tasks"
	"denet-test-task/internal/services/users"
	"denet-test-task/pkg/logctx"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type usersRoutes struct {
	usersService users.Users
	tasksService tasks.Tasks
}

func newUsersRoutes(router chi.Router, usersService users.Users, tasksService tasks.Tasks) {
	routes := &usersRoutes{
		usersService: usersService,
		tasksService: tasksService,
	}

	router.Get("/{user_id}/status", routes.handleGetUserStatus)
	router.Get("/{user_id}/history", routes.handleGetHistory)
	router.Get("/{user_id}/points", routes.handleGetPoints)
	router.Get("/leaderboard", routes.handleGetLeaderboard)

	router.Post("/{user_id}/referrer", routes.handleSetReferrer)
	router.Post("/{user_id}/email", routes.handleSetEmail)

	router.Post("/{user_id}/task/complete", routes.handleCompleteTask)
}

func (r *usersRoutes) handleGetUserStatus(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		logctx.FromContext(req.Context()).Error("usersRoutes.handleGetUserStatus - strconv.Atoi", "err", err)
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := r.usersService.GetInfo(req.Context(), users.UsersGetInfoInput{UserId: userIdInt})
	if err != nil {
		logctx.FromContext(req.Context()).Error("usersRoutes.handleGetUserStatus - usersService.GetInfo", "err", err)
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (r *usersRoutes) handleGetHistory(w http.ResponseWriter, req *http.Request) {

	limit := req.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 || limitInt > 100 {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid limit")
		return
	}

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	history, err := r.usersService.GetHistory(req.Context(), users.UsersGetHistoryInput{UserId: userIdInt, Limit: limitInt})
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(history)
}

func (r *usersRoutes) handleGetPoints(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	points, err := r.usersService.GetPoints(req.Context(), users.UsersGetPointsInput{UserId: userIdInt})
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(points)
}

func (r *usersRoutes) handleGetLeaderboard(w http.ResponseWriter, req *http.Request) {

	limit := req.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 || limitInt > 100 {
		logctx.FromContext(req.Context()).Error("usersRoutes.handleGetLeaderboard - strconv.Atoi", "err", err)
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid limit")
		return
	}

	leaderboard, err := r.usersService.GetLeaderboard(req.Context(), users.UsersGetLeaderboardInput{Limit: limitInt})
	if err != nil {
		logctx.FromContext(req.Context()).Error("usersRoutes.handleGetLeaderboard - usersService.GetLeaderboard", "err", err)
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(leaderboard)
}

func (r *usersRoutes) handleSetReferrer(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	referrer := req.FormValue("referrer")
	referrerInt, err := strconv.Atoi(referrer)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid referrer")
		return
	}

	err = r.usersService.SetReferrer(req.Context(), users.UsersSetReferrerInput{UserId: userIdInt, Referrer: referrerInt})
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

func (r *usersRoutes) handleSetEmail(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	email := req.FormValue("email")
	err = r.usersService.SetEmail(req.Context(), users.UsersSetEmailInput{UserId: userIdInt, Email: email})
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

func (r *usersRoutes) handleCompleteTask(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid user id")
		return
	}

	taskId := req.FormValue("task_id")
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusBadRequest, "invalid task id")
		return
	}

	err = r.usersService.CompleteTask(req.Context(), users.UsersCompleteTaskInput{UserId: userIdInt, TaskId: taskIdInt})
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}
