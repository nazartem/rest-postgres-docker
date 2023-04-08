package prdlist

type ProductList struct {
	ID        int `json:"id"`
	NoteID    int `json:"note_id"`
	ProductID int `json:"product_id"`
	Amount    int `json:"amount"`
}
