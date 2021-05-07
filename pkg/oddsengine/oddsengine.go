package oddsengine

import (
	"github.com/Denton-L/go-fish/pkg"
)

func getRemainingCards(state pkg.GameState) pkg.Deck {
	deck := pkg.Deck{}
	for number := 1; number <= pkg.NumRanks; number++ {
	cardLoop:
		for suit := pkg.Suit(1); suit <= pkg.NumSuits; suit++ {
			card := pkg.Card{number, suit}

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
