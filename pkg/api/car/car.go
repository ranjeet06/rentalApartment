package car

type EmpCar struct {
	Id        int    `json:"id"`
	EmpName   string `json:"emp_name"`
	CarNumber string `json:"car_number"`
	CarModel  string `json:"car_model"`
}
