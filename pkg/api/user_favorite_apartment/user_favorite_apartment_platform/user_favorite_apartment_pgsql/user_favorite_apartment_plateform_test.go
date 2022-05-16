package user_favorite_apartment_pgsql_test

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_platform/user_favorite_apartment_pgsql"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func Test_AddFavApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      user_favorite_apartment.UserFavoriteApartment
		wantData user_favorite_apartment.UserFavoriteApartment
	}{

		{
			name: "Success",
			req: user_favorite_apartment.UserFavoriteApartment{
				Base: model.Base{
					ID: 2,
				},
				UserID:      1,
				ApartmentID: 1,
			},
			wantData: user_favorite_apartment.UserFavoriteApartment{
				Base: model.Base{
					ID: 2,
				},
				UserID:      1,
				ApartmentID: 1,
			},
		},
		{
			name:    "apartment already exists",
			wantErr: true,
			req: user_favorite_apartment.UserFavoriteApartment{
				ApartmentID: 1,
			},
		},
	}

	db := mock.NewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})

	err := mock.InsertMultiple(db,
		&user_favorite_apartment.UserFavoriteApartment{
			Base: model.Base{
				ID: 1,
			},
			UserID: 1,
		})
	if err != nil {
		t.Error(err)
	}

	udb := user_favorite_apartment_pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.AddFavApartmentRepository(echo.New().NewContext(request, nil), tt.req)
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

	defer mock.DeleteNewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})
}

func Test_ViewFavApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData []user_favorite_apartment.UserFavoriteApartment
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			wantData: []user_favorite_apartment.UserFavoriteApartment{
				{
					UserID:      1,
					ApartmentID: 1,
				},
				{
					UserID:      1,
					ApartmentID: 2,
				},
			},
			id: 1,
		},
	}

	db := mock.NewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := user_favorite_apartment_pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			FavApartment, err := udb.ViewFavApartmentRepository(echo.New().NewContext(request, nil), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range FavApartment {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				for i, v := range FavApartment {
					assert.Equal(t, tt.wantData[i].UserID, v.UserID)
					assert.Equal(t, tt.wantData[i].ApartmentID, v.ApartmentID)
				}
			}
		})
	}

	defer mock.DeleteNewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})
}

func Test_ListFavApartmentRepository(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		wantData []user_favorite_apartment.UserFavoriteApartment
	}{

		{
			name: "Success",
			wantData: []user_favorite_apartment.UserFavoriteApartment{
				{
					Base: model.Base{
						ID: 1,
					},
					UserID:      1,
					ApartmentID: 1,
				},
				{
					Base: model.Base{
						ID: 2,
					},
					UserID:      2,
					ApartmentID: 1,
				},
			},
		},
	}

	db := mock.NewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})

	if err := mock.InsertMultiple(db, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := user_favorite_apartment_pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			FavApartment, err := udb.ListFavApartmentRepository(echo.New().NewContext(request, nil))
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range FavApartment {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, FavApartment)
			}
		})
	}
	defer mock.DeleteNewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})
}

func Test_DeleteFavApartmentRepository(t *testing.T) {
	cases := []struct {
		name         string
		wantErr      bool
		favApartment user_favorite_apartment.UserFavoriteApartment
	}{
		{
			name: "Success",
			favApartment: user_favorite_apartment.UserFavoriteApartment{
				Base: model.Base{
					ID:        2,
					DeletedAt: mock.TestTime(2018),
				},
				UserID:      1,
				ApartmentID: 2,
			},
			wantErr: false,
		},
	}

	db := mock.NewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})

	if err := mock.InsertMultiple(db, &cases[0].favApartment); err != nil {
		t.Error(err)
	}

	udb := user_favorite_apartment_pgsql.Initialize(db)
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.DeleteFavApartmentRepository(echo.New().NewContext(request, nil), tt.favApartment.UserID, tt.favApartment.ApartmentID)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
	defer mock.DeleteNewMockDB(t, &user_favorite_apartment.UserFavoriteApartment{})
}
