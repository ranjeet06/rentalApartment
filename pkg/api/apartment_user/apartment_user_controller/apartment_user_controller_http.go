package apartment_user_controller

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/apartmentUserUtl/apartmentUserJwt"
	"github.com/ribice/gorsk/pkg/api/apartment_user"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/api/cache"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

// HTTP represents apartment_user http service

type HTTP struct {
	au    apartment_user_service.ApartmentUserService
	cache cache.ServiceCache
}

// NewHTTP creates new apartment_user http service
func NewHTTP(au apartment_user_service.ApartmentUserService, cache cache.ServiceCache, r *echo.Group, middlewareFunc echo.MiddlewareFunc) {
	h := HTTP{au, cache}
	ap := r.Group("/apartment_users")

	// swagger:route POST /v2/apartment_users  create apartment user
	// Creates new apartment user.
	// responses:
	//  200: apartmentResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	ap.POST("", h.CreateApartmentUser, middlewareFunc)

	// swagger:operation GET /v2/apartment_users list apartment users
	// ---
	// summary: Returns list of apartment users.
	// description: Returns list of apartment users.
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
	ap.GET("", h.ListApartmentUser, middlewareFunc)

	// swagger:operation GET /v2/apartment_users/{id} list apartment user
	// ---
	// summary: Returns a single apartment user.
	// description: Returns a single apartment user by its ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of apartment user
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
	ap.GET("/:id", h.ViewApartmentUser)

	// swagger:operation PATCH /v2/apartment_users/{id} update apartment user
	// ---
	// summary: Updates apartment user information
	// description: Updates apartment user information -> name, user_email, user_address.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of apartment user
	//   type: int
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/apartmentUserUpdate"
	// responses:
	//   "200":
	//     "$ref": "#/responses/apartmentResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ap.PATCH("/:id", h.UpdateApartmentUser)

	// swagger:operation DELETE /v2/apartment_users/{id} delete apartment user
	// ---
	// summary: Deletes a apartment user
	// description: Deletes a apartment user with requested ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of apartment user
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
	ap.DELETE("/:id", h.DeleteApartmentUser)

	ap.GET("", h.GetGithubUserinfo)
}

// CreateApartmentUser create request
// swagger:model apartmentUserCreate
type CreateApartmentUser struct {
	Name        string `json:"name"`
	UserEmail   string `json:"user_email"`
	UserAddress string `json:"user_address"`
	//UserRole
}

func (h HTTP) CreateApartmentUser(c echo.Context) error {

	a := new(CreateApartmentUser)

	if err := c.Bind(a); err != nil {
		log.Err(err)
		return err
	}
	apartmentUser, err := h.au.CreateApartmentUserService(c, apartment_user.ApartmentUser{
		Name:        a.Name,
		UserEmail:   a.UserEmail,
		UserAddress: a.UserAddress,
	})
	if err != nil {
		log.Err(err)
		return err
	}

	var NewApartmentUser apartment_user.ApartmentUser
	NewApartmentUser.Name = a.Name
	NewApartmentUser.UserEmail = a.UserEmail
	var jwtService apartmentUserJwt.Service
	jwtService, err = apartmentUserJwt.JwtNew("HS256", "ranjeetqwertyuiopasdfghjklzxcvbnmmnbvcxzlkjhgfdsapoiuytrewqqwertyuiopasdfghjklmnbvcxz", 15, 64)
	if err != nil {
		log.Err(err)
		return err
	}
	token, err := jwtService.GenerateToken(NewApartmentUser)
	if err != nil {
		log.Err(err)
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	key := "apartment_users"
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.JSON(http.StatusOK, apartmentUser)
}

func (h HTTP) ListApartmentUser(c echo.Context) error {
	var filter apartment_user.UserFilter
	name := c.QueryParam("name")
	userEmail := c.QueryParam("email")

	filter = apartment_user.UserFilter{Name: name, UserEmail: userEmail}

	key := "apartment_users"
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.au.ListApartmentUserService(c, filter)
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

func (h HTTP) ViewApartmentUser(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	key := "apartmentId" + PramId
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.au.ViewApartmentUserService(c, id)
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

// UpdateUser update request
// swagger:model apartmentUserUpdate
type UpdateUser struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UserEmail   string `json:"user_email"`
	UserAddress string `json:"user_address"`
}

func (h HTTP) UpdateApartmentUser(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	a := new(UpdateUser)
	if err := c.Bind(a); err != nil {
		log.Err(err)
		return err
	}

	updateUser, err := h.au.UpdateApartmentUserService(c, apartment_user_service.NewUpdateApartmentUser{
		Id:          id,
		Name:        a.Name,
		UserEmail:   a.UserEmail,
		UserAddress: a.UserAddress,
	})
	if err != nil {
		log.Err(err)
		return err
	}

	key := "apartmentId" + PramId
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.JSON(http.StatusOK, updateUser)
}

func (h HTTP) DeleteApartmentUser(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	if err = h.au.DeleteApartmentUserService(c, id); err != nil {
		log.Err(err)
		return err
	}

	key := "apartmentId" + PramId
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.NoContent(http.StatusOK)
}

func (h HTTP) GetGithubUserinfo(c echo.Context) error {

	//	code := c.QueryParam("code")
	//var cfg = new(config.Configuration)

	//	clientId := cfg.ClientIdSecret.ClientId
	//	clientSecret := cfg.ClientIdSecret.ClientSecret
	/*data := url.Values{
		"code":          {code},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}*/

	url := "https://github.com/login/oauth/access_token?code=b6ed38e21efd1b488057&client_id=ea17c1c4cadace41ba54&client_secret=3eda47f2e395ea18c113c954e246859cd9ce52c1"
	//method := "POST"

	httpClient := &http.Client{}

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Err(err)
		return err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		log.Err(err)
		return err
	}
	defer response.Body.Close()

	var responseBody map[string]interface{}

	json.NewDecoder(response.Body).Decode(&responseBody)

	token := responseBody["access_token"].(string)

	//token := "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"
	accessToken, err := h.au.AccessTokenService(c, apartment_user.AccessToken{Token: token})
	if err != nil {
		log.Err(err)
		return err
	}

	requestUserInfo, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Err(err)
		return err
	}
	a := "Bearer "
	tokenNew := a + accessToken
	requestUserInfo.Header.Set("Authorization", tokenNew)
	requestUserInfo.Header.Set("Content-Type", "application/json")

	responseUserInfo, err := httpClient.Do(requestUserInfo)
	if err != nil {
		log.Err(err)
		return err
	}

	var responseUserBody map[string]interface{}

	json.NewDecoder(responseUserInfo.Body).Decode(&responseUserBody)

	return c.JSON(http.StatusOK, responseUserBody)

}
