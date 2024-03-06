package models

type Country struct {
	Name   string `json:"name" mapstructure:"name" db:"name"`
	Alpha2 string `json:"alpha2" mapstructure:"alpha2" db:"alpha2"`
	Alpha3 string `json:"alpha3" mapstructure:"alpha3" db:"alpha3"`
	Region string `json:"region" mapstructure:"region" db:"region"`
}
