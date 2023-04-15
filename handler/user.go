package handler

import (
	"net/http"

	_ "github.com/novadinn/social-network/service"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) userGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		username := c.Param("username")
		
		user, err := h.Service.GetFollowingUser(username, s.User.ID)
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		posts, err := h.Service.GetPostsByUsername(username)
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		err = h.popErr("follow_user_err")

		if err != nil {
			c.HTML(http.StatusOK, "user.html", gin.H{
				"Session": s,
				"User": user,
				"Posts": posts,
				"Err": err,
				"Error": err.Error(),
			})

			return
		}

		c.HTML(http.StatusOK, "user.html", gin.H{
			"Session": s,
			"User": user,
			"Posts": posts,
		})
	}
}
