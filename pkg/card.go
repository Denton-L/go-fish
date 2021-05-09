package pkg

const (
	NumRanks = 13
	NumSuits = 4
	NumCards = NumRanks * NumSuits

	LowAce = 1
	Jack   = 11
	Queen  = 12
	King   = 13
	Ace    = 14
)

type Suit int

const (
	NoSuit Suit = iota
	Diamond
	Club
	Heart
	Spade
)

type Card struct {
	Rank int
	Suit Suit
}
