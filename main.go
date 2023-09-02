package main

import (
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	m := NewMiddleware()
	e.GET("/", m.SomeMiddleware(SomeHandler, m.SomeErrorHandler))

	e.Logger.Fatal(e.Start(":1323"))
}

type (
	MiddlewareHandler interface {
		SomeMiddleware(next, stop echo.HandlerFunc) echo.HandlerFunc
		SomeErrorHandler(c echo.Context) error
	}

	middleware struct {
		code    int
		message string
	}
)

func NewMiddleware() MiddlewareHandler {
	return &middleware{}
}

func (m *middleware) SetError(code int, message string) {
	m.code = code
	m.message = message
}

func (m *middleware) SomeMiddleware(next, stop echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		n := rand.Intn(100)
		if n%2 == 0 {
			m.SetError(http.StatusBadRequest, "Custom Error")
			return stop(c)
		}
		return next(c)
	}
}

func (m *middleware) SomeErrorHandler(c echo.Context) error {
	return c.JSON(
		m.code,
		map[string]any{"message": m.message},
	)
}

// Handler
func SomeHandler(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{"message": "Hello, 世界 !"},
	)
}
