package transport

import (
	"net/http"
)

func (t *Transporter) JSONResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return t.WriteJSON(w, status, &envelope{Data: data})
}

func (t *Transporter) JSONResponseWithMetadata(w http.ResponseWriter, status int, data any, metadata any) error {
	type envelope struct {
		Data     any `json:"data"`
		Metadata any `json:"metadata"`
	}

	return t.WriteJSON(w, status, &envelope{Data: data, Metadata: metadata})
}

func (t Transporter) EmptyResponse(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}
