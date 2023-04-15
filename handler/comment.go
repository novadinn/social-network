package handler

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) commentsPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		postID := c.PostForm("post_id")
		content := c.PostForm("content")

		if _, err := h.Service.CreateComment(s.User.ID, postID, content, s.IsLoggedIn); err != nil {
			h.pushErr("create_comment_err", err)
			h.pushForm("create_comment_form", c.Request.PostForm)
			http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
			
			return 
		}

		http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusFound)
	}
}
