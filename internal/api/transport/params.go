package transport

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (t *Transporter) GetUrlParam(r *http.Request, key string) (*string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return nil, http.ErrMissingFile
	}
	return &value, nil
}

func (t *Transporter) GetUrlParamInt(r *http.Request, key string) (*int, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return nil, http.ErrMissingFile
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &intValue, nil
}

func (t *Transporter) GetSearchParam(r *http.Request, key string) (*string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, http.ErrMissingFile
	}
	return &value, nil
}

func (t *Transporter) GetSearchParamInt(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, http.ErrMissingFile
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &intValue, nil
}

func (t *Transporter) GetPaginationParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	limitStr := query.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 1
	}

	offsetStr := query.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return offset, limit, nil
}
