package user_favorite_apartment_platform

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
)

// FavApartmentDB represents userFavoriteApartment repository interface
type FavApartmentDB interface {
	AddFavApartmentRepository(echo.Context, user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error)
	ListFavApartmentRepository(echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error)
	ViewFavApartmentRepository(echo.Context, int) ([]user_favorite_apartment.UserFavoriteApartment, error)
	DeleteFavApartmentRepository(echo.Context, int, int) error
}
