package model

// Update request model links an id to a price object to add.
type Update struct {
	ID    int
	Name  string
	Price Price
}
