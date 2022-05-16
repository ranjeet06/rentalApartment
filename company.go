package gorsk

import "github.com/ribice/gorsk/pkg/utl/model"

// Company represents company model
type Company struct {
	model.Base
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Locations []Location `json:"locations,omitempty"`
	Owner     User       `json:"owner"`
}
