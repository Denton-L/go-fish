package oddsengine

import (
	"math/rand"

	"github.com/Denton-L/go-fish/pkg"
)

func getRemainingCards(state pkg.GameState) pkg.Deck {
	deck := pkg.Deck{}
	for rank := pkg.Rank(1); rank <= pkg.NumRanks; rank++ {
	cardLoop:
		for suit := pkg.Suit(1); suit <= pkg.NumSuits; suit++ {
			card := pkg.Card{rank, suit}

			switch state.Phase {
			case pkg.Showdown, pkg.River:
				if card == state.River {
					continue cardLoop
				}
				fallthrough
			case pkg.Turn:
				if card == state.Turn {
					continue cardLoop
				}
				fallthrough
			case pkg.Flop:
				for _, flop := range state.Flop {
					if card == flop {
						continue cardLoop
					}
				}
				fallthrough
			case pkg.PreFlop:
				for _, hand := range state.Hole {
					if card == hand {
						continue cardLoop
					}
				}
			}

			deck = append(deck, card)
		}
	}
	return deck
}

func determineResult(state pkg.GameState, deck pkg.Deck) result {
	return loss
}

type result int

const (
	loss result = iota
	win
	draw
)

type OddsEngine struct {
	Workers    int
	Iterations int
}

func (oe OddsEngine) Calculate(state pkg.GameState) (float64, error) {
	decks := make(chan pkg.Deck, oe.Workers)
	results := make(chan result, oe.Workers)

	go func() {
		remainingCards := getRemainingCards(state)

		for i := 0; i < oe.Iterations; i++ {
			deck := make(pkg.Deck, len(remainingCards))
			copy(deck[:], remainingCards)
			rand.Shuffle(len(deck), func(i, j int) {
				deck[i], deck[j] = deck[j], deck[i]
			})
			decks <- deck
		}
		close(decks)
	}()

	for i := 0; i < oe.Workers; i++ {
		go func() {
			for deck := range decks {
				results <- determineResult(state, deck)
			}
		}()
	}

	losses := 0
	wins := 0
	draws := 0

	for i := 0; i < oe.Iterations; i++ {
		switch <-results {
		case loss:
			losses++
		case win:
			wins++
		case draw:
			draws++
		}
	}

	return float64(wins) / float64(losses+wins), nil
}
