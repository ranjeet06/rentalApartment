package apartment

import (
	"github.com/ribice/gorsk/pkg/utl/model"
)

// Vertex represents vertex for geolocation coordinates
type Vertex struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Pagination represents for pagination
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Apartment represents apartment object
type Apartment struct {
	model.Base
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	FloorArea     float64 `json:"floor_area"`
	PricePerMonth float64 `json:"price_per_month"`
	NumberOfRooms int     `json:"number_of_rooms"`

	GeolocationCoordinates Vertex `json:"geolocation_coordinates" gorm:"foreignKey:lat"`
	AssociatedRealtor      string `json:"associated_realtor"`
}

// FilterApartment represents for the filter
type FilterApartment struct {
	FloorArea     float64 `json:"floor_area"`
	PricePerMonth float64 `json:"price_per_month"`
	NumberOfRooms int     `json:"number_of_rooms"`
}
