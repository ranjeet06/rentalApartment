package apartment_user_platform

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
)

// ApartmentUserDb represents apartment user repository interface
type ApartmentUserDb interface {
	CreateApartmentUserRepository(echo.Context, apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error)
	ListApartmentUserRepository(echo.Context, apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error)
	ViewApartmentUserRepository(echo.Context, int) (apartment_user.ApartmentUser, error)
	UpdateApartmentUserRepository(echo.Context, apartment_user.ApartmentUser) error
	DeleteApartmentUserRepository(echo.Context, apartment_user.ApartmentUser) error
	AccessTokenRepository(echo.Context, apartment_user.AccessToken) (string, error)
}
