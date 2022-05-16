package user_favorite_apartment_controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/cache"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_service"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// HTTP represents user favorite apartment http service
type HTTP struct {
	ufa   user_favorite_apartment_service.UserFavApartmentService
	cache cache.ServiceCache
}

// NewHTTP creates new user favorite apartment http service
func NewHTTP(ufa user_favorite_apartment_service.UserFavApartmentService, cache cache.ServiceCache, r *echo.Group) {
	h := HTTP{ufa, cache}

	fap := r.Group("/favorites")

	// swagger:route POST /v2/favorites add apartment favorite
	// add apartment favorite.
	// responses:
	//  200: apartmentResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	fap.POST("", h.AddUserFavApartment)

	// swagger:operation GET /v2/favorites list user's favorite apartments
	// ---
	// summary: Returns list of apartments and users.
	// description: Returns list of apartments and users.
	// parameters:
	// - name: limit
	//   in: query
	//   description: number of results
	//   type: int
	//   required: false
	// - name: offset
	//   in: query
	//   description: page number
	//   type: int
	//   required: false
	// responses:
	//   "200":
	//     "$ref": "#/responses/apartmentListResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	fap.GET("", h.ListUserFavApartment)

	// swagger:operation DELETE /v2/favorites/{user_id & apartment_id} remove from favorite list
	// ---
	// summary: remove from favorite list
	// description: remove from favorite list with ID.
	// parameters:
	// - name: user Id
	//   in: path
	//   description: id of user
	//   type: int
	//   required: true
	// - name: apartment ID
	//   in: path
	//   description: id of apartment
	//   type: int
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/err"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	fap.DELETE("", h.DeleteUserFavApartment)

	// swagger:operation GET /v2/apartments/{id} get favorite apartments for user
	// ---
	// summary: Returns a user and apartment.
	// description: Returns a apartment for single user .
	// parameters:
	// - name: id
	//   in: path
	//   description: id of apartment
	//   type: int
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/apartmentResp"
	//   "400":
	//     "$ref": "#/responses/err"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "404":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	fap.GET("/:id", h.ViewUserFavApartment)

}

// CreateUserFavApartment add favorite apartment
// swagger:model userFavoriteApartment
type CreateUserFavApartment struct {
	UserId      int `json:"user_id"`
	ApartmentId int `json:"apartment_id"`
}

func (h HTTP) AddUserFavApartment(c echo.Context) error {
	a := new(CreateUserFavApartment)

	if err := c.Bind(a); err != nil {
		log.Err(err)
		return err
	}

	favApartment, err := h.ufa.AddUserFavApartmentService(c, user_favorite_apartment.UserFavoriteApartment{
		UserID:      a.UserId,
		ApartmentID: a.ApartmentId,
	})
	if err != nil {
		log.Err(err)
		return err
	}

	key := "favorites"
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.JSON(http.StatusOK, favApartment)
}

func (h HTTP) ListUserFavApartment(c echo.Context) error {

	key := "apartment_users"
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.ufa.ListUserFavApartmentService(c)
		if err != nil {
			log.Err(err)
			return err
		}
		cacheErr := h.cache.Set(key, result)
		if err != nil {
			log.Err(cacheErr)
			return cacheErr
		}
		return c.JSON(http.StatusOK, result)
	} else {
		return c.JSON(http.StatusOK, cacheResult)
	}

}

// Custom error

var (
	ErrBadReq = echo.NewHTTPError(http.StatusBadRequest, "apartment already exist in favorite.")
)

func (h HTTP) ViewUserFavApartment(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		err = ErrBadReq
		return err
	}

	key := "apartmentId" + PramId
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.ufa.ViewUserFavApartmentService(c, id)
		if err != nil {
			log.Err(err)
			return err
		}
		cacheErr := h.cache.Set(key, result)
		if err != nil {
			log.Err(cacheErr)
			return cacheErr
		}
		return c.JSON(http.StatusOK, result)
	} else {
		return c.JSON(http.StatusOK, cacheResult)
	}

}

func (h HTTP) DeleteUserFavApartment(c echo.Context) error {
	userPramId := c.QueryParam("user_id")
	usedId, err := strconv.Atoi(userPramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	userApartmentId := c.QueryParam("apartment_id")
	apartmentId, err := strconv.Atoi(userApartmentId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	if err := h.ufa.DeleteUserFavApartmentService(c, usedId, apartmentId); err != nil {
		log.Err(err)
		return err
	}

	key := "apartmentId" + userPramId + userApartmentId
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.NoContent(http.StatusOK)
}
