package validatate

import (
	"regexp"
	"solution/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

func IsValidPost(post *models.Post) error {
	if err := validation.ValidateStruct(post,
		validation.Field(&post.ID, validation.Required, validation.Length(1, 100)),
		validation.Field(&post.Content, validation.Required, validation.Length(-1, 1000)),
		validation.Field(&post.Author, validation.Required, validation.Length(4, 30), validation.Match(regexp.MustCompile("[a-zA-Z0-9]+"))),
		validation.Field(&post.Tags, validation.Required, validation.Each(validation.Length(-1, 20))),
		validation.Field(&post.CreatedAt, validation.Required, validation.Date("2006-01-02T15:04:05Z07:00")),
		validation.Field(&post.LikesCount, validation.Required, validation.Min(0)),
		validation.Field(&post.DislikesCount, validation.Required, validation.Min(0)),
	); err != nil {
		return err
	}
	return nil
}
