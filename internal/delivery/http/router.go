package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/azagrivin/testProjectCourses/config"
	er "github.com/azagrivin/testProjectCourses/internal/delivery/http/error"
	"github.com/azagrivin/testProjectCourses/internal/logger"
	"github.com/azagrivin/testProjectCourses/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Chi      *chi.Mux
	cfg      *config.Config
	services *services.Services
	log      logger.HttpLogger
}

func NewRouter(cfg *config.Config, services *services.Services, log logger.HttpLogger) *Router {
	r := &Router{
		Chi:      chi.NewRouter(),
		cfg:      cfg,
		services: services,
		log:      log,
	}

	r.addMiddleware()
	r.addHandlers()

	return r
}

func (r *Router) addMiddleware() {
	r.Chi.Use(middleware.Recoverer)
	r.Chi.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: stdLogFunc(r.log.Info)}))
}

func (r *Router) addHandlers() {
	r.Chi.Route("/api", func(router chi.Router) {

		router.Route("/btcusdt", func(router chi.Router) {
			router.Get("/", NewHandlerGetBtcusdt(r.services.BtcUsdt, r.log))
			router.Post("/", NewHandlerPostBtcusdt(r.services.BtcUsdt, r.log))
		})

		router.Route("/currencies", func(router chi.Router) {
			router.Get("/", NewHandlerGetCurrencies(r.services.BtcUsdt, r.log))
			router.Post("/", NewHandlerPostCurrencies(r.services.BtcUsdt, r.log))
		})

		router.Route("/latest", func(router chi.Router) {
			router.Get("/", NewHandlerGetBtc(r.services.BtcUsdt, r.log))
			router.Post("/", NewHandlerPostBtc(r.services.BtcUsdt, r.log))
		})
	})
}

func (r *Router) Run() {
	r.log.Infof("Start server: %s", r.cfg.App.Url)
	if err := http.ListenAndServe(":"+r.cfg.App.Port, r.Chi); err != nil {
		r.log.Errorf("listen and serve error, %v", err)
	}
	os.Exit(1)
}

func writeResponse(w http.ResponseWriter, resp interface{}, log logger.Logger) {
	if resp != nil {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Errorf("marshal response error, %s", err)
			er.ErrInternalError.Handle(w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
