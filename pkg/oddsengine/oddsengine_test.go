package oddsengine

import (
	"testing"

	"github.com/Denton-L/go-fish/pkg"
)

func TestGetRemainingCards(t *testing.T) {
	missingCards := []pkg.Card{
		pkg.Card{4, pkg.Club},
		pkg.Card{pkg.King, pkg.Spade},
		pkg.Card{pkg.Ace, pkg.Heart},
		pkg.Card{8, pkg.Heart},
		pkg.Card{9, pkg.Club},
		pkg.Card{pkg.Ace, pkg.Club},
		pkg.Card{pkg.King, pkg.Diamond},
	}

	tests := []struct {
		phase  pkg.Phase
		offset int
	}{
		{pkg.PreFlop, 2},
		{pkg.Flop, 5},
		{pkg.Turn, 6},
		{pkg.River, 7},
		{pkg.Showdown, 7},
	}

	for tn, test := range tests {
		hand := pkg.Hole{}
		copy(hand[:], missingCards[0:2])
		flop := [3]pkg.Card{}
		copy(flop[:], missingCards[2:5])

		state := pkg.GameState{
			test.phase,
			hand,
			flop,
			missingCards[5],
			missingCards[6],
		}

		remainingCards := getRemainingCards(state)

		cards := make(map[pkg.Card]struct{})
		for _, c := range remainingCards {
			cards[c] = struct{}{}
		}

		for _, m := range missingCards[:test.offset] {
			if _, ok := cards[m]; ok {
				t.Error(tn, "expected missing card but present:", "expected missing:", m, "remainingCards:", remainingCards)
				break
			}
		}
	}
}
