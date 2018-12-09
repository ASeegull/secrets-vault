package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	uuid "github.com/satori/go.uuid"

	"github.com/ASeegull/secrets-vault/cleanup"
	"github.com/ASeegull/secrets-vault/models"
)

type API struct {
	http.Server
	storage querier
	Cleaner *cleanup.Manager
}

type querier interface {
	Save(uuid.UUID, int, time.Time, string) error
	Get(uuid.UUID) (*models.Secret, error)
	DecrementViews(uuid.UUID) (int, error)
}

func New(db querier, host, port string) *API {
	api := &API{storage: db}
	api.Addr = host + ":" + port
	api.Handler = api.InitRoutes()
	return api
}

func (a *API) InitRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/secret", func(r chi.Router) {
			r.Post("/", a.addSecret)
			r.Get("/{hash}", a.getSecret)
		})
	})
	return r
}
