package handler

import (
	"net/http"
	"errors"

	"github.com/novadinn/social-network/service"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) postGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		post, err := h.Service.GetPostByID(id)
		if err != nil {
			h.Service.Logger.Print(err)
			s := h.sessionFromRequest(c)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		s := h.sessionFromRequest(c)
		comments, err := h.Service.GetComments(id)
		
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		err = h.popErr("create_comment_err")
		form := h.popForm("create_comment_form")

		if err != nil {
			c.HTML(http.StatusOK, "post.html", gin.H{
				"Session": s,
				"Post": post,
				"Comments": comments,
				"Err": err,
				"Error": err.Error(),
				"Form": form,
			})

			return
		}
		
		c.HTML(http.StatusOK, "post.html", gin.H{
			"Session": s,
			"Post": post,
			"Comments": comments,
		})
	}
}

func (h *Handler) postsPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		user, ok := h.userFromContext(c)
		content := c.PostForm("content")

		if err := h.Service.CreatePost(content, user.ID, ok); err != nil {
			if errors.Is(err, service.ErrInvalidPostContent) ||
				errors.Is(err, service.ErrUnauthenticated) {
				h.pushErr("create_post_err", err)
				h.pushForm("create_post_form", c.Request.PostForm)
				http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
				
				return
			}
			
			h.Service.Logger.Printf("cannot create post: %s", err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
	}
}
