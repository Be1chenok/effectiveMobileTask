package handler

import (
	"net/http"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	"github.com/Be1chenok/effectiveMobileTask/internal/service"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type Handler struct {
	conf    *config.Config
	logger  appLogger.Logger
	service *service.Service
}

func New(conf *config.Config, logger appLogger.Logger, service *service.Service) *Handler {
	return &Handler{
		conf:    conf,
		logger:  logger.With(zap.String("component", "handler")),
		service: service,
	}
}

func (h Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/persons", h.FindPersons).Methods("GET")
	router.HandleFunc("/person/{id:[0-9]+}", h.FindPersonById).Methods("GET")
	router.HandleFunc("/person", h.AddPerson).Methods("POST")
	router.HandleFunc("/person/{id:[0-9]+}", h.DeletePerson).Methods("DELETE")
	router.HandleFunc("/person", h.UpdatePerson).Methods("PUT")

	loggedMux := h.LoggerMiddleware(router)

	return loggedMux
}
