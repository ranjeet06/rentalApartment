package apartment_platform

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
)

// ApartmentDB represents apartment repository interface
type ApartmentDB interface {
	CreateApartmentRepository(echo.Context, apartment.Apartment) (apartment.Apartment, error)
	ListApartmentRepository(echo.Context, apartment.Pagination, apartment.FilterApartment) ([]apartment.Apartment, error)
	ViewApartmentRepository(echo.Context, int) (apartment.Apartment, error)
	UpdateApartmentRepository(echo.Context, apartment.Apartment) error
	DeleteApartmentRepository(echo.Context, apartment.Apartment) error
}
