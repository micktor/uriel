package httpd

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/your_org/uriel/internal/config"
	"github.com/your_org/uriel/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	config  config.Config
	service service.Service
	oauth   *oauth2.Config
}

// NewHandler returns a configured Handler.
func NewHandler(config *config.Config, service *service.Service) Handler {
	return Handler{
		config:  *config,
		service: *service,
		oauth: &oauth2.Config{
			ClientID:     config.OAuth.ClientID,
			ClientSecret: config.OAuth.ClientSecret,
			RedirectURL:  config.HTTPServer.Hostname + "/public/auth/callback",
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *Handler) Run() {
	mux := http.NewServeMux()

	handler := useMiddlewares(mux, loggingMiddleware, realIPMiddleware, corsMiddleware)
	setupRoutes(mux, h)

	slog.Info("server running", slog.String("port", h.config.HTTPServer.Port))
	err := http.ListenAndServe(":"+h.config.HTTPServer.Port, handler)
	if err != nil {
		slog.Error("server failed to start", slog.String("error", err.Error()))
	}
}

func setupRoutes(mux *http.ServeMux, a *Handler) {
	mux.HandleFunc("POST /public/auth/register", a.Register)
	mux.HandleFunc("POST /public/auth/login", a.LoginWithCookie)
	//mux.HandleFunc("/public/auth/logout", a.handler.LoginWithCookie)

	// OAuth
	mux.HandleFunc("GET /public/auth/oauth", a.oAuthHandler)
	mux.HandleFunc("GET /public/auth/callback", a.oAuthCallbackHandler)

	privateMux := http.NewServeMux()
	privateHandler := useMiddlewares(privateMux, jwtauth.Verifier(a.service.GetTokenAuth()), jwtAuthenticator, recordUserMiddleware)
	privateMux.HandleFunc("GET /private/viewer", getJSON(a.GetViewer))

	mux.Handle("/private/", privateHandler)
}

func useMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	logger := httplog.NewLogger("uriel", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "msg",
		Tags: map[string]string{
			"env": "prod",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	return httplog.RequestLogger(logger)(next)
}

// Real IP middleware
func realIPMiddleware(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})(next)
}

func recordUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		httplog.LogEntrySetField(ctx, "user", slog.StringValue(getJWT(ctx).Token.Subject()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authenticate requests using JWT tokens.
func jwtAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _ := r.Context().Value(jwtauth.TokenCtxKey).(jwt.Token)
		if token == nil {
			slog.WarnContext(r.Context(), "nil token encountered")
			http.Error(w, "nil jwt token", http.StatusUnauthorized)
			return
		} else if err := jwt.Validate(token); err != nil {
			slog.ErrorContext(r.Context(), "failed to validate token", slog.String("error", err.Error()), slog.String("user", token.Subject()))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), service.JWTAuthCtxKey, service.JWTAuth{Token: token})
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
