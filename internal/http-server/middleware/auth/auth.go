package auth

import (
	"encoding/json"
	"net/http"
)

func BasicAuthMiddleware(handler http.Handler, user, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providedUser, providedPassword, ok := r.BasicAuth()
		if !ok || providedUser != user || providedPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorMessage := map[string]string{"error": "Unauthorized"}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
