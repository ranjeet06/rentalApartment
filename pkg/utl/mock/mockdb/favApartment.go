package mockdb

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
)

type FavApartment struct {
	AddFavApartmentFn    func(echo.Context, user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error)
	ListFavApartmentFn   func(echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error)
	ViewFavApartmentFn   func(echo.Context, int) ([]user_favorite_apartment.UserFavoriteApartment, error)
	DeleteFavApartmentFn func(echo.Context, int, int) error
}

func (fdb FavApartment) AddFavApartmentRepository(c echo.Context, favApartment user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error) {
	return fdb.AddFavApartmentFn(c, favApartment)
}

// ListFavApartmentRepository returns list of all user favorite apartment.
func (fdb FavApartment) ListFavApartmentRepository(c echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	return fdb.ListFavApartmentFn(c)
}

// ViewFavApartmentRepository returns list of single user favorite apartment.
func (fdb FavApartment) ViewFavApartmentRepository(c echo.Context, id int) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	return fdb.ViewFavApartmentFn(c, id)
}

// DeleteFavApartmentRepository remove from favorites
func (fdb FavApartment) DeleteFavApartmentRepository(c echo.Context, userId int, apartmentId int) error {

	return fdb.DeleteFavApartmentFn(c, userId, apartmentId)
}
