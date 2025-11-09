package v1

import (
	"denet-test-task/internal/api/v1/apierrs"
	"denet-test-task/internal/services/tasks"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type tasksRoutes struct {
	tasksService tasks.Tasks
}

func newTasksRoutes(router chi.Router, tasksService tasks.Tasks) {
	routes := &tasksRoutes{
		tasksService: tasksService,
	}

	router.Get("/list", routes.handleGetTasks)
}

func (r *tasksRoutes) handleGetTasks(w http.ResponseWriter, req *http.Request) {

	tasks, err := r.tasksService.GetAllTasks(req.Context())
	if err != nil {
		apierrs.NewErrorResponseHTTP(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}
