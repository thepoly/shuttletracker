package api

import (
	"net/http"
	"net/url"
	"strings"

	gc "gopkg.in/cas.v2"

	"github.com/thepoly/shuttletracker"
	"github.com/thepoly/shuttletracker/auth"
	"github.com/thepoly/shuttletracker/log"
)

// CASClient stores the local cas client and an instance of the database
type CASClient struct {
	cas          auth.AuthenticationService
	us           shuttletracker.UserService
	authenticate bool
}

// CreateCASClient creates an authentication service CASClient using a cas url and database
func CreateCASClient(url *url.URL, us shuttletracker.UserService, authenticate bool) *CASClient {
	client := gc.NewClient(&gc.Options{
		URL:   url,
		Store: nil,
	})

	cli := &CASClient{
		cas: &auth.CAS{
			CAS: client,
		},
		us:           us,
		authenticate: authenticate,
	}
	return cli
}

// InjectMocks allows mock interfaces to be used
func InjectMocks(cli auth.AuthenticationService, us shuttletracker.UserService, auth bool) *CASClient {
	c := &CASClient{
		cas:          cli,
		us:           us,
		authenticate: auth,
	}
	return c
}

func (cli *CASClient) logout(w http.ResponseWriter, r *http.Request) {
	cli.cas.Logout(w, r)
}

func (cli *CASClient) casauth(next http.Handler) http.Handler {
	return cli.cas.HandleFunc(func(w http.ResponseWriter, r *http.Request) {

		if !cli.authenticate {
			// don't authenticate
			next.ServeHTTP(w, r)
			return
		}

		if !cli.cas.Authenticated(r) {
			_, err := w.Write([]byte("redirecting to cas;"))
			if err != nil {
				log.WithError(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			cli.cas.Login(w, r)
		} else {
			auth, err := cli.us.UserExists(strings.ToLower(cli.cas.Username(r)))
			if err != nil {
				log.WithError(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				auth = false
			}
			if auth {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "unauthenticated", 401)

		}

	})
}
