// Package user_favorite_apartment_service contains user favorite apartment application services
package user_favorite_apartment_service

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_platform"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_platform/user_favorite_apartment_pgsql"
	"github.com/rs/zerolog/log"
	"net/http"
)

// UserFavoriteApartment represents user favorite apartment application service
type UserFavoriteApartment struct {
	db               *pg.DB
	favApartmentDb   user_favorite_apartment_platform.FavApartmentDB
	apartmentService apartment_service.ApartmentService
	UserService      apartment_user_service.ApartmentUserService
}

// Initialize initalizes user favorite apartment application service with defaults
func Initialize(db *pg.DB, apartmentService apartment_service.ApartmentService, UserService apartment_user_service.ApartmentUserService) *UserFavoriteApartment {
	favApartment := user_favorite_apartment_pgsql.Initialize(db)
	return &UserFavoriteApartment{db: db, favApartmentDb: favApartment, apartmentService: apartmentService, UserService: UserService}
}
func InitializeMock(favApartment user_favorite_apartment_platform.FavApartmentDB, apartmentService apartment_service.ApartmentService, UserService apartment_user_service.ApartmentUserService) *UserFavoriteApartment {
	return &UserFavoriteApartment{db: nil, favApartmentDb: favApartment, apartmentService: apartmentService, UserService: UserService}
}

// AddUserFavApartmentService add a new apartment in user favorite list
func (fa UserFavoriteApartment) AddUserFavApartmentService(c echo.Context, favApartment user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error) {
	_, err := fa.UserService.ViewApartmentUserService(c, favApartment.UserID)
	if err != nil {
		log.Err(err)
		return favApartment, err
	}
	_, err = fa.apartmentService.ViewApartmentService(c, favApartment.ApartmentID)
	if err != nil {
		log.Err(err)
		return favApartment, err
	}

	return fa.favApartmentDb.AddFavApartmentRepository(c, favApartment)
}

// ListUserFavApartmentService returns list of favorite apartment
func (fa UserFavoriteApartment) ListUserFavApartmentService(c echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	return fa.favApartmentDb.ListFavApartmentRepository(c)
}

// ViewUserFavApartmentService returns favorite apartment for a user
func (fa UserFavoriteApartment) ViewUserFavApartmentService(c echo.Context, a int) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	return fa.favApartmentDb.ViewFavApartmentRepository(c, a)
}

// Custom error
var (
	ErrorNotExists = echo.NewHTTPError(http.StatusBadRequest, "user not exist.")
)

// DeleteUserFavApartmentService remove  apartment from favorite list
func (fa UserFavoriteApartment) DeleteUserFavApartmentService(c echo.Context, userId int, apartmentId int) error {

	return fa.favApartmentDb.DeleteFavApartmentRepository(c, userId, apartmentId)
}
