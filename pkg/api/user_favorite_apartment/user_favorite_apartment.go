package user_favorite_apartment

import "github.com/ribice/gorsk/pkg/utl/model"

// UserFavoriteApartment represents user favorite apartment object
type UserFavoriteApartment struct {
	model.Base
	UserID      int `json:"user_id"`
	ApartmentID int `json:"apartment_id"`
}
