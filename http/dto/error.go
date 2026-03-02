package dto

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/iancoleman/strcase"
)

type ErrResponse struct {
	Err     error  `json:"-"`
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (resp *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if resp.Status == 0 {
		resp.Status = http.StatusInternalServerError
	}

	if resp.Message == "" {
		resp.Message = http.StatusText(resp.Status)
	}

	if resp.Code == "" {
		resp.Code = strcase.ToSnake(http.StatusText(resp.Status))
	}

	render.Status(r, resp.Status)

	return nil
}
