package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/ribice/gorsk/pkg/api/apartment"
	"log"
	"os"
)

func main() {
	dbInsert := `INSERT INTO public.apartments VALUES (1, now(),now(),now()'','','',false,'ranjeet', 'Kolar',700,10000,3,'{236.4,125.3}','rahul');`
	var psn = os.Getenv("DATABASE_URL")

	u, err := pg.ParseURL(psn)
	checkErr(err)
	db := pg.Connect(u)
	_, err = db.Exec("SELECT 1")
	checkErr(err)
	createSchema(db, &apartment.Apartment{})

	_, err = db.Exec(dbInsert)
	if err != nil {
		log.Fatal(err)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		}))
	}
}
