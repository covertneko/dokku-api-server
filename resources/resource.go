package resources

import (
	"net/http"
	"github.com/labstack/echo"
)

type Resource interface {
	Get(c *echo.Context) error
	Post(c *echo.Context) error
	Put(c *echo.Context) error
	Delete(c *echo.Context) error
}

type (
	GetNotSupported struct{}
	PostNotSupported struct{}
	PutNotSupported struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(c * echo.Context) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (PostNotSupported) Post(c * echo.Context) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (PutNotSupported) Put(c * echo.Context) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (DeleteNotSupported) Delete(c * echo.Context) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func HandlerFor(r Resource) echo.HandlerFunc {
	return func(c *echo.Context) error {
		switch c.Request().Method {
			case "GET":
				return r.Get(c)
			case "POST":
				return r.Post(c)
			case "PUT":
				return r.Put(c)
			case "DELETE":
				return r.Delete(c)
			default:
				return c.NoContent(http.StatusMethodNotAllowed)
		}
	}
}
