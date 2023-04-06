package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/NathanielRand/morebytes-templates/go/api/clean-layered/internal/handlers"
	"github.com/NathanielRand/morebytes-templates/go/api/clean-layered/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func SetupRouter() *mux.Router {
	// Initialize a new router from the Gorilla Mux library
	router := mux.NewRouter()

	// Create a new middleware chain using the Alice library
	chain := alice.New()

	// Add middleware to the chain for authentication, rate limiting, caching, and quotas
	chain = chain.Append(middleware.SecurityMiddleware)
	chain = chain.Append(middleware.AuthenticationMiddleware)
	chain = chain.Append(func(next http.Handler) http.Handler {
		return middleware.AuthorizationMiddleware(next, "admin")
	})
	chain = chain.Append(middleware.RateLimitingMiddleware)
	chain = chain.Append(middleware.QuotaMiddleware)
	chain = chain.Append(middleware.CachingMiddleware)
	chain = chain.Append(middleware.LoggingMiddleware)

	// API endpoints to the router

	// General endpoints
	router.Handle("/api/v1/hello", chain.ThenFunc(handlers.HelloHandler)).Methods("GET")
	router.Handle("/api/v1/health", chain.ThenFunc(handlers.HealthHandler)).Methods("GET")

	// User endpoints
	// router.Handle("/api/v1/convert/image/{from}/{to}", chain.ThenFunc(handlers.ConvertHandler)).Methods("POST")

	// Debug endpoints
	// Register pprof endpoints
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return router
}