package dbwork

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB interface {
	Close()
	GetAllQuotes() ([]Quote, error)
	CreateQuote(quote Quote) error
	GetRandomQuote() (Quote, error)
	GetAuthorQuotes(author string) ([]Quote, error)
	DeleteQuote(id int) error
}

type Quote struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type PostgresDbParams struct {
	DbName   string `json:"DbName"`
	Host     string `json:"Host"`
	User     string `json:"User"`
	Password string `json:"Password"`
	SslMode  string `json:"SslMode"`
}

type PostgresDb struct {
	config PostgresDbParams
	db     *sql.DB
}

func (postgres *PostgresDb) Close() {
	postgres.db.Close()
}

func (postgres *PostgresDb) CreateQuote(quote Quote) error {
	query := `INSERT INTO quotes
	         (author, quote)
	         VALUES($1, $2)`
	_, err := postgres.db.Exec(query, quote.Author, quote.Quote)
	if err != nil {
		return err
	}
	return nil
}

func (postgres *PostgresDb) GetAllQuotes() ([]Quote, error) {
	quotes := make([]Quote, 0)
	query := fmt.Sprintf(`SELECT * FROM quotes`)

	rows, err := postgres.db.Query(query)
	if err != nil {
		return quotes, err
	}
	defer rows.Close()

	for rows.Next() {
		temp_quote := Quote{}
		err = rows.Scan(&temp_quote.Id, &temp_quote.Author, &temp_quote.Quote)
		if err != nil {
			return quotes, err
		}

		quotes = append(quotes, temp_quote)
	}

	return quotes, nil
}

func (postgres *PostgresDb) GetRandomQuote() (Quote, error) {
	query := fmt.Sprintf(`SELECT * FROM quotes ORDER BY RANDOM() LIMIT 1`)
	quote := Quote{}
	err := postgres.db.QueryRow(query).Scan(&quote.Id, &quote.Author, &quote.Quote)
	if err != nil {
		return Quote{}, err
	}
	return quote, nil
}

func (postgres *PostgresDb) GetAuthorQuotes(author string) ([]Quote, error) {
	quotes := make([]Quote, 0)
	query := fmt.Sprintf(`SELECT * FROM quotes WHERE author = '%s'`, author)

	rows, err := postgres.db.Query(query)
	if err != nil {
		return quotes, err
	}
	defer rows.Close()

	for rows.Next() {
		temp_quote := Quote{}
		err = rows.Scan(&temp_quote.Id, &temp_quote.Author, &temp_quote.Quote)
		if err != nil {
			return quotes, err
		}
		quotes = append(quotes, temp_quote)
	}

	return quotes, nil
}

func (postgres *PostgresDb) DeleteQuote(id int) error {
	query := fmt.Sprintf(`DELETE FROM quotes WHERE id = '%v'`, id)
	_, err := postgres.db.Query(query)
	if err != nil {
		return err
	}

	return nil
}

func CreatePostgresDb(config PostgresDbParams) (DB, error) {
	connStr := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host,
		config.DbName,
		config.User,
		config.Password,
		config.SslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &PostgresDb{}, err
	}

	if err = db.Ping(); err != nil {
		return &PostgresDb{}, err
	}

	postgres := PostgresDb{db: db, config: config}
	if err = postgres.createTableIfNotExists(); err != nil {
		return &PostgresDb{}, err
	}
	return &postgres, nil
}

func (postgres *PostgresDb) createTableIfNotExists() error {
	const nameTable = "quotes"

	createQueryTable := `CREATE TABLE quotes(
	id BIGSERIAL PRIMARY KEY,
  author TEXT,
  quote TEXT
	);`

	existsTable, err := postgres.verifyTableExists(nameTable)
	if err != nil {
		return err
	}

	if existsTable {
		return nil
	}

	_, err = postgres.db.Exec(createQueryTable)
	if err != nil {
		return err
	}

	return nil
}

func (postgres *PostgresDb) verifyTableExists(table string) (bool, error) {
	var result string

	rows, err := postgres.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() && result != table {
		rows.Scan(&result)
	}
	return result == table, rows.Err()
}
