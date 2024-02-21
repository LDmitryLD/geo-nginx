package storage

import (
	"projects/LDmitryLD/geo-nginx/geo/internal/db/adapter"
	"projects/LDmitryLD/geo-nginx/geo/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=GeoStorager
type GeoStorager interface {
	Select(query string) (models.Address, error)
	Insert(query, lat, lon string) error
}

type GeoStorage struct {
	adapter adapter.SQLAdapterer
}

func NewGeoStorage(sqlAdapter adapter.SQLAdapterer) *GeoStorage {
	return &GeoStorage{
		adapter: sqlAdapter,
	}
}

func (g *GeoStorage) Select(query string) (models.Address, error) {
	address, err := g.adapter.Select(query)

	return address, err
}

func (g *GeoStorage) Insert(query, lat, lon string) error {
	return g.adapter.Insert(query, lat, lon)
}
