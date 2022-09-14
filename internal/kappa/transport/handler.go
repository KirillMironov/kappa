package transport

import (
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct {
	deployer deployer
}

type deployer interface {
	Deploy(domain.Deployment) error
	Cancel(domain.Deployment) error
}

func NewHandler(deployer deployer) *Handler {
	return &Handler{deployer: deployer}
}

func (h Handler) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Route("/api/v1/deploy", func(r chi.Router) {
		router.Post("/", h.deploy)
		router.Delete("/", h.cancel)
	})

	return router
}

func (h Handler) deploy(w http.ResponseWriter, r *http.Request) {
	deployment, err := httputil.StructFromBodyJSON[domain.Deployment](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.deployer.Deploy(deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) cancel(w http.ResponseWriter, r *http.Request) {
	deployment, err := httputil.StructFromBodyJSON[domain.Deployment](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.deployer.Cancel(deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
