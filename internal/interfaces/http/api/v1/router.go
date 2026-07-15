package v1

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
)

func rateLimit() func(http.Handler) http.Handler {
	limiter := httprate.NewRateLimiter(
		10,
		time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			return httprate.CanonicalizeIP(ip), nil
		}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
		}),
	)

	return limiter.Handler
}

func (a *API) NewRouter(jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(rateLimit())
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", a.Register)
		r.Post("/login", a.Login)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtSecret))

			// r.Get("/me", authHandler.Me)

			r.Route("/teams", func(r chi.Router) {
				r.Post("/", a.CreateTeam)
				r.Get("/", a.ListTeams)
				r.Post("/{teamID}/invite", a.InviteUser)

				// r.Group(func(r chi.Router) {
				// 	r.Use(RequireRoles("OWNER", "ADMIN"))
				// 	r.Post("/{teamID}/invite", a.InviteUser)
				// })
			})
			r.Route("/tasks", func(r chi.Router) {
				r.Post("/", a.CreateTask)
				r.Get("/", a.ListTasks)
				r.Put("/{taskID}", a.ListTasks)
				r.Get("/{taskID}/history", a.TaskHistory)
			})
		})
	})

	return r
}
