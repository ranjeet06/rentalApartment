package user_pgsql

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ApartmentUser represents for apartment user repository
type ApartmentUser struct {
	db orm.DB
}

// InitializeUser initializes apartment user repository with default
func InitializeUser(db orm.DB) *ApartmentUser {
	return &ApartmentUser{db: db}
}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusConflict, "user already exist.")
)

// Custom errors
var (
	ErrNotExists = echo.NewHTTPError(http.StatusNotFound, "user Id Not exist.")
)

// CreateApartmentUserRepository creates a new apartment user on database
func (apu ApartmentUser) CreateApartmentUserRepository(c echo.Context, apartmentUserNew apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
	var newApartmentUser = new(apartment_user.ApartmentUser)

	ctx := c.Request().Context()
	err := apu.db.Model(newApartmentUser).Context(ctx).Where("name = ?", apartmentUserNew.Name).Select()
	if err != nil {
		log.Err(err)
	}
	if err == nil || err != pg.ErrNoRows {
		return apartment_user.ApartmentUser{}, ErrAlreadyExists
	}

	err = apu.db.Insert(&apartmentUserNew)
	if err != nil {
		log.Err(err)
	}
	return apartmentUserNew, err
}

// ListApartmentUserRepository returns list of all apartment users.
func (apu ApartmentUser) ListApartmentUserRepository(c echo.Context, filter apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error) {
	var apartmentUser []apartment_user.ApartmentUser
	var q *orm.Query
	ctx := c.Request().Context()
	q = apu.db.Model(&apartmentUser).Context(ctx).Where("deleted_at is null")

	if filter.Name != "" {
		q.Where("name = ?", filter.Name)
	}

	if filter.UserEmail != "" {
		q.Where("user_email = ?", filter.UserEmail)
	}

	err := q.Select()
	if apartmentUser == nil {
		err = ErrNotExists
	}
	if err != nil {
		log.Err(err)
	}
	return apartmentUser, err
}

// ViewApartmentUserRepository returns single apartment user by ID
func (apu ApartmentUser) ViewApartmentUserRepository(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
	var apartmentUser apartment_user.ApartmentUser

	ctx := c.Request().Context()
	err := apu.db.Model(&apartmentUser).Context(ctx).Where("id = ?", id).Select()
	if err != nil {
		log.Err(err)
		err = ErrNotExists
	}

	return apartmentUser, err
}

// UpdateApartmentUserRepository updates apartment user's info
func (apu ApartmentUser) UpdateApartmentUserRepository(c echo.Context, updateUser apartment_user.ApartmentUser) error {
	ctx := c.Request().Context()
	_, err := apu.db.Model(&updateUser).Context(ctx).Where("id = ?", updateUser.ID).Update()
	log.Err(err)
	return err
}

// DeleteApartmentUserRepository sets deleted_at for a apartment user
func (apu ApartmentUser) DeleteApartmentUserRepository(c echo.Context, deleteUser apartment_user.ApartmentUser) error {
	return apu.db.Delete(&deleteUser)
}

func (apu ApartmentUser) AccessTokenRepository(c echo.Context, token apartment_user.AccessToken) (string, error) {
	ctx := c.Request().Context()
	err := apu.db.Model(&token).Context(ctx).Where("token = ?", token.Token).Select()
	if err == nil || err == pg.ErrNoRows {
		_, err := apu.db.Model(&token).Context(ctx).Where("id = ?", 1).Update()
		if err != nil {
			log.Err(err)
		}
	} else {
		log.Err(err)
	}

	return token.Token, err
}
