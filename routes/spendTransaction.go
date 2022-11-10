package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joelpatel/fetch-rewards-backend-exercise/db"
	"github.com/joelpatel/fetch-rewards-backend-exercise/utils"
)

type SpendJSON struct {
	Points int `json:"points"`
}

type Spending struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}

func SpendTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody SpendJSON

		// request.body => Golang struct
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing the request body"})
			return
		}

		remainingPointsToSpend := requestBody.Points
		if remainingPointsToSpend > db.GetTotalPoints() {
			// sends an error to the requestor as its not possible to spend these much points
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to spend these much points, you do not have total points for it"})
			return
		}

		var spendings []Spending
		var markedForDeletion []int
		var updateTransaction db.Transaction // for adding the transaction at the end if any points remains from previous old points
		for i := 0; i < len(db.Transactions); i++ {
			if remainingPointsToSpend == 0 {
				break
			}

			tx := &db.Transactions[i]                      // so that uses less RAM
			if tx.Points <= 0 || db.Payer[tx.Payer] <= 0 { // db.Payer[tx.Payer] ==> current balance
				continue
			}

			amountToSpend := utils.Min(tx.Points, remainingPointsToSpend)
			if db.Payer[tx.Payer]-amountToSpend < 0 {
				continue
			}
			markedForDeletion = append(markedForDeletion, i)
			if tx.Points-amountToSpend > 0 {
				tx.Points -= amountToSpend
				tx.Timestamp = time.Now().UTC()
				updateTransaction = *tx
			}

			spending := Spending{
				Payer:  tx.Payer,
				Points: -amountToSpend,
			}
			spendings = append(spendings, spending)
			remainingPointsToSpend -= amountToSpend
			db.Payer[tx.Payer] -= amountToSpend
		}

		db.RemoveTransactions(markedForDeletion) // this can be replaced with batch removal for DB
		if updateTransaction != (db.Transaction{}) {
			db.InsertTransaction(updateTransaction)
		}
		ctx.JSON(http.StatusOK, spendings)
	}
}
