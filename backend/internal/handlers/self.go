package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(context *gin.Context) {
	{
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
