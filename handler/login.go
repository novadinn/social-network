package handler

import (
	"net/http"
	"errors"

	"github.com/novadinn/social-network/service"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) loginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func (h *Handler) loginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		email := c.PostForm("email")
		username := c.PostForm("username")

		user, err := h.Service.Login(email, username)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) ||
				errors.Is(err, service.ErrUsernameTaken) ||
				errors.Is(err, service.ErrInvalidEmail) ||
				errors.Is(err, service.ErrInvalidUsername) {
				
				c.HTML(http.StatusOK, "login.html", gin.H{
					"Session": s,
					"Err": err,
					"Error": err.Error(),
					"Form": c.Request.PostForm,
				})
				
				return
			}
			
			h.Service.Logger.Printf("cannot login: %s", err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		h.contextWithUser(c, user)
		http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
	}
}
