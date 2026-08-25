module github.com/kohkimakimoto/echo-session/example

go 1.25.0

require (
	github.com/kohkimakimoto/echo-session/v5 v5.0.0-00010101000000-000000000000
	github.com/labstack/echo/v5 v5.3.1
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	golang.org/x/time v0.15.0 // indirect
)

replace github.com/kohkimakimoto/echo-session/v5 => ..
