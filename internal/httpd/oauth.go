package httpd

import (
	"context"
	"encoding/json"
	"github.com/your_org/uriel/internal/dto"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type GoogleUserInfo struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// OAuth handler create an unique OAuth URL using the client ID and client secret.
// then it redirect the user to the OAuth provider website to complete the login.
func (h *Handler) oAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := h.oauth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuth callback handler handle the redirect request from the OAuth provider.
// It read the code query parameter and exchange it to get the access token.
// Then this handler call user info endpoint to get the user public detail eg.,
// Name, Email, Profile picture etc.
func (h *Handler) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// Exchanging the code for an access token
	t, err := h.oauth.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Creating an HTTP client to make authenticated request using the access key.
	// This client method also regenerate the access key using the refresh key.
	client := h.oauth.Client(context.Background(), t)

	// Getting the user public details from google API endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Closing the request body when this function returns.
	// This is h good practice to avoid memory leak
	defer resp.Body.Close()

	var v GoogleUserInfo

	// Reading the JSON body using JSON decoder
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, err := h.service.OAuthLoginUser(r.Context(), dto.RegisterRequest{
		Email:        v.Email,
		AuthProvider: "google",
		OAuthID:      v.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: tokenString,
		Path:  "/",
		//Domain:   "uriel.com",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, h.config.FEHost+"/dashboard", http.StatusSeeOther)
}
