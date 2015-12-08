// Package web contains the HTTP endpoints used by the yovpn-server binary.
package web

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xperimental/yovpn/provisioner"
)

// CreateServer sets up a HTTP server for a provisioner.
func CreateServer(provisioner provisioner.Provisioner) http.Handler {
	s := &yovpnServer{provisioner}

	e := echo.New()
	s.setupHandlers(e)
	return e
}

type yovpnServer struct {
	provisioner provisioner.Provisioner
}

func (s *yovpnServer) setupHandlers(e *echo.Echo) {
	e.Get("/", s.blank)
	e.Get("/cleanup", s.cleanup)

	e.Put("/endpoint", s.createEndpoint)
	e.Get("/endpoint/:id", s.getEndpoint)
	e.Delete("/endpoint/:id", s.deleteEndpoint)

	e.Get("/regions", s.getRegions)
}

func (s *yovpnServer) createEndpoint(c *echo.Context) error {
	region := c.Query("region")
	if len(region) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "You need to provide a region!")
	}

	endpoint := s.provisioner.CreateEndpoint(region)
	c.Response().Header().Set(echo.Location, c.Echo().URI(s.getEndpoint, endpoint.ID))
	c.JSON(http.StatusAccepted, endpoint)
	return nil
}

func (s *yovpnServer) getEndpoint(c *echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID not provided!")
	}

	endpoint, err := s.provisioner.GetEndpoint(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Endpoint not found!")
	}

	c.JSON(http.StatusOK, endpoint)
	return nil
}

func (s *yovpnServer) deleteEndpoint(c *echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID not provided!")
	}

	_, err := s.provisioner.DestroyEndpoint(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Endpoint not found!")
	}

	c.NoContent(http.StatusNoContent)
	return nil
}

func (s *yovpnServer) cleanup(c *echo.Context) error {
	endpoints := s.provisioner.ListEndpoints()
	result := struct {
		Total  int
		Errors int
	}{
		Total:  len(endpoints),
		Errors: 0,
	}
	for _, endpoint := range endpoints {
		_, err := s.provisioner.DestroyEndpoint(endpoint.ID)
		if err != nil {
			result.Errors++
		}
	}

	c.JSON(http.StatusOK, result)
	return nil
}

func (s *yovpnServer) getRegions(c *echo.Context) error {
	regions, err := s.provisioner.ListRegions()
	if err != nil {
		return errors.New("Could not list regions!")
	}

	c.JSON(http.StatusOK, regions)
	return nil
}

func (s *yovpnServer) blank(c *echo.Context) error {
	c.String(http.StatusOK, "yovpn")

	return nil
}
