package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joelpatel/fetch-rewards-backend-exercise/db"
)

func GetAllPayers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, db.Payer)
	}
}
