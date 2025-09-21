package session

import (
	gorillasessions "github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Session struct {
	session     *gorillasessions.Session
	echoContext echo.Context
}

func (s *Session) Internal() *gorillasessions.Session {
	return s.session
}

func (s *Session) Set(key string, value any) {
	s.session.Values[key] = value
}

func (s *Session) Get(key string) any {
	return s.session.Values[key]
}

func (s *Session) Values() map[any]any {
	return s.session.Values
}

func (s *Session) Options() *gorillasessions.Options {
	return s.session.Options
}

func (s *Session) IsNew() bool {
	return s.session.IsNew
}

func (s *Session) AddFlash(value any, vars ...string) {
	s.session.AddFlash(value, vars...)
}

func (s *Session) Flashes(vars ...string) []any {
	return s.session.Flashes(vars...)
}

func (s *Session) Clear() {
	s.session.Values = make(map[any]any)
	s.session.Options.MaxAge = -1
}

func (s *Session) Save() error {
	return s.session.Save(s.echoContext.Request(), s.echoContext.Response())
}
