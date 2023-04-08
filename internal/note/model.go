package note

import "time"

type Note struct {
	Number  int       `json:"number"`
	Date    time.Time `json:"date"`
	BuyerID int       `json:"buyer_id"`
}

type NoteWithPrdList struct {
	Number   int       `json:"number"`
	Date     time.Time `json:"date"`
	BuyerID  int       `json:"buyer_id"`
	PrdLists []PrdList `json:"prd_lists"`
}

type PrdList struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Amount     int     `json:"amount"`
	TotalCount float64 `json:"total_count"`
}
