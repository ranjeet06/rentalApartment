// Package apartment_user_service contains apartment user application services
package apartment_user_service

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_platform"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_platform/user_pgsql"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/rs/zerolog/log"
)

// ApartmentUser represents apartment user application service
type ApartmentUser struct {
	db              *pg.DB
	apartmentUserDb apartment_user_platform.ApartmentUserDb
}

// InitializeUser initalizes Apartment user application service with defaults
func InitializeUser(db *pg.DB) *ApartmentUser {
	apartmentUser := user_pgsql.InitializeUser(db)
	return &ApartmentUser{db: db, apartmentUserDb: apartmentUser}
}

// InitializeUser initalizes Apartment user application service with defaults
func InitializeMockUser(db apartment_user_platform.ApartmentUserDb) *ApartmentUser {
	return &ApartmentUser{db: nil, apartmentUserDb: db}
}

// CreateApartmentUserService creates a new apartment user
func (a ApartmentUser) CreateApartmentUserService(c echo.Context, apartmentUser apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
	return a.apartmentUserDb.CreateApartmentUserRepository(c, apartmentUser)
}

// ListApartmentUserService returns list of apartment users
func (a ApartmentUser) ListApartmentUserService(c echo.Context, filter apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error) {
	return a.apartmentUserDb.ListApartmentUserRepository(c, filter)
}

// ViewApartmentUserService returns single apartment user
func (a ApartmentUser) ViewApartmentUserService(c echo.Context, b int) (apartment_user.ApartmentUser, error) {
	return a.apartmentUserDb.ViewApartmentUserRepository(c, b)
}

// NewUpdateApartmentUser contains apartment user's information used for updating
type NewUpdateApartmentUser struct {
	Id          int
	Name        string
	UserEmail   string
	UserAddress string
}

// UpdateApartmentUserService updates apartment user's information
func (a ApartmentUser) UpdateApartmentUserService(c echo.Context, updateUser NewUpdateApartmentUser) (apartment_user.ApartmentUser, error) {
	if err := a.apartmentUserDb.UpdateApartmentUserRepository(c, apartment_user.ApartmentUser{
		Base:        model.Base{ID: updateUser.Id},
		Name:        updateUser.Name,
		UserEmail:   updateUser.UserEmail,
		UserAddress: updateUser.UserAddress,
	}); err != nil {
		log.Err(err)
		return apartment_user.ApartmentUser{}, err
	}

	return a.apartmentUserDb.ViewApartmentUserRepository(c, updateUser.Id)

}

// DeleteApartmentUserService deletes a apartment user
func (a ApartmentUser) DeleteApartmentUserService(c echo.Context, b int) error {
	apartmentUser, err := a.apartmentUserDb.ViewApartmentUserRepository(c, b)
	if err != nil {
		log.Err(err)
		return err
	}
	return a.apartmentUserDb.DeleteApartmentUserRepository(c, apartmentUser)
}

func (a ApartmentUser) AccessTokenService(c echo.Context, token apartment_user.AccessToken) (string, error) {
	return a.apartmentUserDb.AccessTokenRepository(c, token)
}
