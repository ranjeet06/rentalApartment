package user_favorite_apartment_service_test

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_service"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/mock/mockdb"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func Test_AddUserFavApartmentServiceApartmentService(t *testing.T) {
	type args struct {
		c   echo.Context
		req user_favorite_apartment.UserFavoriteApartment
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData user_favorite_apartment.UserFavoriteApartment
		fdb      *mockdb.FavApartment
		udb      *mockdb.Apartment
		UsrDb    *mockdb.ApartmentUser
	}{{
		name: "Success",
		args: args{req: user_favorite_apartment.UserFavoriteApartment{
			Base: model.Base{
				ID: 1,
			},
			UserID:      1,
			ApartmentID: 1,
		}},
		fdb: &mockdb.FavApartment{
			AddFavApartmentFn: func(c echo.Context, fap user_favorite_apartment.UserFavoriteApartment) (user_favorite_apartment.UserFavoriteApartment, error) {
				fap.CreatedAt = mock.TestTime(2000)
				fap.UpdatedAt = mock.TestTime(2000)
				fap.Base.ID = 1
				return fap, nil
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
		wantData: user_favorite_apartment.UserFavoriteApartment{
			Base: model.Base{
				ID:        1,
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
			},
			UserID:      1,
			ApartmentID: 1,
		}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			s := user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService)
			apartment, err := s.AddUserFavApartmentService(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, apartment)
		})
	}
}

func Test_ListUserFavApartmentService(t *testing.T) {

	cases := []struct {
		name     string
		wantData []user_favorite_apartment.UserFavoriteApartment
		wantErr  bool
		fdb      *mockdb.FavApartment
		udb      *mockdb.Apartment
		UsrDb    *mockdb.ApartmentUser
	}{
		{
			name: "Success",
			fdb: &mockdb.FavApartment{
				ListFavApartmentFn: func(echo.Context) ([]user_favorite_apartment.UserFavoriteApartment, error) {
					return []user_favorite_apartment.UserFavoriteApartment{
						{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
							UserID:      1,
							ApartmentID: 1,
						},
						{
							Base: model.Base{
								ID:        2,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							UserID:      1,
							ApartmentID: 2,
						},
					}, nil
				}},
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
			wantData: []user_favorite_apartment.UserFavoriteApartment{
				{
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
					UserID:      1,
					ApartmentID: 1,
				},
				{
					Base: model.Base{
						ID:        2,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					UserID:      1,
					ApartmentID: 2,
				}},
		},
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			s := user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService)
			favApartments, err := s.ListUserFavApartmentService(echo.New().NewContext(request, nil))
			assert.Equal(t, tt.wantData, favApartments)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func Test_ViewUserFavApartmentService(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}

	cases := []struct {
		name     string
		wantData []user_favorite_apartment.UserFavoriteApartment
		args     args
		wantErr  error
		fdb      *mockdb.FavApartment
		udb      *mockdb.Apartment
		UsrDb    *mockdb.ApartmentUser
	}{
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
				}},
			args: args{
				id: 1,
			},
			fdb: &mockdb.FavApartment{
				ViewFavApartmentFn: func(c echo.Context, id int) ([]user_favorite_apartment.UserFavoriteApartment, error) {
					if id == 1 {
						return []user_favorite_apartment.UserFavoriteApartment{
							{
								UserID:      1,
								ApartmentID: 1,
							},
							{
								UserID:      1,
								ApartmentID: 2,
							},
						}, nil
					}
					return []user_favorite_apartment.UserFavoriteApartment{}, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			s := user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService)
			FavApartment, err := s.ViewUserFavApartmentService(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, FavApartment)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_DeleteUserFavApartmentService(t *testing.T) {
	type args struct {
		c           echo.Context
		usrId       int
		apartmentId int
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
		fdb     *mockdb.FavApartment
		udb     *mockdb.Apartment
		UsrDb   *mockdb.ApartmentUser
	}{

		{
			name: "Success",
			args: args{usrId: 1, apartmentId: 1},
			fdb: &mockdb.FavApartment{
				DeleteFavApartmentFn: func(c echo.Context, usrId int, apartmentId int) error {
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			apartmentService := apartment_service.InitializeMock(tt.udb)
			UsrService := apartment_user_service.InitializeMockUser(tt.UsrDb)
			s := user_favorite_apartment_service.InitializeMock(tt.fdb, apartmentService, UsrService)
			err := s.DeleteUserFavApartmentService(tt.args.c, tt.args.usrId, tt.args.apartmentId)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}
