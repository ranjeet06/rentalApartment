package pgsql

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

// Apartment represents for apartment repository
type Apartment struct {
	db    orm.DB
	newDb *gorm.DB
}

// Initialize initializes apartment repository with default
func Initialize(db orm.DB, newDb *gorm.DB) *Apartment {
	return &Apartment{db: db, newDb: newDb}
}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusConflict, "name already exist.")
)

// CreateApartmentRepository creates a new apartment on database
func (ap Apartment) CreateApartmentRepository(c echo.Context, apartmentNew apartment.Apartment) (apartment.Apartment, error) {
	var newApartment = new(apartment.Apartment)

	ctx := c.Request().Context()
	err := ap.db.Model(newApartment).Context(ctx).Where("name = ?", apartmentNew.Name).Select()
	if err != nil {
		log.Err(err)
	}
	if err == nil || err != pg.ErrNoRows {
		return apartment.Apartment{}, ErrAlreadyExists
	}

	err = ap.db.Insert(&apartmentNew)
	log.Err(err)
	return apartmentNew, err
}

// ListApartmentRepository returns list of all apartments.
func (ap Apartment) ListApartmentRepository(c echo.Context, pagination apartment.Pagination, filterApartment apartment.FilterApartment) ([]apartment.Apartment, error) {
	var apartments []apartment.Apartment
	ctx := c.Request().Context()
	q := ap.db.Model(&apartments).Context(ctx).Limit(pagination.Limit).Offset(pagination.Offset).Where("deleted_at is null")

	if filterApartment.NumberOfRooms != 0 {
		q = q.Where("number_of_rooms = ?", filterApartment.NumberOfRooms)
	}

	if filterApartment.PricePerMonth != 0 {
		q = q.Where("price_per_month = ?", filterApartment.PricePerMonth)
	}

	if filterApartment.FloorArea != 0 {
		q = q.Where("floor_area = ?", filterApartment.FloorArea)
	}

	err := q.Select()
	log.Err(err)
	return apartments, err
}

// Custom errors
var (
	ErrNotExists = echo.NewHTTPError(http.StatusBadRequest, "Apartment Id not exist.")
)

// ViewApartmentRepository returns single apartment by ID
func (ap Apartment) ViewApartmentRepository(c echo.Context, id int) (apartment.Apartment, error) {
	var viewApartment apartment.Apartment
	ctx := c.Request().Context()
	q := ap.db.Model(&viewApartment).Context(ctx).Where("id = ?", id)
	err := q.Select()
	if err != nil {
		log.Err(err)
		err = ErrNotExists
	}

	return viewApartment, err
}

// UpdateApartmentRepository updates apartment's info
func (ap Apartment) UpdateApartmentRepository(c echo.Context, apartment apartment.Apartment) error {
	ctx := c.Request().Context()
	_, err := ap.db.Model(&apartment).Context(ctx).Where(`"id" = ?`, apartment.ID).Update()
	log.Err(err)
	return err
}

// DeleteApartmentRepository sets deleted_at for a apartment
func (ap Apartment) DeleteApartmentRepository(c echo.Context, apartment apartment.Apartment) error {
	ctx := c.Request().Context()

	tx := ap.newDb.Begin(nil).WithContext(ctx)

	err := tx.Select(`"id = ?"`, apartment.ID).Error
	if err != nil {
		tx.Rollback()
		log.Err(err)
		return err
	}

	err = tx.Delete(&apartment, apartment.ID).Error
	if err != nil {
		tx.Rollback()
		log.Err(err)
		return err
	}

	err = tx.Commit().Error
	log.Err(err)
	return err
}
