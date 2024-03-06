package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"solution/models"
	"time"
)

var (
	createPostsTable = `CREATE TABLE IF NOT EXISTS posts (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		content VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL, 
		createdAt TIMESTAMP NOT NULL,
		likesCount INT NOT NULL,
		dislikesCount INT NOT NULL,
		CONSTRAINT posts_auth_fk FOREIGN KEY (author) REFERENCES users(login));`
	createTagsTabel = `CREATE TABLE IF NOT EXISTS tags (
		tag VARCHAR(255) NOT NULL,
		post_id VARCHAR(255) NOT NULL);`

	writePost = `INSERT INTO posts (id, content, author, createdAt, likesCount, dislikesCount) VALUES ($1, $2, $3, $4, $5, $6)`

	getpostByID     = `SELECT * FROM posts WHERE id = $1`
	insertTag       = `INSERT INTO tags (tag, post_id) VALUES ($1, $2)`
	getTagsByPostID = `SELECT tag FROM tags WHERE post_id = $1`

	getPostById = `SELECT p.id, p.content, p.author, p.createdAt, p.likesCount, p.dislikesCount, u.isPublic
					FROM posts AS p
					INNER JOIN users AS u ON $1 = p.author
					WHERE p.id = $2
					AND ( u.isPublic = true OR ( u.isPublic = false AND (u.login = p.author OR EXISTS (SELECT * FROM friends
																WHERE friends.login = $1
																AND friends.user_id = p.author
																	))));`
	getPostByAuthor = `SELECT * FROM posts WHERE author = $1;`
	geMyPosts       = `SELECT * FROM posts
				WHERE posts.author = $1
				ORDER BY createdAt DESC	
				LIMIT $2 OFFSET $3;`
	getAuthorByPosID    = `SELECT author FROM posts WHERE id = $1`
	getUserPostsByLogin = `SELECT * FROM posts p
  							INNER JOIN users AS u ON u.login = p.author
  							WHERE p.author = $1
  							AND (u.isPublic = true OR (u.isPublic = false AND ($2 = p.author OR EXISTS (SELECT * FROM friends WHERE friends.login = $2 AND friends.user_id = p.author))))
							ORDER BY p.createdAt DESC 
							LIMIT $3 OFFSET $4;`
	likePost    = `UPDATE posts SET likesCount = $2 WHERE id = $1`
	dislikePost = `UPDATE posts SET dislikesCount = $2 WHERE id = $1`
)

func (p *PostgresDB) LikePost(userid, post_id string) error {
	err := p.addReaction(context.Background(), post_id, userid, reaction_like)
	if err != nil {
		fmt.Println("ekbmee", err)
		return err
	}

	likecount, err := p.getLikesCount(post_id)
	if err != nil {
		fmt.Println("ebejbejbnje", err)
		return err
	}

	_, err = p.db.Exec(likePost, post_id, likecount)

	return err
}

func (p *PostgresDB) DislikePost(user_id, post_id string) error {
	err := p.addReaction(context.Background(), post_id, user_id, reaction_dislike)
	if err != nil {
		fmt.Println("ekbmee", err)
		return err
	}

	dislikecount, err := p.getDislikesCount(post_id)
	if err != nil {
		fmt.Println("ebejbejbnje", err)
		return err
	}

	_, err = p.db.Exec(dislikePost, post_id, dislikecount)
	return err
}

func (p *PostgresDB) sampleGetPostProc(postid string) (*models.Post, error) {
	res := new(models.Post)
	if err := p.db.QueryRow(getpostByID, postid).Scan(&res.ID, &res.Content, &res.Author, &res.CreatedAt, &res.LikesCount, &res.DislikesCount); err != nil {
		return nil, err
	}

	rows, err := p.db.Query(getTagsByPostID, res.ID)
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	res.Tags = tags
	return res, nil
}

func (p *PostgresDB) GetPostByID(userid, postid string) (*models.Post, error) {
	userlogin, err := p.GetUserLoginByID(userid)
	if err != nil {
		fmt.Println("wwbwkmbwbw")
		return nil, err
	}

	var author string
	err = p.db.QueryRow(getAuthorByPosID, postid).Scan(&author)
	if err != nil {
		fmt.Println("kmwkbmwbw")
		return nil, err
	}

	if author == userlogin {
		return p.sampleGetPostProc(postid)
	}
	if ok, err := p.IsFriend(author, userid); ok {
		if err != nil {
			fmt.Println("q,;qq.;qv,vq")
			return nil, err
		}
		return p.sampleGetPostProc(postid)
	}

	if ok, err := p.IsPublicUserProfile(author); ok {
		if err != nil {
			fmt.Println("o	gmgjgngjg")
			return nil, err
		}
		return p.sampleGetPostProc(postid)
	}
	return nil, errors.New("unknown post or user")
}
func (p *PostgresDB) GetUserPostsByLogin(userlogin, targetlogin string, limit, offset int) ([]*models.Post, error) {
	if userlogin == targetlogin {
		return p.GetUserPosts(targetlogin, limit, offset)
	}

	id, err := p.GetUserIDByLogin(userlogin)
	if err != nil {
		return nil, err
	}
	if ok, err := p.IsFriend(targetlogin, id); ok {
		if err != nil {
			return nil, err
		}
		return p.GetUserPosts(targetlogin, limit, offset)
	}

	if ok, err := p.IsPublicUserProfile(targetlogin); ok {
		if err != nil {
			return nil, err
		}
		return p.GetUserPosts(targetlogin, limit, offset)
	}

	return nil, errors.New("a non-existent or incorrect user")
}

func (p *PostgresDB) setupePostsTable() error {
	_, err := p.db.Exec(createPostsTable)
	if err != nil {
		return err
	}
	_, err = p.db.Exec(createTagsTabel)
	return err
}

func (p *PostgresDB) WritePost(post *models.Post) error {
	for i := 0; i < len(post.Tags); i++ {
		_, err := p.db.Exec(insertTag, post.Tags[i], post.ID)
		if err != nil {
			return err
		}
	}
	_, err := p.db.Exec(writePost, post.ID, post.Content, post.Author, time.Now().Format("2006-01-02T15:04:05Z07:00"), post.LikesCount, post.DislikesCount)
	return err
}

func (post *PostgresDB) GetUserPosts(login string, limit, offset int) ([]*models.Post, error) {
	rows, err := post.db.Query(geMyPosts, login, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	res := []*models.Post{}
	for rows.Next() {
		p := new(models.Post)
		if err := rows.Scan(&p.ID, &p.Content, &p.Author, &p.CreatedAt, &p.LikesCount, &p.DislikesCount); err != nil {
			return nil, err
		}
		rows, err := post.db.Query(getTagsByPostID, p.ID)
		if err != nil {
			return nil, err
		}
		tags := []string{}
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return nil, err
			}
			tags = append(tags, tag)
		}
		p.Tags = tags
		res = append(res, p)
	}
	return res, nil
}
