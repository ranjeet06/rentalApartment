package apartment_controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/cache"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// HTTP represents apartment http service
type HTTP struct {
	ac    apartment_service.ApartmentService
	cache cache.ServiceCache
}

// NewHTTP creates new apartment http service
func NewHTTP(ac apartment_service.ApartmentService, cache cache.ServiceCache, r *echo.Group) {
	h := HTTP{ac, cache}
	ap := r.Group("/apartments")

	// swagger:route POST /v2/apartments apartments Create
	// Creates new apartment.
	// responses:
	//  200: apartmentResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	ap.POST("", h.CreateApartment)

	// swagger:operation GET /v2/apartments apartments listApartments
	// ---
	// summary: Returns list of apartments.
	// description: Returns list of apartments.
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
	ap.GET("", h.ListApartment)

	// swagger:operation GET /v2/apartments/{id} apartments getApartment
	// ---
	// summary: Returns a single apartment.
	// description: Returns a single apartment by its ID.
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
	ap.GET("/:id", h.ViewApartment)

	// swagger:operation PATCH /v2/apartments/{id} apartments apartmentUpdate
	// ---
	// summary: Updates apartment information
	// description: Updates apartment information -> name, description, floor_area, price_per_month, number_of_rooms,associated_realtor.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of apartment
	//   type: int
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/apartmentUpdate"
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
	ap.PATCH("/:id", h.UpdateApartment)

	// swagger:operation DELETE /v2/apartments/{id} apartments apartmentDelete
	// ---
	// summary: Deletes a apartment
	// description: Deletes a apartment with requested ID.
	// parameters:
	// - name: id
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
	ap.DELETE("/:id", h.DeleteApartment)
}

// CreateApartment create request
// swagger:model apartmentCreate
type CreateApartment struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	FloorArea     float64 `json:"floor_area"`
	PricePerMonth float64 `json:"price_per_month"`
	NumberOfRooms int     `json:"number_of_rooms"`

	GeolocationCoordinates apartment.Vertex `json:"geolocation_coordinates"`
	AssociatedRealtor      string           `json:"associated_realtor"`
}

func (h HTTP) CreateApartment(c echo.Context) error {
	a := new(CreateApartment)

	if err := c.Bind(a); err != nil {
		log.Err(err).Msg("json syntax err")
		return err
	}

	apartmentService, err := h.ac.CreateApartmentService(c, apartment.Apartment{
		Name:                   a.Name,
		Description:            a.Description,
		FloorArea:              a.FloorArea,
		PricePerMonth:          a.PricePerMonth,
		NumberOfRooms:          a.NumberOfRooms,
		GeolocationCoordinates: a.GeolocationCoordinates,
		AssociatedRealtor:      a.AssociatedRealtor,
	})

	if err != nil {
		log.Err(err)
		return err
	}

	key := "apartments"
	cacheSetErr := h.cache.Delete(key)
	if cacheSetErr != nil {
		log.Err(cacheSetErr)
		return cacheSetErr
	}

	return c.JSON(http.StatusOK, apartmentService)
}

func (h HTTP) ListApartment(c echo.Context) error {
	var requestPagination apartment.Pagination

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	var filter apartment.FilterApartment

	noOfRooms, err := strconv.Atoi(c.QueryParam("noOfRooms"))
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	pricePerMonth, err := strconv.ParseFloat(c.QueryParam("pricePerMonth"), 64)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	floorArea, err := strconv.ParseFloat(c.QueryParam("floorArea"), 64)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	requestPagination = apartment.Pagination{Limit: limit, Offset: offset}
	filter = apartment.FilterApartment{NumberOfRooms: noOfRooms, FloorArea: floorArea, PricePerMonth: pricePerMonth}

	key := "apartments"
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.ac.ListApartmentService(c, requestPagination, filter)
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

func (h HTTP) ViewApartment(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	key := "apartmentId" + PramId
	cacheResult := h.cache.Get(key)
	if cacheResult == redis.Nil {
		result, err := h.ac.ViewApartmentService(c, id)
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

// NewControllerUpdateApartment update request
// swagger:model apartmentUpdate
type NewControllerUpdateApartment struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	FloorArea         float64 `json:"floor_area"`
	PricePerMonth     float64 `json:"price_per_month"`
	NumberOfRooms     int     `json:"number_of_rooms"`
	AssociatedRealtor string  `json:"associated_realtor"`
}

func (h HTTP) UpdateApartment(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	a := new(NewControllerUpdateApartment)
	if err := c.Bind(a); err != nil {
		log.Err(err).Msg("json syntax err")
		return err
	}

	updateApartment, err := h.ac.UpdateApartmentService(c, apartment_service.NewUpdateApartment{
		ID:                id,
		Name:              a.Name,
		Description:       a.Description,
		FloorArea:         a.FloorArea,
		PricePerMonth:     a.PricePerMonth,
		NumberOfRooms:     a.NumberOfRooms,
		AssociatedRealtor: a.AssociatedRealtor,
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

	return c.JSON(http.StatusOK, updateApartment)

}

func (h HTTP) DeleteApartment(c echo.Context) error {
	PramId := c.Param("id")
	id, err := strconv.Atoi(PramId)
	if err != nil {
		log.Err(err)
		return gorsk.ErrBadRequest
	}

	if err := h.ac.DeleteApartmentService(c, id); err != nil {
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
