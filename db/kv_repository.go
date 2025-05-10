package db

import (
	"database/sql"
	"log"
)

// key-value
type Kv struct {
	Key string
	Value string
}

type KvRepository struct {
	db *sql.DB
}

func NewKvRepository(db *sql.DB) *KvRepository {
	return &KvRepository{db: db}
}

func (r *KvRepository) FindByKey(key string) string {
	query := "SELECT value FROM kvs WHERE key = ?"

	row := r.db.QueryRow(query, key)

	var value string
	err := row.Scan(
		&value,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Panicln(err)
	}

	return value
}

func (r *KvRepository) Insert(kv *Kv) {
	cmd := `INSERT INTO kvs (key, value) VALUES (?, ?)`

	_, err := r.db.Exec(cmd,
		kv.Key,
		kv.Value,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *KvRepository) Update(kv *Kv) {
	cmd := `UPDATE kvs SET value = ? WHERE key = ?;`

	_, err := r.db.Exec(cmd,
		kv.Key,
		kv.Value,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *KvRepository) CreateTable() {
	cmd := `CREATE TABLE IF NOT EXISTS kvs(
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *KvRepository) Drop() {
	cmd := "DROP TABLE IF EXISTS kvs"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *KvRepository) Truncate() {
	cmd := "DELETE from kvs"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
