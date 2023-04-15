package handler

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func (h *Handler) logoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.resetContext(c)
		http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
	}
}
