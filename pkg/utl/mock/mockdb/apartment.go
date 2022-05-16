package mockdb

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
)

// Apartment database mock
type Apartment struct {
	CreateApartmentFn func(echo.Context, apartment.Apartment) (apartment.Apartment, error)
	ListApartmentFn   func(echo.Context, apartment.Pagination, apartment.FilterApartment) ([]apartment.Apartment, error)
	ViewApartmentFn   func(echo.Context, int) (apartment.Apartment, error)
	UpdateApartmentFn func(echo.Context, apartment.Apartment) error
	DeleteApartmentFn func(echo.Context, apartment.Apartment) error
}

// CreateApartmentRepository mock
func (ap *Apartment) CreateApartmentRepository(c echo.Context, apartmentNew apartment.Apartment) (apartment.Apartment, error) {
	return ap.CreateApartmentFn(c, apartmentNew)
}

// ListApartmentRepository mock
func (ap Apartment) ListApartmentRepository(c echo.Context, pagination apartment.Pagination, filterApartment apartment.FilterApartment) ([]apartment.Apartment, error) {
	return ap.ListApartmentFn(c, pagination, filterApartment)
}

// ViewApartmentRepository mock
func (ap Apartment) ViewApartmentRepository(c echo.Context, id int) (apartment.Apartment, error) {
	return ap.ViewApartmentFn(c, id)
}

// UpdateApartmentRepository mock
func (ap Apartment) UpdateApartmentRepository(c echo.Context, apartment apartment.Apartment) error {
	return ap.UpdateApartmentFn(c, apartment)
}

// DeleteApartmentRepository mock
func (ap Apartment) DeleteApartmentRepository(c echo.Context, apartment apartment.Apartment) error {
	return ap.DeleteApartmentFn(c, apartment)
}
