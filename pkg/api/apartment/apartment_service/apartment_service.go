// Package apartment_service contains apartment application services
package apartment_service

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_platform"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_platform/pgsql"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Apartment represents apartment application service
type Apartment struct {
	db          *pg.DB
	apartmentDb apartment_platform.ApartmentDB
}

// Initialize initalizes Apartment application service with defaults

func Initialize(db *pg.DB, newDb *gorm.DB) *Apartment {
	apartmentDb := pgsql.Initialize(db, newDb)
	return &Apartment{db: db, apartmentDb: apartmentDb}
}

func InitializeMock(db apartment_platform.ApartmentDB) *Apartment {
	return &Apartment{db: nil, apartmentDb: db}
}

// CreateApartmentService creates a new apartment
func (a Apartment) CreateApartmentService(c echo.Context, apartment apartment.Apartment) (apartment.Apartment, error) {
	return a.apartmentDb.CreateApartmentRepository(c, apartment)
}

// ListApartmentService returns list of apartment
func (a Apartment) ListApartmentService(c echo.Context, pagination apartment.Pagination, filterApartment apartment.FilterApartment) ([]apartment.Apartment, error) {
	return a.apartmentDb.ListApartmentRepository(c, pagination, filterApartment)
}

// ViewApartmentService returns single apartment
func (a Apartment) ViewApartmentService(c echo.Context, id int) (apartment.Apartment, error) {
	return a.apartmentDb.ViewApartmentRepository(c, id)
}

// NewUpdateApartment contains apartment's information used for updating
type NewUpdateApartment struct {
	ID                int
	Name              string
	Description       string
	FloorArea         float64
	PricePerMonth     float64
	NumberOfRooms     int
	AssociatedRealtor string
}

// UpdateApartmentService updates apartment's information
func (a Apartment) UpdateApartmentService(c echo.Context, newApartment NewUpdateApartment) (apartment.Apartment, error) {

	if err := a.apartmentDb.UpdateApartmentRepository(c, apartment.Apartment{
		Base:              model.Base{ID: newApartment.ID},
		Name:              newApartment.Name,
		Description:       newApartment.Description,
		FloorArea:         newApartment.FloorArea,
		PricePerMonth:     newApartment.PricePerMonth,
		NumberOfRooms:     newApartment.NumberOfRooms,
		AssociatedRealtor: newApartment.AssociatedRealtor,
	}); err != nil {
		log.Err(err)
		return apartment.Apartment{}, err
	}

	return a.apartmentDb.ViewApartmentRepository(c, newApartment.ID)
}

// DeleteApartmentService deletes a apartment
func (a Apartment) DeleteApartmentService(c echo.Context, b int) error {

	apartment, err := a.apartmentDb.ViewApartmentRepository(c, b)
	if err != nil {
		log.Err(err)
		return err
	}

	return a.apartmentDb.DeleteApartmentRepository(c, apartment)
}
