package handler

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) userFollowPostHandler() gin.HandlerFunc {
	return func (c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		id := c.PostForm("user_id")
		
		if err := h.Service.FollowUser(s.User.ID, id); err != nil {
			h.Service.Logger.Print(err)
			h.pushErr("follow_user_err", err)
			http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
			
			return
		}

		http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
	}
}

func (h *Handler) userUnfollowDeleteHandler() gin.HandlerFunc {
	return func (c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		id := c.PostForm("user_id")

		if err := h.Service.UnfollowUser(s.User.ID, id); err != nil {
			h.Service.Logger.Print(err)
			h.pushErr("follow_user_err", err)
			http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
			
			return
		}

		http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
	}
}
