package provisioner

import (
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/pborman/uuid"
)

const (
	Starting  = "starting"
	Running   = "running"
	Failed    = "failed"
	Destroyed = "destroyed"
)

type Endpoint struct {
	ID        string `json:"id"`
	IP        string `json:"-"`
	Config    string `json:"config"`
	DropletID int    `json:"-"`
	Status    string `json:"status"`
}

func newEndpoint() *Endpoint {
	return &Endpoint{
		ID:        uuid.New(),
		IP:        "",
		Config:    "",
		DropletID: 0,
		Status:    Starting,
	}
}

func (p Provisioner) restoreEndpoints() {
	droplets, _, err := p.client.Droplets.List(&godo.ListOptions{})
	if err != nil {
		log.Printf("Error restoring endpoints: %s", err)
		return
	}

	for _, droplet := range droplets {
		if strings.HasPrefix(droplet.Name, baseName) {
			id := strings.TrimPrefix(droplet.Name, baseName)
			if _, err := p.GetEndpoint(id); err == ErrNotFound {
				ip := droplet.Networks.V4[0].IPAddress
				endpoint := &Endpoint{
					ID:        id,
					IP:        ip,
					Config:    "",
					DropletID: droplet.ID,
					Status:    Running,
				}
				p.endpoints[id] = endpoint
				log.Printf("Recovered endpoint with id %s", id)
			}
		}
	}
}

func (p Provisioner) CreateEndpoint(region string) Endpoint {
	endpoint := newEndpoint()
	p.endpoints[endpoint.ID] = endpoint

	go p.provisionEndpoint(endpoint, region)

	return *endpoint
}

func (p Provisioner) ListEndpoints() []Endpoint {
	var result []Endpoint
	for _, endpoint := range p.endpoints {
		result = append(result, *endpoint)
	}
	return result
}

func (p Provisioner) GetEndpoint(id string) (Endpoint, error) {
	if endpoint, ok := p.endpoints[id]; ok {
		return *endpoint, nil
	}
	return Endpoint{}, ErrNotFound
}

func (p Provisioner) DestroyEndpoint(id string) (Endpoint, error) {
	if endpoint, ok := p.endpoints[id]; ok {
		err := deleteDroplet(p.client, endpoint.DropletID)
		if err != nil {
			return Endpoint{}, err
		}
		endpoint.Status = Destroyed
		return *endpoint, nil
	}
	return Endpoint{}, ErrNotFound
}
