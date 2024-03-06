package postgres

import (
	"database/sql"
	"solution/models"
)

var (
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		login VARCHAR(255) NOT NULL CONSTRAINT customers_login_key UNIQUE,
		email VARCHAR(255) NOT NULL CONSTRAINT customers_email_key UNIQUE,
		password VARCHAR NOT NULL,
		created_at TIMESTAMP NOT NULL,
		phone VARCHAR(255) NOT NULL CONSTRAINT customers_phone_key UNIQUE,
		countryCode VARCHAR(255) NOT NULL,
		isPublic BOOLEAN NOT NULL,
		image VARCHAR(255));`

	addUser        = `INSERT INTO users (id, login, email, password, created_at, phone, countryCode, isPublic, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	getbyemail     = `SELECT * FROM users WHERE email = $1`
	getbyid        = `SELECT * FROM users WHERE id = $1`
	deleteUser     = `DELETE FROM users WHERE id = $1`
	updateUser     = `UPDATE users SET login = $2, email = $3, password = $4, created_at = $5, phone = $6, countryCode = $7, isPublic = $8, image = $9 WHERE id = $1`
	existUser      = `SELECT * FROM users WHERE login = $1 AND password = $2 AND phone = $3`
	getallbyname   = `SELECT * FROM users WHERE login = $1`
	getUserByLogin = `SELECT* FROM users u
  					WHERE u.login = $1
					AND (
					u.isPublic = TRUE
					OR EXISTS (
						SELECT 1
						FROM friends f
						WHERE f.user_id = $2
						AND f.login = $1
					)
					);`
	isPublicProfile  = `SELECT isPublic FROM users WHERE login = $1`
	getUserLoginByID = `SELECT login FROM users WHERE id = $1`
	getUserIDByLogin = `SELECT id FROM users WHERE login = $1`
)

func (p *PostgresDB) GetUserIDByLogin(login string) (string, error) {
	row := p.db.QueryRow(getUserIDByLogin, login)
	var id string
	err := row.Scan(&id)
	return id, err
}

func (p *PostgresDB) IsPublicUserProfile(targetlogin string) (bool, error) {
	res := p.db.QueryRow(isPublicProfile, targetlogin)
	var isPublic bool
	err := res.Scan(&isPublic)
	return isPublic, err
}

func (p *PostgresDB) setupeUsersTable() error {
	_, err := p.db.Exec(createUsersTable)
	return err
}

func (p *PostgresDB) UpdateUser(user *models.User) error {
	_, err := p.db.Exec(updateUser, user.ID, user.Login, user.Email, user.Password, user.CreatedAt, user.Phone, user.CountryCode, user.IsPublic, user.Image)
	return err
}

func (p *PostgresDB) GetUserLoginByID(user_id string) (string, error) {
	row := p.db.QueryRow(getUserLoginByID, user_id)
	var login string
	err := row.Scan(&login)
	return login, err
}

func (p *PostgresDB) GetUserByLoginIfPublic(userid string, targetlogin string) (*models.User, error) {
	res := p.db.QueryRow(getUserByLogin, targetlogin, userid)
	out := new(models.User)
	err := res.Scan(&out.ID, &out.Login, &out.Email, &out.Password, &out.CreatedAt, &out.Phone, &out.CountryCode, &out.IsPublic, &out.Image)
	return out, err
}

func (p *PostgresDB) AddUser(user *models.User) error {
	_, err := p.db.Exec(addUser, user.ID, user.Login, user.Email, user.Password, user.CreatedAt, user.Phone, user.CountryCode, user.IsPublic, user.Image)
	//если мне кто-то подскажет как называется ошибка уникальности в драйвере то буду очень рад
	if err != nil {
		return ErrUniqueViolation
	}
	return err
}

func (p *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	res, err := p.db.Query(getbyemail, email)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	err = res.Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.CreatedAt, &user.Phone, &user.CountryCode, &user.IsPublic, &user.Image)

	return user, err
}

func (p *PostgresDB) UserIsExist(user *models.User) (bool, error) {
	var u models.User
	err := p.db.QueryRow(existUser, user.Login, user.Password, user.Phone).Scan(&u.ID, &u.Login, &u.Email, &u.Password, &u.CreatedAt, &u.Phone, &u.CountryCode, &u.IsPublic, &u.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *PostgresDB) GetUserByID(id string) (*models.User, error) {

	res, err := p.db.Query(getbyid, id)
	if err != nil {
		return nil, err
	}

	ok := res.Next()
	if !ok {
		return nil, sql.ErrNoRows
	}

	user := new(models.User)
	err = res.Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.CreatedAt, &user.Phone, &user.CountryCode, &user.IsPublic, &user.Image)

	return user, err
}

func (p *PostgresDB) DeleteUserByID(id string) error {
	_, err := p.db.Exec(deleteUser, id)
	return err
}

func (p *PostgresDB) GetUserByLogin(name string) (*models.User, error) {
	res, err := p.db.Query(getallbyname, name)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	res.Next()
	err = res.Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.CreatedAt, &user.Phone, &user.CountryCode, &user.IsPublic, &user.Image)
	if err != nil {
		return nil, err
	}

	return user, nil
}
