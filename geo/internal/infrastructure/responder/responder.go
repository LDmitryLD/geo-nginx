package responder

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type Respond struct {
	log *zap.Logger
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{log: logger}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		//log.Println("Ошибка при отправке ответа:", err.Error())
		r.log.Error("responder json encode error", zap.Error(err))
		r.ErrorInternal(w, err)
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.log.Error("http response internal error", zap.Error(err))
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
