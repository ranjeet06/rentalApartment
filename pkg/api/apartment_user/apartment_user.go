package apartment_user

import "github.com/ribice/gorsk/pkg/utl/model"

// ApartmentUser represents apartment user object
type ApartmentUser struct {
	model.Base
	Name        string `json:"name"`
	UserEmail   string `json:"user_email"`
	UserAddress string `json:"user_address"`
}

// UserFilter represent for filter in apartment user
type UserFilter struct {
	Name      string `json:"name"`
	UserEmail string `json:"user_email"`
}

type AccessToken struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}
