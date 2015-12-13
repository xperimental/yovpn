package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/xperimental/yovpn/provisioner"
)

func mustNewRequest(method, urlStr string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return req
}

func TestBlank(t *testing.T) {
	e := echo.New()
	s := &yovpnServer{}
	req := mustNewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)
	err := s.blank(c)

	if err != nil {
		t.Errorf("Expected no error, but got %s", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %v", w.Code)
	}
	if w.Body.Len() == 0 {
		t.Error("Expected body, but got none.")
	}
}

type mockProvisioner struct {
	regions []provisioner.Region
}

func (p mockProvisioner) CreateEndpoint(string) provisioner.Endpoint {
	return provisioner.Endpoint{}
}

func (p mockProvisioner) GetEndpoint(string) (provisioner.Endpoint, error) {
	return provisioner.Endpoint{}, nil
}

func (p mockProvisioner) ListEndpoints() []provisioner.Endpoint {
	return []provisioner.Endpoint{}
}

func (p mockProvisioner) DestroyEndpoint(string) (provisioner.Endpoint, error) {
	return provisioner.Endpoint{}, nil
}

func (p mockProvisioner) ListRegions() ([]provisioner.Region, error) {
	return p.regions, nil
}

func (p mockProvisioner) Signal() chan struct{} {
	return nil
}

func createProvisioner(regions []provisioner.Region) mockProvisioner {
	return mockProvisioner{
		regions: regions,
	}
}

func TestListRegions(t *testing.T) {
	for _, test := range []struct {
		p    mockProvisioner
		err  error
		code int
		body string
	}{
		{
			createProvisioner([]provisioner.Region{}),
			nil,
			200,
			"[]",
		},
		{
			createProvisioner([]provisioner.Region{{"name", "description", "country"}}),
			nil,
			200,
			"[{\"name\":\"name\",\"description\":\"description\",\"country\":\"country\"}]",
		},
	} {
		s := &yovpnServer{test.p}

		e := echo.New()
		req := mustNewRequest("GET", "/regions", nil)
		w := httptest.NewRecorder()
		c := echo.NewContext(req, echo.NewResponse(w, e), e)
		err := s.getRegions(c)

		assert.Equal(t, test.err, err)
		assert.Equal(t, test.code, w.Code)
		assert.Equal(t, test.body, w.Body.String())
	}
}
