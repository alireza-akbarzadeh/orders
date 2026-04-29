package application

import (
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddlware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/techies/orders-api/handler"
	middleware "github.com/techies/orders-api/middleware"
	"github.com/techies/orders-api/repository/auth"
	Orders "github.com/techies/orders-api/repository/orders"
	"github.com/techies/orders-api/repository/users"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	// middlewares
	router.Use(chiMiddlware.Logger)
	router.Use(httprate.LimitByIP(a.Config.RateLimit, time.Minute))
	router.Use(chiMiddlware.Recoverer)

	// book routes
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes (no auth)
	router.Get("/", handler.ShowLandingPage)
	router.Route("/auth", a.loadAuthRoutes)

	// Protected routes (require authentication)
	router.Group(func(r chi.Router) {
		// Apply authentication middleware to this whole group
		r.Use(middleware.Authenticate([]byte(a.Config.JWTSecret)))

		// Sub-routes that require authentication
		r.Route("/users", a.loadUsersRoutes)
		r.Route("/orders", a.loadOrderRoutes)
	})

	a.Router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &Orders.Repo{DB: a.DB},
	}

	router.Get("/", orderHandler.ListOrders)

	router.With(middleware.Authenticate([]byte(a.Config.JWTSecret))).Post("/", orderHandler.CreateOrder)
	router.With(middleware.Authenticate([]byte(a.Config.JWTSecret))).Put("/", orderHandler.UpdateOrder)
	router.With(middleware.Authenticate([]byte(a.Config.JWTSecret))).Put("/{id}/cancel", orderHandler.CancelOrder)
	router.With(middleware.Authenticate([]byte(a.Config.JWTSecret))).Get("/my", orderHandler.GetUserOrders)
	router.With(middleware.RequireRole("admin")).Delete("/{id}", orderHandler.DeleteOrder)
}

func (a *App) loadUsersRoutes(router chi.Router) {
	userHandler := &handler.UsersHandler{
		Repo: &users.Repo{DB: a.DB},
	}

	router.With(middleware.RequireRole("admin")).Get("/", userHandler.GetUserList)
	router.Get("/{id}", userHandler.GetUserById) // user can view own profile? add ownership check inside handler
	router.With(middleware.RequireRole("admin")).Put("/{id}", userHandler.UpdateUser)
	router.With(middleware.RequireRole("admin")).Delete("/{id}", userHandler.DeleteUser)
}

func (a *App) loadAuthRoutes(router chi.Router) {
	authHandler := &handler.Auth{
		Repo:      &auth.Repo{DB: a.DB},
		JWTSecret: []byte(a.Config.JWTSecret),
	}

	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate([]byte(a.Config.JWTSecret)))
		r.Use(middleware.RequireRole("admin"))
		r.Post("/find-by-email", authHandler.FindByEmailHandler)
		r.Post("/find-by-username", authHandler.FindByUserNameHandler)
	})
}
