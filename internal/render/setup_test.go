package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/easymomo/go-bookings/internal/config"
	"github.com/easymomo/go-bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// Things we can store in the session, primitives and string can be added by default but we need
	// to tell the application about structs we created ourselves
	gob.Register(models.Reservation{})

	// Change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// Create a type with interfaces that satisfies the requirements for a ResponseWriter
//
//	type ResponseWriter interface {
//	   Header() Header
//	   Write([]byte) (int, error)
//	   WriteHeader(statusCode int)
//	}
//
// A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
// A ResponseWriter may not be used after the Handler.ServeHTTP method has returned.
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (tw *myWriter) WriteHeader(statusCode int) {}
