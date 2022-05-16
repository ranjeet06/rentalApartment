package user_favorite_apartment_pgsql

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	"github.com/rs/zerolog/log"
	"net/http"
)

// FavoriteApartment represents for userFavoriteApartment repository
type FavoriteApartment struct {
	db orm.DB
}

// Initialize initializes userFavoriteApartment repository with default
func Initialize(db orm.DB) *FavoriteApartment {
	return &FavoriteApartment{db: db}
}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusConflict, "apartment already exist in favorite.")
)

// Custom errors
var (
	ErrNotExists = echo.NewHTTPError(http.StatusNotFound, "user id not exist in favorite.")
)

// AddFavApartmentRepository add a favorite apartment  on database
func (fdb FavoriteApartment) AddFavApartmentRepository(c echo.Context, favApartment user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error) {
	var newFavApartment = new(user_favorite_apartment.UserFavoriteApartment)

	ctx := c.Request().Context()
	err := fdb.db.Model(newFavApartment).Context(ctx).Where("apartment_id = ?", favApartment.ApartmentID).Select()
	if err != nil {
		log.Err(err)
	}
	if err == nil || err != pg.ErrNoRows {
		return user_favorite_apartment.UserFavoriteApartment{}, ErrAlreadyExists
	}

	err = fdb.db.Insert(&favApartment)
	log.Err(err)
	return favApartment, err

}

// ListFavApartmentRepository returns list of all user favorite apartment.
func (fdb FavoriteApartment) ListFavApartmentRepository(c echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	var favApartment []user_favorite_apartment.UserFavoriteApartment

	ctx := c.Request().Context()
	err := fdb.db.Model(&favApartment).Context(ctx).Select()
	log.Err(err)
	return favApartment, err
}

// ViewFavApartmentRepository returns list of single user favorite apartment.
func (fdb FavoriteApartment) ViewFavApartmentRepository(c echo.Context, id int) ([]user_favorite_apartment.UserFavoriteApartment, error) {
	var favApartment []user_favorite_apartment.UserFavoriteApartment

	sql := `SELECT "user_id" , "apartment_id" FROM "user_favorite_apartments"
			WHERE "user_favorite_apartments"."user_id" = ?;`

	a, err := fdb.db.Query(&favApartment, sql, id)
	if err != nil {
		log.Err(err)
	}
	if a.RowsReturned() == 0 {
		err = ErrNotExists
	}
	return favApartment, err
}

// DeleteFavApartmentRepository remove from favorites
func (fdb FavoriteApartment) DeleteFavApartmentRepository(c echo.Context, userId int, apartmentId int) error {

	var favApartment user_favorite_apartment.UserFavoriteApartment
	sql := `DELETE FROM "user_favorite_apartments"
			WHERE "user_favorite_apartments"."user_id" = ? AND "user_favorite_apartments"."apartment_id" = ?;`

	_, err := fdb.db.Query(&favApartment, sql, userId, apartmentId)
	log.Err(err)
	return err
}
