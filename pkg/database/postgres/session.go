package postgres

import (
	"context"
	"solution/models"
)

var (
	createSessionTable = `CREATE TABLE IF NOT EXISTS sessions (
		GUID VARCHAR(255) NOT NULL,
		UserID VARCHAR(255) NOT NULL,
		UserAgent VARCHAR(255) NOT NULL,
		RefreshToken VARCHAR(255) NOT NULL,
		IP VARCHAR(255) NOT NULL,
		CreatedAt TIMESTAMP NOT NULL,
		ExpiresAt TIMESTAMP NOT NULL,
		CONSTRAINT userid_fk FOREIGN KEY (UserID) REFERENCES users(id));`

	writeSession = `INSERT INTO sessions (GUID, UserID, UserAgent, RefreshToken, IP, CreatedAt, ExpiresAt) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	findByGuid   = `SELECT * FROM sessions WHERE GUID = $1`
	deleteByGuid = `DELETE FROM sessions WHERE GUID = $1`
	updateByGuid = `UPDATE sessions SET UserID = $2, UserAgent = $3, RefreshToken = $4, IP = $5, CreatedAt = $6, ExpiresAt = $7 WHERE GUID = $1`
)

func (p *PostgresDB) setupeSessionTable() error {
	_, err := p.db.Exec(createSessionTable)
	return err
}

func (p *PostgresDB) WriteSession(ctx context.Context, session *models.Session) error {
	_, err := p.db.Exec(writeSession, session.SessionGUID, session.UserID, session.UserAgent, session.RefreshToken, session.IP, session.CreatedAt, session.ExpiresAt)
	return err
}

func (p *PostgresDB) FindSessionByGUID(ctx context.Context, SessionGUID string) (*models.Session, error) {
	var session *models.Session
	err := p.db.QueryRow(findByGuid, SessionGUID).Scan(&SessionGUID, &session.UserID, &session.UserAgent, &session.RefreshToken, &session.IP, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (p *PostgresDB) DeleteSessionByGUID(ctx context.Context, SessionGUID string) error {
	_, err := p.db.Exec(deleteByGuid, SessionGUID)
	return err
}

func (p *PostgresDB) UpdateSessionByGUID(ctx context.Context, session *models.Session, SessionGUID string) error {
	_, err := p.db.Exec(updateByGuid, SessionGUID, session.UserID, session.UserAgent, session.RefreshToken, session.IP, session.CreatedAt, session.ExpiresAt)
	return err
}
