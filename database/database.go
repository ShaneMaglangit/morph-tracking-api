package database

import (
	_ "embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
)

type AxieDB struct {
	conn *sqlx.DB
}

//go:embed queries/createTables.sql
var createTablesQuery string

func New() *AxieDB {
	// Create connection
	db, err := initConnection()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(createTablesQuery)
	if err != nil {
		log.Fatal(err)
	}
	return &AxieDB{db}
}

func initConnection() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db AxieDB) SaveAxieMultiple(axies []Axie) {
	if len(axies) == 0 {
		return
	}
	query := "INSERT IGNORE INTO axie_morphed(hash, blockNumber, timestamp, tokenId) VALUES (:hash, :blockNumber, :timestamp, :tokenId)"
	res, err := db.conn.NamedExec(query, axies)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SaveTxMultiple:", rows, "row(s) affected")
}

func (db AxieDB) GetLatestBlock() uint64 {
	var latestBlock uint64
	row := db.conn.QueryRow("SELECT blockNumber FROM axie_morphed ORDER BY blockNumber DESC LIMIT 1")
	err := row.Scan(&latestBlock)
	if err != nil {
		return 6000000
	}
	return latestBlock
}

func (db AxieDB) SelectAxies(page int, ascending bool, byId bool) []Axie {
	// Get offset param
	offset := page * 100

	// Get order by param
	orderBy := "timestamp"
	if byId {
		orderBy = "tokenId"
	}

	// Get order param
	order := "DESC"
	if ascending {
		order = "ASC"
	}

	// Construct the query
	query := fmt.Sprintf("SELECT tokenId, timestamp FROM axie_morphed ORDER BY %s %s LIMIT 100 OFFSET %d", orderBy, order, offset)

	// Select Axies from the DB
	var axies []Axie
	err := db.conn.Select(&axies, query)
	if err != nil {
		log.Fatal(err)
	}
	return axies
}
