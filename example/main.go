package main

import (
	"fmt"
	"log/slog"
	"net/http"

	session "github.com/kohkimakimoto/echo-session/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.Use(session.Middleware(session.NewCookieStore([]byte("12345678901234567890123456789012"))))

	e.GET("/", func(c *echo.Context) error {
		s := session.MustGet(c)
		counter := 0
		if val := s.Get("counter"); val != nil {
			if count, ok := val.(int); ok {
				counter = count
			}
		}
		counter++

		s.Set("counter", counter)
		if err := s.Save(); err != nil {
			return err
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("Counter: %d", counter))
	})

	e.GET("/refresh", func(c *echo.Context) error {
		s := session.MustGet(c)
		s.Clear()
		if err := s.Save(); err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/")
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
