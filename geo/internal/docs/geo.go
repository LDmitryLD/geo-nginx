package docs

import "projects/LDmitryLD/geo-nginx/geo/internal/modules/geo/controller"

//swagger:route Post /api/address/geocode geocode GeocodeRequest
// Вычисление адресса по широте и долготе.
// responses:
//	200: GeocodeResponse

//swagger:parameters GeocodeRequest
type GeocodeRequest struct {
	// Lat - широта
	// Lng - долгота
	// in: body
	// required: true
	Body controller.GeocodeRequest
}

//swagger:response GeocodeResponse
type GeocodeResponse struct {
	// in: body
	Body controller.GeocodeResponse
}

//swagger:route Post /api/address/search search SearchRequest
// Вычисление местанахождения по адрессу.
// responses:
//	200: SearchResponse
//

//swagger:parameters SearchRequest
type SearchRequest struct {
	//Qury - запрос, представляющий собой адрес
	//in: body
	Body controller.SearchRequest
}

//swagger:response SearchResponse
type SearchResponse struct {
	// Addresses содержит список адрессов
	// in: body
	Body controller.SearchResponse
}
