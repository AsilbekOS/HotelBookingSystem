package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func PostgreSQL(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Printf("Databaseni ochishda xatolik - %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Databasega ulanish(Ping)da xatolik - %v", err)
		return nil, err
	}

	fmt.Println("Databasega muvaffaqiyatli ulandi!")
	return db, nil
}
