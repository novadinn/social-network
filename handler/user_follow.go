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

func (h *Handler) userFollowingGetHandler() gin.HandlerFunc {
	return func (c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		username := c.Param("username")
		users, err := h.Service.Queries.GetFollowingByUsername(username)
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		c.HTML(http.StatusOK, "following.html", gin.H{
			"Session": s,
			"Users": users,
		})
	}
}

func (h *Handler) userFollowersGetHandler() gin.HandlerFunc {
	return func (c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		username := c.Param("username")
		users, err := h.Service.Queries.GetFollowersByUsername(username)
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		c.HTML(http.StatusOK, "followers.html", gin.H{
			"Session": s,
			"Users": users,
		})
	}
}
