/*
in memory pseudo database
*/

package db

import (
	"sort"
	"time"

	"github.com/joelpatel/fetch-rewards-backend-exercise/utils"
)

var (
	Transactions []Transaction
	Payer        = make(map[string]int)
)

type Transaction struct {
	Payer     string    `json:"payer" validate:"required"`
	Points    int       `json:"points" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}

func InsertTransaction(t Transaction) {
	if t.Points < 0 {
		t.Points = -t.Points                    // making it positive so its more intuitive to understand & write code
		remainingPoints := SearchAndDestroy(&t) // remaining points from the previous transactions after subtract
		if remainingPoints <= 0 {
			return
		} else {
			t.Points = remainingPoints // else assign remaining points to current incoming transaction and insert that based on timestamp of the current tx.
		}
	}

	/*
		- search for the right place to insert using binary search
		- move rest to the right
		- assign new
		for a database it can be replaced with
		just a simple insertion where the table is kept
		sorted by a clustered index
	*/
	i := sort.Search(len(Transactions), func(i int) bool {
		return Transactions[i].Timestamp.After(t.Timestamp)
	})
	Transactions = append(Transactions, Transaction{})
	copy(Transactions[i+1:], Transactions[i:])
	Transactions[i] = t
}

/*
Goes over transactions from beginning
and updates/deletes the "rowes" if the payer
is same. i.e., if previously had payer "A" with points 10
and a new request to "A" with -5 points. value comes in,
then "A" will have points 5 and timestamp will be the one
of the second request (-5 one).
*/
func SearchAndDestroy(tx *Transaction) int {
	updateValue := 0
	remaining := tx.Points
	var markedForDeletion []int
	for i := 0; i < len(Transactions); i++ {
		if remaining == 0 {
			break
		}
		if Transactions[i].Payer != tx.Payer {
			continue
		}

		// update this index's transaction
		// check if it zeroes out after subtraction
		// 		if so then remove this element and continue
		// 		else remove this element but return the remainder

		previousPoint := Transactions[i].Points
		if previousPoint-utils.Min(previousPoint, remaining) != 0 {
			// means that leftovers from previous transaction
			updateValue = previousPoint - utils.Min(previousPoint, remaining)
		}
		markedForDeletion = append(markedForDeletion, i)
		remaining -= utils.Min(previousPoint, remaining)

	}

	// cleanup ops
	RemoveTransactions(markedForDeletion)
	/*
	 note: above line is not needed if we use a database
	 in which case we can remove rows in line or batch them
	*/
	return updateValue
}

/*
Goes over the markedForDeletion in reverse.
Removes the elements which were marked for deletion previously.
*/
func RemoveTransactions(markedForDeletion []int) {
	for i := len(markedForDeletion) - 1; i >= 0; i-- {
		index := markedForDeletion[i]
		copy(Transactions[index:], Transactions[index+1:])
		Transactions = Transactions[:len(Transactions)-1]
	}
}

func GetTotalPoints() int {
	total := 0
	for _, points := range Payer {
		total += points
	}
	return total
}
