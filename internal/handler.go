package internal

import (
	"app/internal/apperror"
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
	router.HandlerFunc(http.MethodGet, tasksURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodPost, tasksURL, apperror.Middleware(h.CreateTask))
	router.HandlerFunc(http.MethodGet, taskURL, apperror.Middleware(h.GetTaskByUUID))
	router.HandlerFunc(http.MethodPut, taskURL, apperror.Middleware(h.UpdateTask))
	router.HandlerFunc(http.MethodPatch, taskURL, apperror.Middleware(h.PartiallyUpdateTask))
	router.HandlerFunc(http.MethodDelete, taskURL, apperror.Middleware(h.DeleteTask))

}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	w.Write([]byte("this is list of task"))
	return nil
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) error {
	w.WriteHeader(201)
	w.Write([]byte("this is list of task"))
	return nil
}
func (h *handler) GetTaskByUUID(w http.ResponseWriter, r *http.Request, param httprouter.Params) error {
	w.WriteHeader(200)
	w.Write([]byte("this is list by id"))
	return nil
}
func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) error {
	w.WriteHeader(204)
	w.Write([]byte("this is list of task"))
	return nil
}
func (h *handler) PartiallyUpdateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) error {
	w.WriteHeader(204)
	w.Write([]byte("this is list of task"))
	return nil
}
func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) error {
	w.WriteHeader(204)
	w.Write([]byte("this is list of task"))
	return nil
}
