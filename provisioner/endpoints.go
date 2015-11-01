package provisioner

import "code.google.com/p/go-uuid/uuid"

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

func (p Provisioner) CreateEndpoint(region string) Endpoint {
	endpoint := newEndpoint()
	p.endpoints[endpoint.ID] = endpoint

	go p.provisionEndpoint(endpoint, region)

	return *endpoint
}

func (p Provisioner) GetEndpoint(id string) (Endpoint, error) {
	if endpoint, ok := p.endpoints[id]; ok {
		return *endpoint, nil
	}
	return Endpoint{}, ErrNotFound
}
