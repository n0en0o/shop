package shared

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func WaitForDB(db *sql.DB, maxAttempts int, interval time.Duration) error {
	for i := 1; i <= maxAttempts; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("WaitForDB: попытка %d/%d - %v", i, maxAttempts, err)
			time.Sleep(interval)
			continue
		}
		log.Printf("WaitForDB: БД доступна попытка(%d/%d)", i, maxAttempts)
		return nil
	}

	return fmt.Errorf("WaitForDB: БД недоступна после %d попыток", maxAttempts)
}
