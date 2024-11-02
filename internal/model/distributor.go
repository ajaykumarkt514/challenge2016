package model

import (
	"challenge2016/internal/errors"
	"net/http"
	"strings"
)

type Permission struct {
	Include map[string]*Country
	Exclude map[string]*Country
}

type Distributor struct {
	Name    string   `json:"name"`
	Parent  *string  `json:"parent,omitempty"` // Parent is optional
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type DistributorResponse struct {
	Name      string              `json:"name"`
	Locations map[string]*Country `json:"locations"`
}

type DistributorAccess struct {
	Name string
	Permission
}

func (d *Distributor) Validate() error {
	d.Name = Sanitize(d.Name)

	if d.Name == "" {
		return &errors.Response{
			Code:   http.StatusBadRequest,
			Reason: "Invalid Name",
		}
	}

	if d.Parent != nil {
		*d.Parent = Sanitize(*d.Parent)

		if *d.Parent == "" {
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: "Invalid Parent Name",
			}
		}
	}

	return nil
}

// Sanitize trims leading and trailing spaces from the input string and converts it to uppercase.
func Sanitize(s string) string {
	return strings.TrimSpace(strings.ToUpper(s))
}
