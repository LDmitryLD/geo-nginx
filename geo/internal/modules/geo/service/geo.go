package service

import (
	"context"
	"projects/LDmitryLD/geo-nginx/geo/internal/models"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules/geo/storage"

	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/client"
	"go.uber.org/zap"
)

type Geo struct {
	storage storage.GeoStorager
	logger  *zap.Logger
}

func NewGeo(storage storage.GeoStorager, logger *zap.Logger) Georer {
	return &Geo{storage: storage, logger: logger}
}

func (g *Geo) GeoCode(in GeoCodeIn) GeoCodeOut {
	return GeoCodeOut{
		Lat: in.Lat,
		Lng: in.Lng,
	}
}

func (g *Geo) SearchAddresses(in SearchAddressesIn) SearchAddressesOut {

	address, err := g.storage.Select(in.Query)
	if err != nil {
		res, err := searchFromAPI(in.Query)
		if err != nil {
			return SearchAddressesOut{
				Err: err,
			}
		}

		if err = g.storage.Insert(in.Query, res.Lat, res.Lon); err != nil {
			g.logger.Error("ошибка при добавлении данных в бд:", zap.Error(err))
		} else {
			g.logger.Error("Данные добавлены в бд", zap.Error(err))
		}

		return SearchAddressesOut{
			Address: res,
		}
	}

	return SearchAddressesOut{
		Address: address,
		Err:     nil,
	}
}

func searchFromAPI(query string) (models.Address, error) {

	api := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "d538755936a28def6bca48517dd287303cb0dae7",
		SecretKeyValue: "81081aa1fa5ca90caa8a69b14947b5876f58b8db",
	}))

	addresses, err := api.Address(context.Background(), query)
	if err != nil {
		return models.Address{}, err
	}

	res := models.Address{
		Lat: addresses[0].GeoLat,
		Lon: addresses[0].GeoLon,
	}

	return res, nil
}
