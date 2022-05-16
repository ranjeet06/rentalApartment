package mock

import (
	"github.com/ribice/gorsk/pkg/api/apartment_user"
)

type JWTUser struct {
	GenerateTokenFn func(apartment_user.ApartmentUser) (string, error)
}

// GenerateToken mock
func (j JWTUser) GenerateToken(user apartment_user.ApartmentUser) (string, error) {
	return j.GenerateTokenFn(user)
}
