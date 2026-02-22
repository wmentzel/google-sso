package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google-sso/auth"
)

type googleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Callback handles the OAuth2 redirect from Google, validates the state,
// exchanges the code for a token, fetches the user profile, and stores the
// session before redirecting to /hello.
func Callback(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, auth.SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Validate state
	expectedState, _ := session.Values["oauth_state"].(string)
	if r.URL.Query().Get("state") != expectedState {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := auth.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch user info from Google
	client := auth.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read user info", http.StatusInternalServerError)
		return
	}

	var user googleUser
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Persist auth in session
	delete(session.Values, "oauth_state")
	session.Values["authenticated"] = true
	session.Values["user_email"] = user.Email
	session.Values["user_name"] = user.Name
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/hello", http.StatusTemporaryRedirect)
}
