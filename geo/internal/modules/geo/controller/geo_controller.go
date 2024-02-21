package controller

import (
	"encoding/json"
	"net/http"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/component"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/responder"
	"projects/LDmitryLD/geo-nginx/geo/internal/models"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules/geo/service"
)

type Georer interface {
	Geocode(http.ResponseWriter, *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type GeoController struct {
	geo service.Georer
	responder.Responder
}

func NewGeoController(service service.Georer, components *component.Components) Georer {
	return &GeoController{
		geo:       service,
		Responder: components.Responder,
	}
}

func (g *GeoController) Geocode(w http.ResponseWriter, r *http.Request) {

	var geocodeRequest GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&geocodeRequest); err != nil {
		g.ErrorBadRequest(w, err)
		return
	}

	geo := g.geo.GeoCode(service.GeoCodeIn{Lat: geocodeRequest.Lat, Lng: geocodeRequest.Lng})

	geocodeResponse := GeocodeResponse{
		Addresses: []*models.Address{{Lat: geo.Lat, Lon: geo.Lng}},
	}

	g.OutputJSON(w, geocodeResponse)
}

func (g *GeoController) Search(w http.ResponseWriter, r *http.Request) {

	var searchRequest SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
		g.ErrorBadRequest(w, err)
		return
	}

	out := g.geo.SearchAddresses(service.SearchAddressesIn{Query: searchRequest.Query})
	if out.Err != nil {
		g.ErrorInternal(w, out.Err)
		return
	}

	searchResponse := SearchResponse{
		Addresses: []*models.Address{&out.Address},
	}

	g.OutputJSON(w, searchResponse)

}
