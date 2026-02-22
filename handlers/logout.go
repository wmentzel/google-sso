package handlers

import (
	"net/http"

	"google-sso/auth"
)

// Logout clears the auth session and redirects to /login.
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, auth.SessionName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Expire the session immediately
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	_ = session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
