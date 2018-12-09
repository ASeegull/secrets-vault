package storage

import (
	"fmt"
	"time"

	"github.com/ASeegull/secrets-vault/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type DB struct {
	Conn *sqlx.DB
}

const connString = "user=%s password=%s dbname=%s sslmode=%s"

func InitDB(host, user, pass, dbname, sslmode string) (DB, error) {
	var (
		db  DB
		err error
	)

	dataSource := fmt.Sprintf(connString, user, pass, dbname, sslmode)
	for i := 0; i < 10; i++ {
		db.Conn, err = sqlx.Connect("postgres", dataSource)
		if err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	return db, err
}

func (db DB) Save(hash uuid.UUID, views int, expires time.Time, msg string) error {
	_, err := db.Conn.Exec("INSERT INTO secrets (hash, expireAfterViews, expireAfter, message) VALUES (?, ?, ?, ?)", hash, views, expires, msg)
	return err
}

func (db DB) Get(hash uuid.UUID) (*models.Secret, error) {
	var (
		s   *models.Secret
		err error
	)

	err = db.Conn.Get(s, "SELECT from secrets WHERE hash = ?", hash)
	return s, err
}

func (db DB) DecrementViews(hash uuid.UUID) (int, error) {
	var views int
	err := db.Conn.QueryRowx("UPDATE secrets SET expireAfterViews = expireAfterViews - 1) WHERE hash = ? RETURNING expireAfterViews", hash).Scan(&views)
	return views, err
}

func (db DB) Delete(uuid.UUID) error {
	return nil
}

func (db DB) ClearExpired() (string, error) {
	now := time.Now()
	res, err := db.Conn.Exec("DELETE from secrets WHERE expireAfterViews = 0 OR (expireAfter < ? AND expiredAfter NOT NULL)", now)
	return fmt.Sprintf("Executing clenup job result: %v", res), err
}
