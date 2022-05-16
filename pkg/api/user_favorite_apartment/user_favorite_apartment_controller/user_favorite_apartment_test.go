package user_favorite_apartment_controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	FavApartmentController "github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_controller"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_service"
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

func TestHTTP_AddUserFavApartment(t *testing.T) {

	favApartmentMock := user_favorite_apartment.UserFavoriteApartment{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2022),
			UpdatedAt: mock.TestTime(2022),
		},
		UserID:      1,
		ApartmentID: 1,
	}

	request, err := json.Marshal(favApartmentMock)
	if err != nil {
		log.Fatal(err)
	}

	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *user_favorite_apartment.UserFavoriteApartment
		fdb        *mockdb.FavApartment
		udb        *mockdb.Apartment
		UsrDb      *mockdb.ApartmentUser
	}{

		{
			name: "Success",
			req:  string(request),
			fdb: &mockdb.FavApartment{
				AddFavApartmentFn: func(c echo.Context, favApartment user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error) {
					favApartment.ID = 1
					favApartment.CreatedAt = mock.TestTime(2022)
					favApartment.UpdatedAt = mock.TestTime(2022)
					return favApartment, nil
				},
			},
			udb: &mockdb.Apartment{
				ViewApartmentFn: func(c echo.Context, id int) (apartment.Apartment, error) {
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
						GeolocationCoordinates: apartment.Vertex{Lat: 25, Long: 52},
						AssociatedRealtor:      "rahul",
					}, nil
				},
			},
			UsrDb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, i int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:        "karan",
						UserEmail:   "karan@123",
						UserAddress: "danish",
					}, nil
				},
			},
			wantResp:   &favApartmentMock,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			FavApartmentController.NewHTTP(user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/favorites"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(user_favorite_apartment.UserFavoriteApartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func Test_ListUsrFavApartment(t *testing.T) {

	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *[]user_favorite_apartment.UserFavoriteApartment
		fdb        *mockdb.FavApartment
		udb        *mockdb.Apartment
		UsrDb      *mockdb.ApartmentUser
	}{

		{
			name: "Success",
			req:  ``,
			fdb: &mockdb.FavApartment{
				ListFavApartmentFn: func(c echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error) {
					return []user_favorite_apartment.UserFavoriteApartment{
						{
							Base: model.Base{
								ID:        10,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							UserID:      1,
							ApartmentID: 1,
						},
						{
							Base: model.Base{
								ID:        11,
								CreatedAt: mock.TestTime(2004),
								UpdatedAt: mock.TestTime(2005),
							},
							UserID:      2,
							ApartmentID: 2,
						},
					}, nil
					return nil, gorsk.ErrGeneric
				},
			},
			udb: &mockdb.Apartment{
				CreateApartmentFn: func(c echo.Context, ap apartment.Apartment) (apartment.Apartment, error) {
					ap.CreatedAt = mock.TestTime(2000)
					ap.UpdatedAt = mock.TestTime(2000)
					ap.Base.ID = 1
					return ap, nil
				},
			},
			UsrDb: &mockdb.ApartmentUser{
				CreateApartmentUserFn: func(c echo.Context, ap apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
					ap.CreatedAt = mock.TestTime(2000)
					ap.UpdatedAt = mock.TestTime(2000)
					ap.Base.ID = 1
					return ap, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &[]user_favorite_apartment.UserFavoriteApartment{
				{
					Base: model.Base{
						ID:        10,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					UserID:      1,
					ApartmentID: 1,
				},
				{
					Base: model.Base{
						ID:        11,
						CreatedAt: mock.TestTime(2004),
						UpdatedAt: mock.TestTime(2005),
					},
					UserID:      2,
					ApartmentID: 2,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			FavApartmentController.NewHTTP(user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/favorites" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new([]user_favorite_apartment.UserFavoriteApartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func Test_ViewUsrFavApartment(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   []user_favorite_apartment.UserFavoriteApartment
		fdb        *mockdb.FavApartment
		udb        *mockdb.Apartment
		UsrDb      *mockdb.ApartmentUser
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			req:  `1`,
			fdb: &mockdb.FavApartment{
				ViewFavApartmentFn: func(c echo.Context, id int) ([]user_favorite_apartment.UserFavoriteApartment, error) {
					return []user_favorite_apartment.UserFavoriteApartment{
						{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							UserID:      1,
							ApartmentID: 1,
						}, {
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							UserID:      1,
							ApartmentID: 2,
						},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: []user_favorite_apartment.UserFavoriteApartment{
				{
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
					UserID:      1,
					ApartmentID: 1,
				}, {
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
					UserID:      1,
					ApartmentID: 2,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			FavApartmentController.NewHTTP(user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/favorites/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if tt.wantResp != nil {
				response := new([]user_favorite_apartment.UserFavoriteApartment)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				for i, v := range *response {
					assert.Equal(t, tt.wantResp[i].UserID, v.UserID)
					assert.Equal(t, tt.wantResp[i].ApartmentID, v.ApartmentID)
				}

			}

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func Test_DeleteUsrFavApartment(t *testing.T) {
	cases := []struct {
		name       string
		res        string
		wantStatus int
		fdb        *mockdb.FavApartment
		udb        *mockdb.Apartment
		UsrDb      *mockdb.ApartmentUser
	}{
		{
			name:       "Invalid request",
			res:        `?a`,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "Success",
			res:  `?user_id=2&apartment_id=2`,
			udb: &mockdb.Apartment{
				CreateApartmentFn: func(c echo.Context, ap apartment.Apartment) (apartment.Apartment, error) {
					ap.CreatedAt = mock.TestTime(2000)
					ap.UpdatedAt = mock.TestTime(2000)
					ap.Base.ID = 1
					return ap, nil
				},
			},
			UsrDb: &mockdb.ApartmentUser{
				CreateApartmentUserFn: func(c echo.Context, ap apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
					ap.CreatedAt = mock.TestTime(2000)
					ap.UpdatedAt = mock.TestTime(2000)
					ap.Base.ID = 1
					return ap, nil
				},
			},
			fdb: &mockdb.FavApartment{
				DeleteFavApartmentFn: func(c echo.Context, i int, j int) error {
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
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			FavApartmentController.NewHTTP(user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/favorites" + tt.res
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
