package handler

import (
	"net/http"

	"github.com/novadinn/social-network/model"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) homeGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.sessionFromRequest(c)
		
		posts, err := h.Service.GetPosts()
		if err != nil {
			h.Service.Logger.Print(err)
			c.HTML(http.StatusOK, "err.html", gin.H{
				"Session": s,
				"Error": err.Error(),
			})
			
			return
		}

		var followingPosts []model.Post
		if s.IsLoggedIn {

			users, err := h.Service.Queries.GetFollowingByUsername(s.User.Username)
			ids := make([]interface{}, 0)
			for _, v := range users {
				ids = append(ids, v.ID)
			}

			followingPosts, err = h.Service.GetFollowingPosts(ids)
			if err != nil {
				h.Service.Logger.Print(err)
				c.HTML(http.StatusOK, "err.html", gin.H{
					"Session": s,
					"Error": err.Error(),
				})
				
				return
			}
		}

		err = h.popErr("create_post_err")
		form := h.popForm("create_post_form")

		if err != nil {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"Session": s,
				"Err": err,
				"Error": err.Error(),
				"Form": form,
				"Posts": posts,
				"FollowPosts": followingPosts,
			})

			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"Session": s,
			"Form": form,
			"Posts": posts,
			"FollowPosts": followingPosts,
		})
	}
}
