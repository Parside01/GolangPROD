package models

import (
	"time"
)

type Friend struct {
	UserId      string    `json:"user_id" hash:"user_id" mapstructure:"user_id"`
	FriendLogin string    `json:"friend_login" hash:"friend_login" mapstructure:"friend_login"`
	AddedAt     time.Time `json:"added_at" hash:"added_at" mapstructure:"added_at"`
}

type FriendResult struct {
	FriendLogin string    `json:"friend_login" hash:"friend_login" mapstructure:"friend_login"`
	AddedAt     time.Time `json:"added_at" hash:"added_at" mapstructure:"added_at"`
}

func GetFriendResult(input []*Friend) []*FriendResult {
	res := []*FriendResult{}
	for _, i := range input {
		res = append(res, &FriendResult{FriendLogin: i.FriendLogin, AddedAt: i.AddedAt})
	}
	return res
}
