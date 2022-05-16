package service

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/car"
	"github.com/ribice/gorsk/pkg/api/car/platform/pgsql"
)

type Service interface {
	Create(echo.Context, car.EmpCar) (car.EmpCar, error)
	List(echo.Context, gorsk.Pagination) ([]car.EmpCar, error)
	View(echo.Context, int) (car.EmpCar, error)
	Delete(echo.Context, int) error
	Update(echo.Context, Update) (car.EmpCar, error)
}

type Car struct {
	db  *pg.DB
	udb UDB
}

func (u Car) Create(c echo.Context, req car.EmpCar) (car.EmpCar, error) {
	return u.udb.Create(u.db, req)
}

func (u Car) List(c echo.Context, p gorsk.Pagination) ([]car.EmpCar, error) {
	return u.udb.List(u.db, p)
}

func (u Car) View(c echo.Context, id int) (car.EmpCar, error) {
	return u.udb.View(u.db, id)
}

type Update struct {
	Id        int
	EmpName   string
	CarNumber string
	CarModel  string
}

func (u Car) Update(c echo.Context, r Update) (car.EmpCar, error) {

	if err := u.udb.Update(u.db, car.EmpCar{
		Id:        r.Id,
		EmpName:   r.EmpName,
		CarNumber: r.CarNumber,
		CarModel:  r.CarModel,
	}); err != nil {
		return car.EmpCar{}, err
	}

	return u.udb.View(u.db, r.Id)
}

func (u Car) Delete(c echo.Context, id int) error {
	empcar, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}
	return u.udb.Delete(u.db, empcar)
}

func Initialize(db *pg.DB) *Car {
	empCar := pgsql.Initialize(db)
	return &Car{db: db, udb: empCar}
}

type UDB interface {
	Create(orm.DB, car.EmpCar) (car.EmpCar, error)
	View(orm.DB, int) (car.EmpCar, error)
	List(orm.DB, gorsk.Pagination) ([]car.EmpCar, error)
	Update(orm.DB, car.EmpCar) error
	Delete(orm.DB, car.EmpCar) error
}
