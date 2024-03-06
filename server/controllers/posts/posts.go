package posts

import (
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	validatate "solution/pkg/validate"
	"solution/server/controllers"
	"solution/server/middleware"

	"github.com/labstack/echo/v4"
	"github.com/twharmon/gouid"
)

type PostController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

func NewPostController(logger *slog.Logger, db *postgres.PostgresDB) *PostController {
	return &PostController{
		logger: logger,
		db:     db,
	}
}

func (p *PostController) GetGroup() string {
	return "api/posts"
}

func (c *PostController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.AuthMiddleware(c.db),
	}
}

func (c *PostController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "POST",
			Path:    "/new",
			Handler: c.newPost,
		},
		&controllers.Handler{
			Method:  "GET",
			Path:    "/:postId",
			Handler: c.getPostById,
		},
		&controllers.Handler{
			Method:  "GET",
			Path:    "/feed/my",
			Handler: c.getMyPosts,
		},
		&controllers.Handler{
			Method:  "GET",
			Path:    "/feed/:login",
			Handler: c.getPostsByUser,
		},
		&controllers.Handler{
			Method:  "POST",
			Path:    "/:postId/like",
			Handler: c.likePost,
		},
		&controllers.Handler{
			Method:  "POST",
			Path:    "/:postId/dislike",
			Handler: c.dislikePost,
		},
	}
}

func (c *PostController) dislikePost(e echo.Context) error {
	postid := e.Param("postId")
	userid := e.Get("UserID").(string)
	if err := c.db.DislikePost(userid, postid); err != nil {
		c.logger.Error("posts.dislikePost: failed to dislike post: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}
	post, err := c.db.GetPostByID(userid, postid)
	if err != nil {
		c.logger.Error("posts.dislikePost: failed to get post: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, post)
}

func (c *PostController) likePost(e echo.Context) error {
	postid := e.Param("postId")
	userid := e.Get("UserID").(string)
	if err := c.db.LikePost(userid, postid); err != nil {
		c.logger.Error("posts.likePost: failed to like post: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}

	post, err := c.db.GetPostByID(userid, postid)
	if err != nil {
		c.logger.Error("posts.likePost: failed to get post: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}
	return e.JSON(http.StatusOK, post)
}

func (c *PostController) getPostsByUser(e echo.Context) error {
	user_id := e.Get("UserID").(string)
	tlogin := e.Param("login")

	limit, offset := e.QueryParam("limit"), e.QueryParam("offset")
	if limit == "" {
		limit = "5"
	}
	if offset == "" {
		offset = "0"
	}
	l, o, ok := validatate.IsValidPaginationParams(limit, offset)
	if !ok {
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "invalid pagination params"})
	}

	userlogin, err := c.db.GetUserLoginByID(user_id)
	if err != nil {
		c.logger.Error("posts.newPost: failed to get user login: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
	}
	posts, err := c.db.GetUserPostsByLogin(userlogin, tlogin, l, o)
	if err != nil {
		c.logger.Error("posts.getUserPosts: failed to get posts: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}
	return e.JSON(http.StatusOK, posts)
}

func (c *PostController) getMyPosts(e echo.Context) error {
	user_id := e.Get("UserID").(string)
	userlogin, err := c.db.GetUserLoginByID(user_id)
	if err != nil {
		c.logger.Error("posts.newPost: failed to get user login: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
	}

	limit, offset := e.QueryParam("limit"), e.QueryParam("offset")
	if limit == "" {
		limit = "5"
	}
	if offset == "" {
		offset = "0"
	}
	l, o, ok := validatate.IsValidPaginationParams(limit, offset)
	if !ok {
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "invalid pagination params"})
	}

	posts, err := c.db.GetUserPosts(userlogin, l, o)
	if err != nil {
		c.logger.Error("posts.getMyPosts: failed to get posts: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, posts)
}

func (c *PostController) getPostById(e echo.Context) error {
	user_id := e.Get("UserID").(string)
	postid := e.Param("postId")
	if postid == "" {
		c.logger.Error("posts.getPosById: postid param is nil")
		return echo.NewHTTPError(http.StatusBadRequest, "PostId param is nil")
	}

	post, err := c.db.GetPostByID(user_id, postid)
	if err != nil {
		c.logger.Error("posts.getPostById: failed to get post: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrorResponse{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, post)
}

func (c *PostController) newPost(e echo.Context) error {
	user_id := e.Get("UserID").(string)
	var post *models.Post
	if err := e.Bind(&post); err != nil {
		c.logger.Error("posts.newPost: failed to bind post: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	userlogin, err := c.db.GetUserLoginByID(user_id)
	if err != nil {
		c.logger.Error("posts.newPost: failed to get user login: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
	}

	post.Author = userlogin
	post.ID = gouid.String(16, gouid.Secure64Char)
	if err := c.db.WritePost(post); err != nil {
		c.logger.Error("posts.newPost: failed to write post: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}
	return e.JSON(http.StatusCreated, post)
}
