package apartment_service_test

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/mock/mockdb"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateApartmentService(t *testing.T) {
	type args struct {
		c   echo.Context
		req apartment.Apartment
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData apartment.Apartment
		udb      *mockdb.Apartment
	}{{
		name: "Success",
		args: args{req: apartment.Apartment{
			Base: model.Base{
				ID: 1,
			},
			Name:                   "karan",
			Description:            "kolar road",
			FloorArea:              500,
			PricePerMonth:          5000,
			NumberOfRooms:          5,
			GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
			AssociatedRealtor:      "rahul",
		}},
		udb: &mockdb.Apartment{
			CreateApartmentFn: func(c echo.Context, ap apartment.Apartment) (apartment.Apartment, error) {
				ap.CreatedAt = mock.TestTime(2000)
				ap.UpdatedAt = mock.TestTime(2000)
				ap.Base.ID = 1
				return ap, nil
			},
		},
		wantData: apartment.Apartment{
			Base: model.Base{
				ID:        1,
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
			},
			Name:                   "karan",
			Description:            "kolar road",
			FloorArea:              500,
			PricePerMonth:          5000,
			NumberOfRooms:          5,
			GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
			AssociatedRealtor:      "rahul",
		}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_service.InitializeMock(tt.udb)
			apartment, err := s.CreateApartmentService(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, apartment)
		})
	}
}

func Test_ListApartmentService(t *testing.T) {
	type args struct {
		c      echo.Context
		pgn    apartment.Pagination
		filter apartment.FilterApartment
	}
	cases := []struct {
		name     string
		args     args
		wantData []apartment.Apartment
		wantErr  bool
		udb      *mockdb.Apartment
	}{
		{
			name: "Success",
			args: args{c: nil, pgn: apartment.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			udb: &mockdb.Apartment{
				ListApartmentFn: func(echo.Context, apartment.Pagination, apartment.FilterApartment) ([]apartment.Apartment, error) {
					return []apartment.Apartment{
						{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
							Name:                   "karan",
							Description:            "kolar road",
							FloorArea:              500,
							PricePerMonth:          5000,
							NumberOfRooms:          5,
							GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
							AssociatedRealtor:      "rahul",
						},
						{
							Base: model.Base{
								ID:        2,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							Name:                   "rahul",
							Description:            "kolar road",
							FloorArea:              500,
							PricePerMonth:          5000,
							NumberOfRooms:          5,
							GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
							AssociatedRealtor:      "rahul",
						},
					}, nil
				}},
			wantData: []apartment.Apartment{
				{
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
					Name:                   "karan",
					Description:            "kolar road",
					FloorArea:              500,
					PricePerMonth:          5000,
					NumberOfRooms:          5,
					GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
					AssociatedRealtor:      "rahul",
				},
				{
					Base: model.Base{
						ID:        2,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					Name:                   "rahul",
					Description:            "kolar road",
					FloorArea:              500,
					PricePerMonth:          5000,
					NumberOfRooms:          5,
					GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
					AssociatedRealtor:      "rahul",
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_service.InitializeMock(tt.udb)
			apartments, err := s.ListApartmentService(tt.args.c, tt.args.pgn, tt.args.filter)
			assert.Equal(t, tt.wantData, apartments)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func Test_ViewApartmentService(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name     string
		args     args
		wantData apartment.Apartment
		wantErr  error
		udb      *mockdb.Apartment
	}{
		{
			name: "Success",
			args: args{id: 1},
			wantData: apartment.Apartment{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					if id == 1 {
						return apartment.Apartment{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							Name:                   "karan",
							Description:            "kolar road",
							FloorArea:              500,
							PricePerMonth:          5000,
							NumberOfRooms:          5,
							GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
							AssociatedRealtor:      "rahul",
						}, nil
					}
					return apartment.Apartment{}, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_service.InitializeMock(tt.udb)
			apartment, err := s.ViewApartmentService(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, apartment)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_UpdateApartmentService(t *testing.T) {
	type args struct {
		c   echo.Context
		upd apartment_service.NewUpdateApartment
	}
	cases := []struct {
		name     string
		args     args
		wantData apartment.Apartment
		wantErr  error
		udb      *mockdb.Apartment
	}{

		{
			name: "Fail on Update",
			args: args{upd: apartment_service.NewUpdateApartment{
				ID: 1,
			}},
			wantErr: gorsk.ErrGeneric,
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					return apartment.Apartment{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(1991),
						},
						Name:                   "karan",
						Description:            "kolar road",
						FloorArea:              500,
						PricePerMonth:          5000,
						NumberOfRooms:          5,
						GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
						AssociatedRealtor:      "rahul",
					}, nil
				},
				UpdateApartmentFn: func(c echo.Context, apartment apartment.Apartment) error {
					return gorsk.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{upd: apartment_service.NewUpdateApartment{
				ID:                1,
				Name:              "mohit",
				Description:       "danish road",
				FloorArea:         500,
				PricePerMonth:     5000,
				NumberOfRooms:     4,
				AssociatedRealtor: "rahul",
			}},
			wantData: apartment.Apartment{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(1990),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:                   "mohit",
				Description:            "danish road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          4,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					return apartment.Apartment{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:                   "mohit",
						Description:            "danish road",
						FloorArea:              500,
						PricePerMonth:          5000,
						NumberOfRooms:          4,
						GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
						AssociatedRealtor:      "rahul",
					}, nil
				},
				UpdateApartmentFn: func(c echo.Context, apartment apartment.Apartment) error {
					apartment.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_service.InitializeMock(tt.udb)
			usr, err := s.UpdateApartmentService(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_DeleteApartmentService(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
		udb     *mockdb.Apartment
	}{
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: gorsk.ErrGeneric,
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					if id != 1 {
						return apartment.Apartment{}, nil
					}
					return apartment.Apartment{}, gorsk.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					return apartment.Apartment{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:                   "karan",
						Description:            "kolar road",
						FloorArea:              500,
						PricePerMonth:          5000,
						NumberOfRooms:          5,
						GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
						AssociatedRealtor:      "rahul",
					}, nil
				},
				DeleteApartmentFn: func(c echo.Context, apartment apartment.Apartment) error {
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_service.InitializeMock(tt.udb)
			err := s.DeleteApartmentService(tt.args.c, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}
