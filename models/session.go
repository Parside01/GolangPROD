package models

import (
	"encoding/json"
	"time"
)

type Session struct {
	SessionGUID  string    `json:"session_guid" hash:"session_guid" mapstructure:"session_guid"`
	UserID       string    `json:"user_id" hash:"user_id" mapstructure:"user_id"`
	UserAgent    string    `json:"user_agent" hash:"user_agent" mapstructure:"user_agent"`
	RefreshToken string    `json:"refresh_token" hash:"refresh_token" mapstructure:"refresh_token"`
	IP           string    `json:"ip" hash:"ip" mapstructure:"ip"`
	ExpiresAt    time.Time `json:"expires_at" hash:"expires_at" mapstructure:"expires_at"`
	CreatedAt    time.Time `json:"created_at" hash:"created_at" mapstructure:"created_at"`
}

func (m *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
