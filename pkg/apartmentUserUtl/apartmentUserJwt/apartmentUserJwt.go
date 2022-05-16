package apartmentUserJwt

import (
	"fmt"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"strings"
	"time"

	"github.com/ribice/gorsk"

	"github.com/dgrijalva/jwt-go"
)

var minSecretLen = 128

// New generates new JWT service necessary for auth middleware
func JwtNew(algo, secret string, ttlMinutes, minSecretLength int) (Service, error) {
	if minSecretLength > 0 {
		minSecretLen = minSecretLength
	}
	if len(secret) < minSecretLen {
		return Service{}, fmt.Errorf("jwt secret length is %v, which is less than required %v", len(secret), minSecretLen)
	}
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		return Service{}, fmt.Errorf("invalid jwt signing method: %s", algo)
	}

	return Service{
		key:  []byte(secret),
		algo: signingMethod,
		ttl:  time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the apartmentUserJwt token is valid.
	ttl time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// ParseToken parses token from Authorization header
func (s Service) ParseToken(authHeader string) (*jwt.Token, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, gorsk.ErrGeneric
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if s.algo != token.Method {
			return nil, gorsk.ErrGeneric
		}
		return s.key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (s Service) GenerateToken(apartmentUser apartment_user.ApartmentUser) (string, error) {
	return jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"userName":  apartmentUser.Name,
		"userEmail": apartmentUser.UserEmail,
		"exp":       time.Now().Add(s.ttl).Unix(),
	}).SignedString(s.key)

}

// GetDummyToken generates dummy token

func (s Service) GetDummyToken() (string, error) {
	apartmentUser := apartment_user.ApartmentUser{Name: "ranjeet", UserEmail: "ranjeet@123"}
	DummyToken, err := s.GenerateToken(apartmentUser)

	return DummyToken, err
}
