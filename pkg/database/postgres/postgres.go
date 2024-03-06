package postgres

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUniqueViolation = errors.New("unique violation")
)

type PostgresDB struct {
	logger *slog.Logger
	db     *sqlx.DB
	URL    string
}

func New(logger *slog.Logger) *PostgresDB {
	url := os.Getenv("POSTGRES_CONN")
	postgres := &PostgresDB{
		logger: logger,
		URL:    url,
	}
	err := postgres.setupConnection()
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info("postgres.New: connection established")
	return postgres
}

func (p *PostgresDB) setupConnection() error {
	db, err := sqlx.Connect("pgx", p.URL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	p.db = db
	err = p.setupeUsersTable()
	if err != nil {
		return err
	}

	if err := p.setupeTokensTable(); err != nil {
		return err
	}

	if err := p.setupeFriendTable(); err != nil {
		return err
	}

	if err := p.setupePostsTable(); err != nil {
		return err
	}

	if err := p.setupReactionTable(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Close() {
	p.db.Close()
}
