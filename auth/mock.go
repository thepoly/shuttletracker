package auth

import (
	"github.com/stretchr/testify/mock"
	"github.com/thepoly/shuttletracker/log"
	"net/http"
)

// Mock implements the auth interface.
type Mock struct {
	mock.Mock
}

// Authenticated returns the mock response to the server
func (auth *Mock) Authenticated(request *http.Request) bool {
	return true
}

// Logout writes logout to the ResponseWriter
func (auth *Mock) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("logout"))
	if err != nil {
		log.WithError(err)
	}

}

// Login writes login to the ResponseWriter
func (auth *Mock) Login(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("login"))
	if err != nil {
		log.WithError(err)
	}
}

// Username returns the mock response to the server
func (auth *Mock) Username(request *http.Request) string {
	return "lyonj4"
}

// HandleFunc returns an http handler for the request
func (auth *Mock) HandleFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}
