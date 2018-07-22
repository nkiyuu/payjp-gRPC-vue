package domain

type Item struct {
	ID          int64
	Name        string
	Discription string
	Amount      int64
}

type Items []Item
