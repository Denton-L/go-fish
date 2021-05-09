package pkg

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
