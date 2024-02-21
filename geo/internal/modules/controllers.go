package modules

import (
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/component"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules/geo/controller"
)

type Controllers struct {
	Geo controller.Georer
}

func NewControllers(services *Services, components *component.Components) *Controllers {
	geoController := controller.NewGeoController(services.Geo, components)

	return &Controllers{
		Geo: geoController,
	}
}
