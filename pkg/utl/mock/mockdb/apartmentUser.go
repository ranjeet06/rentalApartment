package mockdb

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
)

type ApartmentUser struct {
	CreateApartmentUserFn func(echo.Context, apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error)
	ListApartmentUserFn   func(echo.Context, apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error)
	ViewApartmentUserFn   func(echo.Context, int) (apartment_user.ApartmentUser, error)
	UpdateApartmentUserFn func(echo.Context, apartment_user.ApartmentUser) error
	DeleteApartmentUserFn func(echo.Context, apartment_user.ApartmentUser) error
}

func (apu ApartmentUser) CreateApartmentUserRepository(c echo.Context, apartmentUserNew apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
	return apu.CreateApartmentUserFn(c, apartmentUserNew)
}

// ListApartmentUserRepository returns list of all apartment users.
func (apu ApartmentUser) ListApartmentUserRepository(c echo.Context, filter apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error) {
	return apu.ListApartmentUserFn(c, filter)
}

// ViewApartmentUserRepository returns single apartment user by ID
func (apu ApartmentUser) ViewApartmentUserRepository(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
	return apu.ViewApartmentUserFn(c, id)
}

// UpdateApartmentUserRepository updates apartment user's info
func (apu ApartmentUser) UpdateApartmentUserRepository(c echo.Context, updateUser apartment_user.ApartmentUser) error {
	return apu.UpdateApartmentUserFn(c, updateUser)
}

// DeleteApartmentUserRepository sets deleted_at for a apartment user
func (apu ApartmentUser) DeleteApartmentUserRepository(c echo.Context, deleteUser apartment_user.ApartmentUser) error {
	return apu.DeleteApartmentUserFn(c, deleteUser)
}
