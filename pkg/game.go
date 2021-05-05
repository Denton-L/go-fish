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

type Hole [2]Card
type Deck []Card

type Phase int

const (
	PreFlop Phase = iota
	Flop
	Turn
	River
	Showdown
)

type GameState struct {
	Phase Phase
	Hole  Hole
	Flop  [3]Card
	Turn  Card
	River Card
}

type Engine interface {
	Calculate(state GameState) (float64, error)
}
