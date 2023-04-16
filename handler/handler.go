package handler

import (
	"encoding/gob"
	"net/url"
	"errors"

	"github.com/novadinn/social-network/service"
	"github.com/novadinn/social-network/model"
	
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Handler struct {
	Router *gin.Engine
	Service *service.Service
}

type session struct {
	IsLoggedIn bool
	User *model.User
}

func New() *Handler {
	var h Handler
	h.Service = service.New()
	gob.Register(&model.User{})

	h.Router = gin.Default()

	h.Router.LoadHTMLGlob("templates/*.html")

	h.Router.Use(methodOverride(h.Router))
	
	h.Router.GET("/", h.homeGetHandler())
	h.Router.GET("/login", h.loginGetHandler())
	h.Router.POST("/login", h.loginPostHandler())
	h.Router.POST("/logout", h.logoutGetHandler())
	h.Router.POST("/posts", h.postsPostHandler())
	h.Router.POST("/comments", h.commentsPostHandler())
	h.Router.POST("/user-follows", h.userFollowPostHandler())
	h.Router.GET("/@:username/following", h.userFollowingGetHandler())
	h.Router.GET("/@:username/followers", h.userFollowersGetHandler())
	h.Router.DELETE("/user-follows", h.userUnfollowDeleteHandler())
	h.Router.GET("/@:username", h.userGetHandler())
	h.Router.GET("/p/:id", h.postGetHandler())

	return &h
}

func (h *Handler) Run() error {
	err := h.Router.Run(":8080")
	if err != nil {
		return err
	}
	
	return nil
}

func (h *Handler) sessionFromRequest(c *gin.Context) session {
	user, ok := h.userFromContext(c)
	return session{IsLoggedIn: ok, User: user}
}

func (h *Handler) userFromContext(c *gin.Context) (*model.User, bool) {
	session, _ := h.Service.Store.Get(c.Request, "session")
	var user = &model.User{}
	var ok bool
	val := session.Values["user"]

	user, ok = val.(*model.User)

	return user, ok
}

func (h *Handler) contextWithUser(c *gin.Context, user model.User) {
	session, _ := h.Service.Store.Get(c.Request, "session")
	session.Values["user"] = user
	session.Save(c.Request, c.Writer)
}

func (h *Handler) resetContext(c *gin.Context) {
	session, _ := h.Service.Store.Get(c.Request, "session")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(c.Request, c.Writer)
}

func (h *Handler) pushErr(key string, err error) {
	h.Service.Push(key, err.Error())
}

func (h *Handler) popErr(key string) error {
	value, ok := h.Service.Pop(key)
	if !ok {
		return nil
	}
	
	s, ok := value.(string)
	if !ok {
		return nil
	}
	
	return errors.New(s)
}

func (h *Handler) pushForm(key string, value url.Values) {
	h.Service.Push(key, value)
}

func (h *Handler) popForm(key string) url.Values {
	value, ok := h.Service.Pop(key)
	if !ok {
		return nil
	}

	v, ok := value.(url.Values)
	if !ok {
		return nil
	}
	
	return v
}

func methodOverride(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			return
		}

		method := c.PostForm("_method")
		if method == "" {
			return
		}

		switch method {
		case "DELETE":
			c.Request.Method = "DELETE"

		case "PUT":
			c.Request.Method = "PUT"

		case "PATCH":
			c.Request.Method = "PATCH"
			
		default:
			return
		}

		c.Abort()
		r.HandleContext(c)
	}
}
