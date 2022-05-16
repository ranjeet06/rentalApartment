package user_pgsql_test

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_platform/user_pgsql"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func Test_CreateApartmentUserRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      apartment_user.ApartmentUser
		wantData apartment_user.ApartmentUser
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 1,
				},
				Name:        "ranjeet",
				UserEmail:   "ranjeet@123",
				UserAddress: "kolar",
			},
		},
		{
			name: "Success",
			req: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "ranjeet",
				UserEmail:   "ranjeet@123",
				UserAddress: "kolar",
			},
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "ranjeet",
				UserEmail:   "ranjeet@123",
				UserAddress: "kolar",
			},
		},
		{
			name:    "User already exists",
			wantErr: true,
			req: apartment_user.ApartmentUser{
				Name: "karan",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment_user.ApartmentUser{})

	err := mock.InsertMultiple(db,
		&apartment_user.ApartmentUser{
			Base: model.Base{
				ID: 1,
			},
			Name: "neha",
		})
	if err != nil {
		t.Error(err)
	}

	udb := user_pgsql.InitializeUser(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.CreateApartmentUserRepository(echo.New().NewContext(request, nil), tt.req)
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

	defer mock.DeleteNewMockDB(t, &apartment_user.ApartmentUser{})
}

func Test_ViewApartmentUserRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData apartment_user.ApartmentUser
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "karan",
				UserAddress: "kolar",
				UserEmail:   "karan@123",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment_user.ApartmentUser{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := user_pgsql.InitializeUser(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartmentUsr, err := udb.ViewApartmentUserRepository(echo.New().NewContext(request, nil), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				if apartmentUsr.ID == 0 {
					t.Errorf("response was empty due to: %v", err)
				} else {
					tt.wantData.CreatedAt = apartmentUsr.CreatedAt
					tt.wantData.UpdatedAt = apartmentUsr.UpdatedAt
					assert.Equal(t, tt.wantData, apartmentUsr)
				}
			}
		})
	}

	defer mock.DeleteNewMockDB(t, &apartment_user.ApartmentUser{})
}

func Test_ListApartmentUserRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		filter   apartment_user.UserFilter
		wantData []apartment_user.ApartmentUser
	}{
		{
			name: "fail on filter",
			filter: apartment_user.UserFilter{
				Name:      "mohit",
				UserEmail: "",
			},
			wantErr: true,
		},
		{
			name: "Success",

			filter: apartment_user.UserFilter{
				Name:      "",
				UserEmail: "",
			},
			wantData: []apartment_user.ApartmentUser{
				{
					Base: model.Base{
						ID: 2,
					},
					Name:        "karan",
					UserEmail:   "karan@123",
					UserAddress: "kolar",
				},
				{
					Base: model.Base{
						ID: 3,
					},
					Name:        "rahul",
					UserEmail:   "rahul@123",
					UserAddress: "kolar",
				},
			},
		},
	}

	db := mock.NewMockDB(t, &apartment_user.ApartmentUser{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := user_pgsql.InitializeUser(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			newApartmentUsers, err := udb.ListApartmentUserRepository(echo.New().NewContext(request, nil), tt.filter)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range newApartmentUsers {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, newApartmentUsers)
			}
		})
	}
	defer mock.DeleteNewMockDB(t, &apartment_user.ApartmentUser{})
}

func Test_UpdateApartmentRepository(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   bool
		apartment apartment_user.ApartmentUser
		wantData  apartment_user.ApartmentUser
	}{
		{
			name: "Success",
			apartment: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "karan",
				UserEmail:   "karan",
				UserAddress: "kolar",
			},
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "karan",
				UserEmail:   "karan",
				UserAddress: "kolar",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment_user.ApartmentUser{})

	if err := mock.InsertMultiple(db, &cases[0].apartment); err != nil {
		t.Error(err)
	}

	udb := user_pgsql.InitializeUser(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.UpdateApartmentUserRepository(echo.New().NewContext(request, nil), tt.wantData)
			if tt.wantErr != (err != nil) {
				fmt.Println(tt.wantErr, err)
			}
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				apartmentUser := apartment_user.ApartmentUser{
					Base: model.Base{
						ID: tt.apartment.ID,
					},
				}
				if err := db.Select(&apartmentUser); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = apartmentUser.UpdatedAt
				tt.wantData.CreatedAt = apartmentUser.CreatedAt
				tt.wantData.DeletedAt = apartmentUser.DeletedAt
				assert.Equal(t, tt.wantData, apartmentUser)
			}
		})
	}
	defer mock.DeleteNewMockDB(t, &apartment_user.ApartmentUser{})
}

func Test_DeleteApartmentRepository(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   bool
		apartment apartment_user.ApartmentUser
		wantData  apartment_user.ApartmentUser
	}{
		{
			name: "Success",
			apartment: apartment_user.ApartmentUser{
				Base: model.Base{
					ID:        2,
					DeletedAt: mock.TestTime(2018),
				},
			},
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID: 2,
				},
				Name:        "karan",
				UserEmail:   "karan@123",
				UserAddress: "kolar",
			},
		},
	}

	db := mock.NewMockDB(t, &apartment_user.ApartmentUser{})

	if err := mock.InsertMultiple(db, &cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := user_pgsql.InitializeUser(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.DeleteApartmentUserRepository(echo.New().NewContext(request, nil), tt.apartment)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
	defer mock.DeleteNewMockDB(t, &apartment_user.ApartmentUser{})
}
