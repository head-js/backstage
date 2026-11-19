package vault

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zalando/go-keyring"
)

func NewHandler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveGet(store, w, r)
	})
}

func serveGet(store *Store, w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	n := r.URL.Query().Get("n")
	k := r.URL.Query().Get("k")
	role := r.URL.Query().Get("role")

	if secret == "" || n == "" || k == "" {
		log.Printf("vault: missing param: secret present=%t n present=%t k present=%t", secret != "", n != "", k != "")
		writeUnauthorized(w)
		return
	}

	if !store.Lookup(secret) {
		log.Printf("vault: secret miss")
		writeUnauthorized(w)
		return
	}

	value, err := keyring.Get(n, k)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			log.Printf("vault: keyring not found: n=%s k=%s", n, k)
		} else {
			log.Printf("vault: keyring error: %v", err)
		}
		writeUnauthorized(w)
		return
	}

	store.Append(secret, role, n, k)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"value": value}); err != nil {
		log.Printf("vault: encode response: %v", err)
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"}); err != nil {
		log.Printf("vault: encode error response: %v", err)
	}
}
