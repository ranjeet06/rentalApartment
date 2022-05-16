package apartment_user_service_test

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/utl/mock"
	"github.com/ribice/gorsk/pkg/utl/mock/mockdb"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateApartmentUserService(t *testing.T) {
	type args struct {
		c   echo.Context
		req apartment_user.ApartmentUser
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData apartment_user.ApartmentUser
		udb      *mockdb.ApartmentUser
	}{{
		name: "Success",
		args: args{req: apartment_user.ApartmentUser{
			Base: model.Base{
				ID: 1,
			},
			Name:        "karan",
			UserEmail:   "karan@123",
			UserAddress: "kolar",
		}},
		udb: &mockdb.ApartmentUser{
			CreateApartmentUserFn: func(c echo.Context, ap apartment_user.ApartmentUser) (apartment_user.ApartmentUser, error) {
				ap.CreatedAt = mock.TestTime(2000)
				ap.UpdatedAt = mock.TestTime(2000)
				ap.Base.ID = 1
				return ap, nil
			},
		},
		wantData: apartment_user.ApartmentUser{
			Base: model.Base{
				ID:        1,
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
			},
			Name:        "karan",
			UserEmail:   "karan@123",
			UserAddress: "kolar",
		}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_user_service.InitializeMockUser(tt.udb)
			apartmentUser, err := s.CreateApartmentUserService(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, apartmentUser)
		})
	}
}

func Test_ListApartmentUserService(t *testing.T) {
	type args struct {
		c      echo.Context
		filter apartment_user.UserFilter
	}
	cases := []struct {
		name     string
		args     args
		wantData []apartment_user.ApartmentUser
		wantErr  bool
		udb      *mockdb.ApartmentUser
	}{

		{
			name: "Success",
			args: args{c: nil},
			udb: &mockdb.ApartmentUser{
				ListApartmentUserFn: func(echo.Context, apartment_user.UserFilter) ([]apartment_user.ApartmentUser, error) {
					return []apartment_user.ApartmentUser{
						{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
							Name:        "karan",
							UserEmail:   "karan@123",
							UserAddress: "koalr",
						},
						{
							Base: model.Base{
								ID:        2,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							Name:        "rahul",
							UserEmail:   "rahul@123",
							UserAddress: "koalr",
						},
					}, nil
				}},
			wantData: []apartment_user.ApartmentUser{
				{
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
					Name:        "karan",
					UserEmail:   "karan@123",
					UserAddress: "koalr",
				},
				{
					Base: model.Base{
						ID:        2,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					Name:        "rahul",
					UserEmail:   "rahul@123",
					UserAddress: "koalr",
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_user_service.InitializeMockUser(tt.udb)
			apartmentsUsr, err := s.ListApartmentUserService(tt.args.c, tt.args.filter)
			assert.Equal(t, tt.wantData, apartmentsUsr)
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
		wantData apartment_user.ApartmentUser
		wantErr  error
		udb      *mockdb.ApartmentUser
	}{
		{
			name: "Success",
			args: args{id: 1},
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:        "karan",
				UserEmail:   "karan@123",
				UserAddress: "kolar",
			},
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					if id == 1 {
						return apartment_user.ApartmentUser{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							Name:        "karan",
							UserEmail:   "karan@123",
							UserAddress: "kolar",
						}, nil
					}
					return apartment_user.ApartmentUser{}, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_user_service.InitializeMockUser(tt.udb)
			apartmentUsr, err := s.ViewApartmentUserService(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, apartmentUsr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_UpdateApartmentService(t *testing.T) {
	type args struct {
		c   echo.Context
		upd apartment_user_service.NewUpdateApartmentUser
	}
	cases := []struct {
		name     string
		args     args
		wantData apartment_user.ApartmentUser
		wantErr  error
		udb      *mockdb.ApartmentUser
	}{

		{
			name: "Fail on Update",
			args: args{upd: apartment_user_service.NewUpdateApartmentUser{
				Id: 1,
			}},
			wantErr: gorsk.ErrGeneric,
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(1991),
						},
						Name:        "karan",
						UserEmail:   "karan@123",
						UserAddress: "kolar",
					}, nil
				},
				UpdateApartmentUserFn: func(c echo.Context, apartmentUsr apartment_user.ApartmentUser) error {
					return gorsk.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{upd: apartment_user_service.NewUpdateApartmentUser{
				Id:          1,
				Name:        "karan",
				UserEmail:   "karan@123",
				UserAddress: "kolar",
			}},
			wantData: apartment_user.ApartmentUser{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(1990),
					UpdatedAt: mock.TestTime(2000),
				},
				Name:        "karan",
				UserEmail:   "karan@123",
				UserAddress: "kolar",
			},
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:        "karan",
						UserEmail:   "karan@123",
						UserAddress: "kolar",
					}, nil
				},
				UpdateApartmentUserFn: func(c echo.Context, apartmentUsr apartment_user.ApartmentUser) error {
					apartmentUsr.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_user_service.InitializeMockUser(tt.udb)
			usr, err := s.UpdateApartmentUserService(tt.args.c, tt.args.upd)
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
		udb     *mockdb.ApartmentUser
	}{
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: gorsk.ErrGeneric,
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					if id != 1 {
						return apartment_user.ApartmentUser{}, nil
					}
					return apartment_user.ApartmentUser{}, gorsk.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockdb.ApartmentUser{
				ViewApartmentUserFn: func(c echo.Context, id int) (apartment_user.ApartmentUser, error) {
					return apartment_user.ApartmentUser{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						Name:        "karan",
						UserEmail:   "karan@123",
						UserAddress: "kolar",
					}, nil
				},
				DeleteApartmentUserFn: func(c echo.Context, apartmentUsr apartment_user.ApartmentUser) error {
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := apartment_user_service.InitializeMockUser(tt.udb)
			err := s.DeleteApartmentUserService(tt.args.c, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}
