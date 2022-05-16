package apartment_controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	_ "github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/apartment"
	apartmentController "github.com/ribice/gorsk/pkg/api/apartment/apartment_controller"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/mock/mockdb"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/ribice/gorsk/pkg/utl/server"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateApartment(t *testing.T) {

	apartmentMock := apartment.Apartment{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2022),
			UpdatedAt: mock.TestTime(2022),
		},
		Name:                   "ranjeet",
		Description:            "kolar road",
		FloorArea:              500,
		PricePerMonth:          5000,
		NumberOfRooms:          5,
		GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
		AssociatedRealtor:      "rahul",
	}

	apartmentMockNameEmpty := apartment.Apartment{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2022),
			UpdatedAt: mock.TestTime(2022),
		},
		Name:                   "",
		Description:            "kolar road",
		FloorArea:              500,
		PricePerMonth:          5000,
		NumberOfRooms:          5,
		GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
		AssociatedRealtor:      "rahul",
	}

	request, err := json.Marshal(apartmentMock)
	if err != nil {
		log.Fatal(err)
	}

	requestEmptyName, err := json.Marshal(apartmentMockNameEmpty)
	if err != nil {
		log.Fatal(err)
	}

	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *apartment.Apartment
		udb        *mockdb.Apartment
	}{
		{
			name:       "validation test",
			req:        `{"id":1,"created_at":"2022-05-19T01:02:03.000000004Z","updated_at":"2022-05-19T01:02:03.000000004Z","deleted_at":"0001-01-01T00:00:00Z","created_by":"","updated_by":"","deleted_by":"","is_deleted":false,"name":"ranjeet","description":"kolar road","floor_area":500,"price_per_month":5000,"number_of_rooms":5,"geolocation_coordinates:{"lat":23,"long":56},"associated_realtor":"rahul"}`,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "name empty test",
			req:  string(requestEmptyName),
			udb: &mockdb.Apartment{
				CreateApartmentFn: func(c echo.Context, apartment apartment.Apartment) (apartment.Apartment, error) {
					apartment.ID = 1
					apartment.CreatedAt = mock.TestTime(2022)
					apartment.UpdatedAt = mock.TestTime(2022)
					return apartment, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Success",
			req:  string(request),
			udb: &mockdb.Apartment{
				CreateApartmentFn: func(c echo.Context, apartment apartment.Apartment) (apartment.Apartment, error) {
					apartment.ID = 1
					apartment.CreatedAt = mock.TestTime(2022)
					apartment.UpdatedAt = mock.TestTime(2022)
					return apartment, nil
				},
			},
			wantResp:   &apartmentMock,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentController.NewHTTP(apartment_service.InitializeMock(tt.udb), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartments"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(apartment.Apartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestListApartment(t *testing.T) {

	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *[]apartment.Apartment
		udb        *mockdb.Apartment
	}{
		{
			name:       "Invalid request",
			req:        `?limit=10&offset=-1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing parameter",
			req:        `?limit=10&offset=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			req:  `?limit=10&offset=1&noOfRooms=5&pricePerMonth=5000&floorArea=500`,
			udb: &mockdb.Apartment{
				ListApartmentFn: func(c echo.Context, p apartment.Pagination, filterApartment apartment.FilterApartment) ([]apartment.Apartment, error) {
					if p.Limit == 10 && p.Offset == 1 {
						return []apartment.Apartment{
							{
								Base: model.Base{
									ID:        10,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
									CreatedBy: "ranjeet",
									UpdatedBy: "ranjeet",
								},
								Name:                   "ranjeet",
								Description:            "kolar road",
								FloorArea:              500,
								PricePerMonth:          5000,
								NumberOfRooms:          5,
								GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
								AssociatedRealtor:      "rahul",
							},
							{
								Base: model.Base{
									ID:        11,
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
									CreatedBy: "karan",
									UpdatedBy: "karan",
								},
								Name:                   "karan",
								Description:            "kolar road",
								FloorArea:              500,
								PricePerMonth:          5000,
								NumberOfRooms:          5,
								GeolocationCoordinates: apartment.Vertex{Lat: 25, Long: 52},
								AssociatedRealtor:      "rahul",
							},
						}, nil
					}
					return nil, gorsk.ErrGeneric
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &[]apartment.Apartment{
				{
					Base: model.Base{
						ID:        10,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
						CreatedBy: "ranjeet",
						UpdatedBy: "ranjeet",
					},
					Name:                   "ranjeet",
					Description:            "kolar road",
					FloorArea:              500,
					PricePerMonth:          5000,
					NumberOfRooms:          5,
					GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
					AssociatedRealtor:      "rahul",
				},
				{
					Base: model.Base{
						ID:        11,
						CreatedAt: mock.TestTime(2004),
						UpdatedAt: mock.TestTime(2005),
						CreatedBy: "karan",
						UpdatedBy: "karan",
					},
					Name:                   "karan",
					Description:            "kolar road",
					FloorArea:              500,
					PricePerMonth:          5000,
					NumberOfRooms:          5,
					GeolocationCoordinates: apartment.Vertex{Lat: 25, Long: 52},
					AssociatedRealtor:      "rahul",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentController.NewHTTP(apartment_service.InitializeMock(tt.udb), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartments" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new([]apartment.Apartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestViewApartment(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   apartment.Apartment
		udb        *mockdb.Apartment
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			req:  `1`,
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					return apartment.Apartment{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
							CreatedBy: "karan",
							UpdatedBy: "karan",
						},
						Name:                   "karan",
						Description:            "kolar road",
						FloorArea:              500,
						PricePerMonth:          5000,
						NumberOfRooms:          5,
						GeolocationCoordinates: apartment.Vertex{Lat: 25, Long: 52},
						AssociatedRealtor:      "rahul",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: apartment.Apartment{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
					CreatedBy: "karan",
					UpdatedBy: "karan",
				},
				Name:                   "karan",
				Description:            "kolar road",
				FloorArea:              500,
				PricePerMonth:          5000,
				NumberOfRooms:          5,
				GeolocationCoordinates: apartment.Vertex{Lat: 25, Long: 52},
				AssociatedRealtor:      "rahul",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentController.NewHTTP(apartment_service.InitializeMock(tt.udb), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartments/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp.ID != 0 {
				response := new(apartment.Apartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, &tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdateApartment(t *testing.T) {

	apartmentMock := apartment.Apartment{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2000),
			UpdatedAt: mock.TestTime(2000),
		},
		Name:                   "ranjeet",
		Description:            "kolar road",
		FloorArea:              500,
		PricePerMonth:          5000,
		NumberOfRooms:          5,
		GeolocationCoordinates: apartment.Vertex{Lat: 23, Long: 56},
		AssociatedRealtor:      "rahul",
	}

	request, err := json.Marshal(apartmentMock)
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   apartment.Apartment
		udb        *mockdb.Apartment
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation test",
			id:         `1`,
			req:        `{"name":"ranjeet","description":"colar_road","floor_area":500,"price_per_month":5000,"number_of_rooms":4,"geolocation_coordinates:{"lat":23,"long":56},"associated_realtor":"rahul"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			id:   `1`,
			req:  string(request),
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
					return apartmentMock, nil
				},
				UpdateApartmentFn: func(c echo.Context, apartmentMock apartment.Apartment) error {
					apartmentMock.UpdatedAt = mock.TestTime(2010)
					apartmentMock.PricePerMonth = 5000
					return nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: apartment.Apartment{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
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
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentController.NewHTTP(apartment_service.InitializeMock(tt.udb), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartments/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp.ID != 0 {
				response := new(apartment.Apartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, &tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestDeleteApartment(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.Apartment
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, i int) (apartment.Apartment, error) {
					return apartment.Apartment{}, nil
				},
				DeleteApartmentFn: func(c echo.Context, apartment2 apartment.Apartment) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentController.NewHTTP(apartment_service.InitializeMock(tt.udb), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartments/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
