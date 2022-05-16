package user_favorite_apartment_service

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
)

// UserFavApartmentService represents user favorite apartment application interface
type UserFavApartmentService interface {
	AddUserFavApartmentService(echo.Context, user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error)
	ListUserFavApartmentService(echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error)
	ViewUserFavApartmentService(echo.Context, int) ([]user_favorite_apartment.UserFavoriteApartment, error)
	DeleteUserFavApartmentService(echo.Context, int, int) error
}
