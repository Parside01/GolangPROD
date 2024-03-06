package postgres

import (
	"solution/models"
	"time"
)

var (
	createFriendTable = `CREATE TABLE IF NOT EXISTS friends (
		login VARCHAR(255) NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		addedAt TIMESTAMP NOT NULL,
		CONSTRAINT login_fk FOREIGN KEY (login) REFERENCES users(login),
		CONSTRAINT user_id_friend_fk FOREIGN KEY (user_id) REFERENCES users(id));`

	writeFriend = `MERGE INTO friends AS f
					USING (
					SELECT $1 AS login, $2 AS user_id, $3::timestamp AS addedAt
					) AS t
					ON f.login = t.login AND f.user_id = t.user_id
					WHEN MATCHED THEN
					UPDATE SET addedAt = t.addedAt
					WHEN NOT MATCHED THEN
					INSERT (login, user_id, addedAt) VALUES (t.login, t.user_id, t.addedAt);
	`

	//writeFriend           = `INSERT INTO friends (login, user_id, addedAt) VALUES ($1, $2, $3)`
	getFriendByUserID     = `SELECT * FROM friends WHERE user_id = $1`
	deleteFriendsBuUserID = `DELETE FROM friends WHERE user_id = $1`
	getByUserIDAndLimit   = `SELECT * FROM friends
								WHERE user_id  = $1
								ORDER BY addedAt DESC
								LIMIT $2 OFFSET $3;`
	isFriend = `SELECT 1 FROM friends WHERE login = $1 AND user_id = $2`

	deleteFriendByLogin = `DELETE FROM friends WHERE login = $1 AND user_id = $2`
)

func (p *PostgresDB) DeleteFriendByLogin(login, userID string) error {
	_, err := p.db.Exec(deleteFriendByLogin, login, userID)
	return err
}
func (p *PostgresDB) IsFriend(login, userID string) (bool, error) {
	var res bool

	err := p.db.QueryRow(isFriend, login, userID).Scan(&res)
	return res, err
}

func (p *PostgresDB) setupeFriendTable() error {
	_, err := p.db.Exec(createFriendTable)
	return err
}
func (p *PostgresDB) WriteFriend(userId, friendLogin string) error {
	_, err := p.db.Exec(writeFriend, friendLogin, userId, time.Now().Format("2006-01-02T15:04:05Z07:00"))
	return err
}

func (p *PostgresDB) GetUserFriend(userID string) ([]*models.Friend, error) {
	rows, err := p.db.Query(getFriendByUserID, userID)
	if err != nil {
		return nil, err
	}
	res := []*models.Friend{}

	for rows.Next() {
		f := new(models.Friend)
		err := rows.Scan(&f.UserId, &f.FriendLogin, &f.AddedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (p *PostgresDB) GetUserFriendByLimit(user_id string, LIMIT, OFFSET int) ([]*models.Friend, error) {
	rows, err := p.db.Query(getByUserIDAndLimit, user_id, LIMIT, OFFSET)
	if err != nil {
		return nil, err
	}

	res := []*models.Friend{}
	for rows.Next() {
		f := new(models.Friend)
		if err := rows.Scan(&f.FriendLogin, &f.UserId, &f.AddedAt); err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (p *PostgresDB) DeleteFriendByUserID(userID string) error {
	_, err := p.db.Exec(deleteFriendsBuUserID, userID)
	return err
}
