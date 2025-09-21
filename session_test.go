package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMustGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := MiddlewareWithConfig(MiddlewareConfig{
		Store: NewCookieStore([]byte("12345678901234567890123456789012")),
	})(func(c echo.Context) error {
		s := MustGet(c)
		assert.NotNil(t, s)
		return c.String(http.StatusOK, "test")
	})
	err := h(c)
	assert.NoError(t, err)
}

func TestSession_Save(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := MiddlewareWithConfig(MiddlewareConfig{
		Store: NewCookieStore([]byte("12345678901234567890123456789012")),
	})(func(c echo.Context) error {
		s := MustGet(c)
		assert.NotNil(t, s)

		s.Set("test", "test-value")
		if err := s.Save(); err != nil {
			return err
		}
		return c.String(http.StatusOK, "test")
	})
	err := h(c)
	assert.Regexp(t, "^session=.+", rec.Header().Get("Set-Cookie"))
	assert.NoError(t, err)
}
