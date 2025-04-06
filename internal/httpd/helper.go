package httpd

import (
	"context"
	"encoding/json"
	"github.com/your_org/uriel/internal/service"
	"net/http"
	"net/url"
)

type ctxKey string

const (
	requestKey ctxKey = "uriel:request"
)

type Request struct {
	request     *http.Request
	queryValues url.Values
}

func jsonHandler[TIN any, TOUT any](fn func(ctx context.Context, req TIN) (TOUT, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestKey, Request{
			request:     r,
			queryValues: r.URL.Query(),
		})

		var req TIN
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request", http.StatusUnprocessableEntity)
			return
		}

		resp, err := fn(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getJSON[T any](fn func(ctx context.Context) (T, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestKey, Request{
			request:     r,
			queryValues: r.URL.Query(),
		})

		resp, err := fn(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getParam(ctx context.Context, key string) string {
	return ctx.Value(requestKey).(Request).request.PathValue(key)
}

func getQuery(ctx context.Context, key string) string {
	return ctx.Value(requestKey).(Request).queryValues.Get(key)
}

func getJWT(ctx context.Context) service.JWTAuth {
	jwtAuth, _ := ctx.Value(service.JWTAuthCtxKey).(service.JWTAuth)
	return jwtAuth
}
