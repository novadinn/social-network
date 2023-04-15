package handler

import (
	"net/http"

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

		err = h.popErr("create_post_err")
		form := h.popForm("create_post_form")

		if err != nil {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"Session": s,
				"Err": err,
				"Error": err.Error(),
				"Form": form,
				"Posts": posts,
			})

			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"Session": s,
			"Form": form,
			"Posts": posts,
		})
	}
}
