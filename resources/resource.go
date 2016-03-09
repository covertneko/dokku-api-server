package resources

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/nikelmwann/dokku-api/dokku"
)

type Resource interface {
	Get(c *echo.Context, s *dokku.Dokku) error
	Post(c *echo.Context, s *dokku.Dokku) error
	Put(c *echo.Context, s *dokku.Dokku) error
	Delete(c *echo.Context, s *dokku.Dokku) error
}

type (
	GetNotSupported struct{}
	PostNotSupported struct{}
	PutNotSupported struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(c * echo.Context, _ *dokku.Dokku) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (PostNotSupported) Post(c * echo.Context, _ *dokku.Dokku) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (PutNotSupported) Put(c * echo.Context, _ *dokku.Dokku) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func (DeleteNotSupported) Delete(c * echo.Context, _ *dokku.Dokku) error {
	return c.NoContent(http.StatusMethodNotAllowed)
}

func HandlerFor(r Resource, s *dokku.Dokku) echo.HandlerFunc {
	return func(c *echo.Context) error {
		switch c.Request().Method {
			case "GET":
				return r.Get(c, s)
			case "POST":
				return r.Post(c, s)
			case "PUT":
				return r.Put(c, s)
			case "DELETE":
				return r.Delete(c, s)
			default:
				return c.NoContent(http.StatusMethodNotAllowed)
		}
	}
}
