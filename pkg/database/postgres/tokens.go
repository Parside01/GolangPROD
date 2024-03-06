package postgres

import (
	"database/sql"
	"solution/models"
)

var (
	createTokensTable = `CREATE TABLE IF NOT EXISTS tokens (
		user_id VARCHAR(255) NOT NULL,
		token VARCHAR(255) NOT NULL,
		CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id));`
	addToken    = `INSERT INTO tokens (user_id, token) VALUES ($1, $2)`
	deleteToken = `DELETE FROM tokens WHERE user_id = $1`
	canUse      = `SELECT * FROM tokens WHERE token = $1 and user_id = $2`
)

func (p *PostgresDB) setupeTokensTable() error {
	_, err := p.db.Exec(createTokensTable)
	return err
}

func (p *PostgresDB) WriteToken(token *models.Token) error {
	_, err := p.db.Exec(addToken, token.UserID, token.Token)
	return err
}

func (p *PostgresDB) DeleteTokenByUserID(id string) error {
	_, err := p.db.Exec(deleteToken, id)
	return err
}

func (p *PostgresDB) CanUseToken(user_id, token string) error {
	res, err := p.db.Query(canUse, token, user_id)
	if err != nil {
		return err
	}
	ok := res.Next()
	if !ok {
		return sql.ErrNoRows
	}
	return nil
}
