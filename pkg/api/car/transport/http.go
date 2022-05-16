package transport

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/car"
	"github.com/ribice/gorsk/pkg/api/car/service"
	"net/http"
	"strconv"
)

type HTTP struct {
	svc service.Service
}

func NewHTTP(svc service.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/car")

	ur.POST("", h.create)

	ur.GET("", h.list)

	ur.GET("/:id", h.view)

	ur.PATCH("/:id", h.update)

	ur.DELETE("/:id", h.delete)
}

type createReq struct {
	Id        int    `json:"id"`
	EmpName   string `json:"emp_name"`
	CarNumber string `json:"car_number"`
	CarModel  string `json:"car_model"`
}

func (h HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {

		return err
	}

	cr, err := h.svc.Create(c, car.EmpCar{
		Id:        r.Id,
		EmpName:   r.EmpName,
		CarNumber: r.CarNumber,
		CarModel:  r.CarModel,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cr)
}

type listResponse struct {
	Users []car.EmpCar `json:"emp_car"`
	Page  int          `json:"page"`
}

func (h HTTP) list(c echo.Context) error {
	var req gorsk.PaginationReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	result, err := h.svc.List(c, req.Transform())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listResponse{result, req.Page})
}

func (h HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

type updateReq struct {
	Id        int    `json:"id"`
	EmpName   string `json:"emp_name"`
	CarNumber string `json:"car_number"`
	CarModel  string `json:"car_model"`
}

func (h HTTP) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	r := new(updateReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	ucar, err := h.svc.Update(c, service.Update{
		Id:        id,
		EmpName:   r.EmpName,
		CarNumber: r.CarNumber,
		CarModel:  r.CarModel,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ucar)
}

func (h HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return gorsk.ErrBadRequest
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
