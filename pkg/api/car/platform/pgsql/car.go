package pgsql

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/api/car"
	"net/http"
	"strings"
)

type EmpCar struct {
	db orm.DB
}

func Initialize(db orm.DB) *EmpCar {
	return &EmpCar{db: db}
}

var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "employee name or car number already exists.")
)

// Create creates a new user on database
func (u EmpCar) Create(db orm.DB, cr car.EmpCar) (car.EmpCar, error) {
	var ecar = new(car.EmpCar)
	err := u.db.Model(ecar).Where("emp_name = ? or car_number = ? ",
		strings.ToLower(cr.EmpName), strings.ToLower(cr.CarNumber)).Select()
	if err == nil || err != pg.ErrNoRows {
		return car.EmpCar{}, ErrAlreadyExists
	}

	err = u.db.Insert(&cr)
	return cr, err
}

func (u EmpCar) List(db orm.DB, p gorsk.Pagination) ([]car.EmpCar, error) {
	var ecar []car.EmpCar
	q := db.Model(&ecar).Order("emp_car.emp_name desc")
	err := q.Select()
	return ecar, err
}

func (u EmpCar) View(db orm.DB, id int) (car.EmpCar, error) {
	var empcar car.EmpCar
	sql := `SELECT "emp_cars".* 
	FROM "emp_cars"
	WHERE "emp_cars"."id" = ?`
	_, err := db.QueryOne(&empcar, sql, id)
	return empcar, err
}

func (u EmpCar) Update(db orm.DB, empcar car.EmpCar) error {
	_, err := db.Model(&empcar).Where(`"id" = ?`, empcar.Id).Update()
	return err
}

func (u EmpCar) Delete(db orm.DB, empcar car.EmpCar) error {
	return db.Delete(&empcar)
}
