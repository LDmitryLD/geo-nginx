package modules

import (
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/component"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules/geo/service"
	"projects/LDmitryLD/geo-nginx/geo/internal/storages"
)

type Services struct {
	Geo service.Georer
}

func NewServices(storages *storages.Storages, components *component.Components) *Services {
	geoService := service.NewGeo(storages.Geo, components.Logger)

	return &Services{
		Geo: geoService,
	}
}
