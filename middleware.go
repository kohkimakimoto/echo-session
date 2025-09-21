package session

import (
	"fmt"
	"net/http"

	gorillasessions "github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const contextKey = "_session"

type InvalidSessionErrorHandlerFunc func(err error, store gorillasessions.Store, name string, c echo.Context) error

func DefaultInvalidSessionErrorHandler(err error, store gorillasessions.Store, name string, c echo.Context) error {
	// Remove the invalid session cookie
	http.SetCookie(c.Response(), &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})

	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid session: %v", err))
}

func Middleware(store gorillasessions.Store) echo.MiddlewareFunc {
	return MiddlewareWithConfig(MiddlewareConfig{
		Store: store,
	})
}

type MiddlewareConfig struct {
	Skipper middleware.Skipper
	// Store is the session store used to get and save sessions.
	Store gorillasessions.Store
	// Name is the name of the session cookie. The default is "session".
	Name string
	// InvalidSessionErrorHandler is called when an invalid session is detected.
	InvalidSessionErrorHandler InvalidSessionErrorHandlerFunc
}

func MiddlewareWithConfig(config MiddlewareConfig) echo.MiddlewareFunc {
	if config.Store == nil {
		panic("session: Store must be provided")
	}
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.Name == "" {
		config.Name = "session"
	}
	if config.InvalidSessionErrorHandler == nil {
		config.InvalidSessionErrorHandler = DefaultInvalidSessionErrorHandler
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Unlike the `labstack/echo-contrib/session` package, this middleware retrieves the session instead of setting a store in the context.
			// This is by design to perform fundamental validation of the session, such as detecting cookie tampering, before using it in the handler.
			sess, err := config.Store.Get(c.Request(), config.Name)
			if err != nil {
				// This error typically occurs when the session cookie is invalid
				return config.InvalidSessionErrorHandler(err, config.Store, config.Name, c)
			}

			c.Set(contextKey, &Session{
				session:     sess,
				echoContext: c,
			})
			return next(c)
		}
	}
}

var ErrNoSession = fmt.Errorf("no session in context")

// Get retrieves the session from the echo.Context.
func Get(c echo.Context) (*Session, error) {
	sess, ok := c.Get(contextKey).(*Session)
	if !ok {
		return nil, ErrNoSession
	}
	return sess, nil
}

// MustGet retrieves the session from the echo.Context and panics if it fails.
// If your handlers run after the session middleware, the session should always be available.
// This function is provided as a convenience to avoid error handling in such cases.
func MustGet(c echo.Context) *Session {
	sess, err := Get(c)
	if err != nil {
		panic(err)
	}
	return sess
}

// NewCookieStore creates a new CookieStore using the given secret key.
// It is a thin wrapper around gorilla/sessions.NewCookieStore with recommended options.
func NewCookieStore(secret []byte) *gorillasessions.CookieStore {
	s := gorillasessions.NewCookieStore(secret)
	// It is recommended to set Secure and HttpOnly to true for security reasons.
	s.Options.HttpOnly = true
	s.Options.Secure = true
	return s
}
