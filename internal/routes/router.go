package routes

import (
	// "net/http"
	"net/http/pprof"

	"github.com/NathanielRand/webchest-image-converter-api/internal/handlers"
	"github.com/NathanielRand/webchest-image-converter-api/internal/middleware"
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
	// chain = chain.Append(func(next http.Handler) http.Handler {
	// 	return middleware.AuthorizationMiddleware(next, "admin")
	// })
	// chain = chain.Append(middleware.RateLimitingMiddleware)
	// chain = chain.Append(middleware.QuotaMiddleware)
	// chain = chain.Append(middleware.CachingMiddleware)
	chain = chain.Append(middleware.LoggingMiddleware)

	// API endpoints to the router

	// General endpoints
	router.Handle("/api/v1/hello", chain.ThenFunc(handlers.HelloHandler)).Methods("GET")
	router.Handle("/api/v1/health", chain.ThenFunc(handlers.HealthHandler)).Methods("GET")

	// User endpoints
	router.Handle("/api/v1/image/convert", chain.ThenFunc(handlers.ImageConvertHandler)).Methods("POST")

	// router.Handle("/api/v1/resize/image/{height}/{width}", chain.ThenFunc(handlers.ImageResizeHandler)).Methods("POST")
	// router.Handle("/api/v1/crop/image/{height}/{width}/{x}/{y}", chain.ThenFunc(handlers.ImageCropHandler)).Methods("POST")
	// router.Handle("/api/v1/zoom/image/{factor}", chain.ThenFunc(handlers.ImageZoomHandler)).Methods("POST")
	// router.Handle("/api/v1/rotate/image/{degrees}", chain.ThenFunc(handlers.ImageRotateHandler)).Methods("POST")
	// router.Handle("/api/v1/flip/image/{direction}", chain.ThenFunc(handlers.ImageFlipHandler)).Methods("POST")
	// router.Handle("/api/v1/blur/image/{sigma}", chain.ThenFunc(handlers.ImageBlurHandler)).Methods("POST")
	// router.Handle("/api/v1/contrast/image/{factor}", chain.ThenFunc(handlers.ImageContrastHandler)).Methods("POST")
	// router.Handle("/api/v1/brightness/image/{factor}", chain.ThenFunc(handlers.ImageBrightnessHandler)).Methods("POST")
	// router.Handle("/api/v1/sharpen/image/{sigma}/{radius}", chain.ThenFunc(handlers.ImageSharpenHandler)).Methods("POST")
	// router.Handle("/api/v1/median/image/{radius}", chain.ThenFunc(handlers.ImageMedianHandler)).Methods("POST")
	// router.Handle("/api/v1/emboss/image/{radius}", chain.ThenFunc(handlers.ImageEmbossHandler)).Methods("POST")
	// router.Handle("/api/v1/edge/image/{radius}", chain.ThenFunc(handlers.ImageEdgeHandler)).Methods("POST")
	// router.Handle("/api/v1/normalize/image/{percent}", chain.ThenFunc(handlers.ImageNormalizeHandler)).Methods("POST")
	// router.Handle("/api/v1/grayscale/image/{percent}", chain.ThenFunc(handlers.ImageGrayscaleHandler)).Methods("POST")
	// router.Handle("/api/v1/sepia/image/{percent}", chain.ThenFunc(handlers.ImageSepiaHandler)).Methods("POST")
	// router.Handle("/api/v1/invert/image/", chain.ThenFunc(handlers.ImageInvertHandler)).Methods("POST")

	// Debug endpoints
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return router
}
