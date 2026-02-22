package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const SessionName = "auth-session"

var Store *sessions.CookieStore
var OAuthConfig *oauth2.Config

func Init() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "change-me-in-production-32bytes!"
	}
	Store = sessions.NewCookieStore([]byte(secret))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// IsAuthenticated returns true if the request has a valid auth session.
func IsAuthenticated(r *http.Request) bool {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return false
	}
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}
