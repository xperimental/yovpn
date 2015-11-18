package provisioner

import "github.com/digitalocean/godo"

type Region struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Country     string `json:"country"`
}

func containsSize(region godo.Region, search string) bool {
	for _, size := range region.Sizes {
		if size == search {
			return true
		}
	}
	return false
}

var countryMap = map[string]string{
	"sfo1": "us",
	"nyc1": "us",
	"nyc2": "us",
	"nyc3": "us",
	"ams1": "nl",
	"ams2": "nl",
	"ams3": "nl",
	"sgp1": "sg",
	"lon1": "uk",
	"fra1": "de",
	"tor1": "ca",
}

func country(region godo.Region) string {
	country, ok := countryMap[region.Slug]
	if ok {
		return country
	}
	return ""
}

func (p provisioner) ListRegions() ([]Region, error) {
	doRegions, _, err := p.client.Regions.List(&godo.ListOptions{})
	if err != nil {
		return nil, err
	}

	var regions []Region
	for _, region := range doRegions {
		if region.Available && containsSize(region, defaultSize) {
			regions = append(regions, Region{
				Name:        region.Slug,
				Description: region.Name,
				Country:     country(region),
			})
		}
	}
	return regions, nil
}
