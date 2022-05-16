package apartment_service

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
)

// ApartmentService represents apartment application interface
type ApartmentService interface {
	CreateApartmentService(echo.Context, apartment.Apartment) (apartment.Apartment, error)
	ListApartmentService(echo.Context, apartment.Pagination, apartment.FilterApartment) ([]apartment.Apartment, error)
	ViewApartmentService(echo.Context, int) (apartment.Apartment, error)
	DeleteApartmentService(echo.Context, int) error
	UpdateApartmentService(echo.Context, NewUpdateApartment) (apartment.Apartment, error)
}
