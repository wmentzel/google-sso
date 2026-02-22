package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"google-sso/auth"
)

// Login generates a random state token, stores it in the session, and
// redirects the browser to Google's consent screen.
func Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.Get(r, auth.SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.Values["oauth_state"] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	url := auth.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
