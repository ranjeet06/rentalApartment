package pgsql_test

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_platform/pgsql"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func Test_CreateApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      apartment.Apartment
		wantData apartment.Apartment
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: apartment.Apartment{
				Base: model.Base{
					ID: 1,
				},
				Name:                   "ranjeet",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
		},
		{
			name: "Success",
			req: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
			wantData: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
		},
		{
			name:    "User already exists",
			wantErr: true,
			req: apartment.Apartment{
				Name: "karan",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment.Apartment{})

	err := mock.InsertMultiple(db,
		&apartment.Apartment{
			Base: model.Base{
				ID: 1,
			},
			Name: "neha",
		})
	if err != nil {
		t.Error(err)
	}

	udb := pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.CreateApartmentRepository(echo.New().NewContext(request, nil), tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				if resp.ID == 0 {
					t.Error("expected data, but got empty struct.")
					return
				}
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.UpdatedAt = resp.UpdatedAt
				assert.Equal(t, tt.wantData, resp)
			}

		})
	}

	defer mock.DeleteNewMockDB(t, &apartment.Apartment{})
}

func Test_ViewApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData apartment.Apartment
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment.Apartment{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartment, err := udb.ViewApartmentRepository(echo.New().NewContext(request, nil), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				if apartment.ID == 0 {
					t.Errorf("response was empty due to: %v", err)
				} else {
					tt.wantData.CreatedAt = apartment.CreatedAt
					tt.wantData.UpdatedAt = apartment.UpdatedAt
					assert.Equal(t, tt.wantData, apartment)
				}
			}
		})
	}

	defer mock.DeleteNewMockDB(t, &apartment.Apartment{})
}

func Test_ListApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		pg       apartment.Pagination
		filter   apartment.FilterApartment
		wantData []apartment.Apartment
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: apartment.Pagination{
				Limit: -100,
			},
			filter: apartment.FilterApartment{
				NumberOfRooms: 5,
				FloorArea:     500,
				PricePerMonth: 5000,
			},
		},
		{
			name: "Success",
			pg: apartment.Pagination{
				Limit:  10,
				Offset: 0,
			},
			filter: apartment.FilterApartment{
				NumberOfRooms: 5,
				PricePerMonth: 5000,
				FloorArea:     500,
			},
			wantData: []apartment.Apartment{
				{
					Base: model.Base{
						ID: 2,
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
						ID: 3,
					},
					Name:                   "rahul",
					Description:            "kolar road",
					FloorArea:              500,
					PricePerMonth:          5000,
					NumberOfRooms:          5,
					GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
					AssociatedRealtor:      "rahul",
				},
			},
		},
	}

	db := mock.NewMockDB(t, &apartment.Apartment{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			newApartments, err := udb.ListApartmentRepository(echo.New().NewContext(request, nil), tt.pg, tt.filter)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range newApartments {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, newApartments)
			}
		})
	}
	defer mock.DeleteNewMockDB(t, &apartment.Apartment{})
}

func Test_UpdateApartmentRepository(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   bool
		apartment apartment.Apartment
		wantData  apartment.Apartment
	}{
		{
			name: "Success",
			apartment: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
			wantData: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment.Apartment{})

	if err := mock.InsertMultiple(db, &cases[0].apartment); err != nil {
		t.Error(err)
	}

	udb := pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.UpdateApartmentRepository(echo.New().NewContext(request, nil), tt.wantData)
			if tt.wantErr != (err != nil) {
				fmt.Println(tt.wantErr, err)
			}
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				apartment := apartment.Apartment{
					Base: model.Base{
						ID: tt.apartment.ID,
					},
				}
				if err := db.Select(&apartment); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = apartment.UpdatedAt
				tt.wantData.CreatedAt = apartment.CreatedAt
				tt.wantData.DeletedAt = apartment.DeletedAt
				assert.Equal(t, tt.wantData, apartment)
			}
		})
	}
	defer mock.DeleteNewMockDB(t, &apartment.Apartment{})
}

func Test_DeleteApartmentRepository(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   bool
		apartment apartment.Apartment
		wantData  apartment.Apartment
	}{
		{
			name: "Success",
			apartment: apartment.Apartment{
				Base: model.Base{
					ID:        2,
					DeletedAt: mock.TestTime(2018),
				},
			},
			wantData: apartment.Apartment{
				Base: model.Base{
					ID: 2,
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
				AssociatedRealtor:      "rahul",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment.Apartment{})

	if err := mock.InsertMultiple(db, &cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.DeleteApartmentRepository(echo.New().NewContext(request, nil), tt.apartment)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
	defer mock.DeleteNewMockDB(t, &apartment.Apartment{})
}
