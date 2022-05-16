package gorsk

import "github.com/ribice/gorsk/pkg/utl/model"

// Location represents company location model
type Location struct {
	model.Base
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Address string `json:"address"`

	CompanyID int `json:"company_id"`
}
