package handlers

import (
	"fmt"
	"net/http"

	"google-sso/auth"
)

// Hello is a protected endpoint that greets the authenticated user.
func Hello(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, auth.SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	name, _ := session.Values["user_name"].(string)
	email, _ := session.Values["user_email"].(string)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Hello World</title></head>
<body>
  <h1>Hello, World! 👋</h1>
  <p>Welcome, <strong>%s</strong> (%s)</p>
  <p><a href="/logout">Logout</a></p>
</body>
</html>`, name, email)
}
