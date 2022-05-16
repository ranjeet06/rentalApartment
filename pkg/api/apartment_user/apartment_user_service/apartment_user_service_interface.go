package apartment_user_service

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
)

// ApartmentUserService represents apartment user application interface
type ApartmentUserService interface {
	CreateApartmentUserService(echo.Context, apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error)
	ListApartmentUserService(echo.Context, apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error)
	ViewApartmentUserService(echo.Context, int) (apartment_user.ApartmentUser, error)
	DeleteApartmentUserService(echo.Context, int) error
	UpdateApartmentUserService(echo.Context, NewUpdateApartmentUser) (apartment_user.ApartmentUser, error)
	AccessTokenService(echo.Context, apartment_user.AccessToken) (string, error)
}
