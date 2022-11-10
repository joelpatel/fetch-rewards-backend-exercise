package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joelpatel/fetch-rewards-backend-exercise/db"
)

var validate = validator.New()

func AddTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transaction db.Transaction

		// request.body => the Golang struct variable
		if err := ctx.BindJSON(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing the request body"}) // sends response
			return
		}

		// Performs validation of the fields
		// currently just checking if fields
		// provided. Look at the struct Transaction
		// in db packge.
		validationErr := validate.Struct(transaction)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
			return
		}

		newTotal := db.Payer[transaction.Payer] + transaction.Points
		if newTotal < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot add this transaction as it'll make total for " + transaction.Payer + " less than 0"})
			return
		}

		db.Payer[transaction.Payer] = newTotal
		db.InsertTransaction(transaction)

		ctx.JSON(http.StatusOK, gin.H{"message": "added transaction successfully"})
	}
}
