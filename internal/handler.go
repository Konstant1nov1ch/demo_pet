package internal

import (
	"app/internal/handlers"
	"app/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	tasksURL = "/task"
	taskURL  = "/task/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(tasksURL, h.GetList)
	router.POST(tasksURL, h.CreateTask)
	router.GET(taskURL, h.GetTaskByUUID)
	router.PUT(taskURL, h.UpdateTask)
	router.PATCH(taskURL, h.PartiallyUpdateTask)
	router.DELETE(taskURL, h.DeleteTask)

}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list of task"))
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list of task"))
}
func (h *handler) GetTaskByUUID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list by id"))
}
func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list of task"))
}
func (h *handler) PartiallyUpdateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list of task"))
}
func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte("this is list of task"))
}
