package apartment_user_controller_test

import (
	"encoding/json"
	"github.com/labstack/echo"
	_ "github.com/ribice/gorsk"
	apartmentUserAuthMw "github.com/ribice/gorsk/pkg/apartmentUserUtl/middleware/auth"
	_ "github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	apartmentUserController "github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_controller"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/mock/mockdb"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/ribice/gorsk/pkg/utl/server"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateApartmentUser(t *testing.T) {

	apartmentUserMock := apartment_user.ApartmentUser{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2022),
			UpdatedAt: mock.TestTime(2022),
		},
		Name:        "ranjeet",
		UserEmail:   "ranjeet@123",
		UserAddress: "danish",
	}

	apartmentUserMockNameEmpty := apartment_user.ApartmentUser{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2022),
			UpdatedAt: mock.TestTime(2022),
		},
		Name:        "",
		UserEmail:   "ranjeet@123",
		UserAddress: "danish",
	}

	request, err := json.Marshal(apartmentUserMock)
	if err != nil {
		log.Fatal(err)
	}

	requestEmptyName, err := json.Marshal(apartmentUserMockNameEmpty)
	if err != nil {
		log.Fatal(err)
	}
	jwt, err := mock.JwtNew("HS384", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 30, 64)
	if err != nil {
		log.Fatal(err)
	}
	Token, err := jwt.GetDummyToken()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *apartment_user.ApartmentUser
		udb        *mockdb.ApartmentUser
		jwtToken   mock.JWTUser
	}{
		{
			name:       "validation test",
			req:        `string(request) + "}"`,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "name empty test",
			req:  string(requestEmptyName),
			udb: &mockdb.ApartmentUser{
				CreateApartmentUserFn: func(c echo.Context, apartmentUser apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
					apartmentUser.ID = 1
					apartmentUser.CreatedAt = mock.TestTime(2022)
					apartmentUser.UpdatedAt = mock.TestTime(2022)
					return apartmentUser, nil
				},
			},

			wantStatus: http.StatusOK,
		},

		{
			name: "Success",
			req:  string(request),
			udb: &mockdb.ApartmentUser{
				CreateApartmentUserFn: func(c echo.Context, apartmentUser apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
					apartmentUser.ID = 1
					apartmentUser.CreatedAt = mock.TestTime(2022)
					apartmentUser.UpdatedAt = mock.TestTime(2022)
					apartmentUser.Name = "ranjeet"
					apartmentUser.UserEmail = "ranjeet@123"
					apartmentUser.UserAddress = "danish"
					return apartmentUser, nil
				},
			},

			wantResp:   &apartmentUserMock,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			middleware := apartmentUserAuthMw.Middleware(jwt)
			apartmentUserController.NewHTTP(apartment_user_service.InitializeMockUser(tt.udb), rg, middleware)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartment_users"
			req, err := http.NewRequest("POST", path, strings.NewReader(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			a := "Bearer "
			tokenNew := a + Token
			req.Header.Add("Authorization", tokenNew)
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(apartment_user.ApartmentUser)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestListApartmentUser(t *testing.T) {

	jwt, err := mock.JwtNew("HS384", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 30, 64)
	if err != nil {
		log.Fatal(err)
	}
	Token, err := jwt.GetDummyToken()
	if err != nil {
		log.Fatal(err)
	}

	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *[]apartment_user.ApartmentUser
		udb        *mockdb.ApartmentUser
	}{
		{
			name: "Success",
			req:  ``,
			udb: &mockdb.ApartmentUser{
				ListApartmentUserFn: func(c echo.Context, filter apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error) {
					return []apartment_user.ApartmentUser{
						{
							Base: model.Base{
								ID:        10,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
								CreatedBy: "ranjeet",
								UpdatedBy: "ranjeet",
							},
							Name:        "ranjeet",
							UserEmail:   "ranjeet@123",
							UserAddress: "bhopal",
						},
						{
							Base: model.Base{
								ID:        11,
								CreatedAt: mock.TestTime(2004),
								UpdatedAt: mock.TestTime(2005),
								CreatedBy: "karan",
								UpdatedBy: "karan",
							},
							Name:        "karan",
							UserEmail:   "karan@123",
							UserAddress: "bhopal",
						},
					}, nil

				},
			},
			wantStatus: http.StatusOK,
			wantResp: &[]apartment_user.ApartmentUser{
				{
					Base: model.Base{
						ID:        10,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
						CreatedBy: "ranjeet",
						UpdatedBy: "ranjeet",
					},
					Name:        "ranjeet",
					UserEmail:   "ranjeet@123",
					UserAddress: "bhopal",
				},
				{
					Base: model.Base{
						ID:        11,
						CreatedAt: mock.TestTime(2004),
						UpdatedAt: mock.TestTime(2005),
						CreatedBy: "karan",
						UpdatedBy: "karan",
					},
					Name:        "karan",
					UserEmail:   "karan@123",
					UserAddress: "bhopal",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			middleware := apartmentUserAuthMw.Middleware(jwt)
			apartmentUserController.NewHTTP(apartment_user_service.InitializeMockUser(tt.udb), rg, middleware)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartment_users" + tt.req
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			a := "Bearer "
			tokenNew := a + Token
			req.Header.Add("Authorization", tokenNew)
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new([]apartment_user.ApartmentUser)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestViewApartmentUser(t *testing.T) {

	jwt, err := mock.JwtNew("HS384", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 30, 64)
	if err != nil {
		log.Fatal(err)
	}
	Token, err := jwt.GetDummyToken()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   apartment_user.ApartmentUser
		udb        *mockdb.ApartmentUser
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			req:  `1`,
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:        "karan",
						UserEmail:   "karan@123",
						UserAddress: "bhopal",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: apartment_user.ApartmentUser{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:        "karan",
				UserEmail:   "karan@123",
				UserAddress: "bhopal",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			middleware := apartmentUserAuthMw.Middleware(jwt)
			apartmentUserController.NewHTTP(apartment_user_service.InitializeMockUser(tt.udb), rg, middleware)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartment_users/" + tt.req
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			a := "Bearer "
			tokenNew := a + Token
			req.Header.Add("Authorization", tokenNew)
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp.ID != 0 {
				response := new(apartment_user.ApartmentUser)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, &tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdateApartmentUser(t *testing.T) {

	jwt, err := mock.JwtNew("HS384", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 30, 64)
	if err != nil {
		log.Fatal(err)
	}
	Token, err := jwt.GetDummyToken()
	if err != nil {
		log.Fatal(err)
	}

	apartmentUsrMock := apartment_user.ApartmentUser{
		Base: model.Base{
			ID:        1,
			CreatedAt: mock.TestTime(2000),
			UpdatedAt: mock.TestTime(2000),
		},
		Name:        "ranjeet",
		UserEmail:   "ranjeet@123",
		UserAddress: "bhopal",
	}

	request, err := json.Marshal(apartmentUsrMock)
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   apartment_user.ApartmentUser
		udb        *mockdb.ApartmentUser
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation test",
			id:         `1`,
			req:        `{"name:"ranjeet","description":"colar_road","floor_area":500,"price_per_month":5000,"number_of_rooms":4,"geolocation_coordinates:{"lat":23,"long":56},"associated_realtor":"rahul"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			id:   `1`,
			req:  string(request),
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					return apartmentUsrMock, nil
				},
				UpdateApartmentUserFn: func(c echo.Context, apartmentUsrMock apartment_user.ApartmentUser) error {
					apartmentUsrMock.UpdatedAt = mock.TestTime(2010)
					apartmentUsrMock.UserAddress = "5000"
					return nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: apartment_user.ApartmentUser{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:        "ranjeet",
				UserEmail:   "ranjeet@123",
				UserAddress: "bhopal",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			middleware := apartmentUserAuthMw.Middleware(jwt)
			apartmentUserController.NewHTTP(apartment_user_service.InitializeMockUser(tt.udb), rg, middleware)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/apartment_users/" + tt.id
			req, err := http.NewRequest("PATCH", path, strings.NewReader(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			a := "Bearer "
			tokenNew := a + Token
			req.Header.Add("Authorization", tokenNew)
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp.ID != 0 {
				response := new(apartment_user.ApartmentUser)
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

	jwt, err := mock.JwtNew("HS384", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 30, 64)
	if err != nil {
		log.Fatal(err)
	}
	Token, err := jwt.GetDummyToken()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.ApartmentUser
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, i int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{}, nil
				},
				DeleteApartmentUserFn: func(c echo.Context, apartmentUsr apartment_user.ApartmentUser) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			middleware := apartmentUserAuthMw.Middleware(jwt)
			apartmentUserController.NewHTTP(apartment_user_service.InitializeMockUser(tt.udb), rg, middleware)
			ts := httptest.NewServer(r)
			path := ts.URL + "/apartment_users/" + tt.id
			req, err := http.NewRequest("DELETE", path, strings.NewReader(tt.id))
			if err != nil {
				t.Fatal(err)
			}
			a := "Bearer "
			tokenNew := a + Token
			req.Header.Add("Authorization", tokenNew)
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
