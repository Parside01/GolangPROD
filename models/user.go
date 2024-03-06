package models

import "time"

type User struct {
	ID          string    `json:"id" mapstructure:"id" db:"id"`
	Login       string    `json:"login" mapstructure:"login" db:"login"`
	Email       string    `json:"email" mapstructure:"email" db:"email"`
	Password    string    `json:"password" mapstructure:"password" db:"password"` //Будем хешировать его.
	CreatedAt   time.Time `json:"created_at" mapstructure:"created_at" db:"created_at"`
	Phone       string    `json:"phone" mapstructure:"phone" db:"phone"`
	CountryCode string    `json:"countryCode" mapstructure:"countryCode" db:"countryCode"`
	IsPublic    bool      `json:"isPublic" mapstructure:"isPublic" db:"isPublic"`
	Image       string    `json:"image" mapstructure:"image" db:"image"`
}

type UserProfile struct {
	Login       string `json:"login" mapstructure:"login" db:"login"`
	Email       string `json:"email" mapstructure:"email" db:"email"`
	CountryCode string `json:"countryCode" mapstructure:"countryCode" db:"countryCode"`
	IsPublic    bool   `json:"isPublic" mapstructure:"isPublic" db:"isPublic"`
	Phone       string `json:"phone" mapstructure:"phone" db:"phone"`
}

func (u *User) GetProfile() *UserProfile {
	return &UserProfile{
		Login:       u.Login,
		Email:       u.Email,
		Phone:       u.Phone,
		CountryCode: u.CountryCode,
		IsPublic:    u.IsPublic,
	}
}
