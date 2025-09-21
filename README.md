# Echo Session

[![test](https://github.com/kohkimakimoto/echo-session/actions/workflows/test.yml/badge.svg)](https://github.com/kohkimakimoto/echo-session/actions/workflows/test.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kohkimakimoto/echo-session/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/kohkimakimoto/echo-session.svg)](https://pkg.go.dev/github.com/kohkimakimoto/echo-session)


This is session middleware for [Echo](https://github.com/labstack/echo), provided as an alternative implementation inspired by [labstack/echo-contrib/session](https://github.com/labstack/echo-contrib/tree/master/session).

## Usage

```go
package main

import (
	"fmt"
	"net/http"

	session "github.com/kohkimakimoto/echo-session"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(session.Middleware(session.NewCookieStore([]byte("12345678901234567890123456789012"))))

	e.GET("/", func(c echo.Context) error {
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

	e.GET("/refresh", func(c echo.Context) error {
		s := session.MustGet(c)
		s.Clear()
		if err := s.Save(); err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
