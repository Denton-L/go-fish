package oddsengine

import (
	"testing"

	"github.com/Denton-L/go-fish/pkg"
)

func TestGetRemainingCards(t *testing.T) {
	missingCards := pkg.StringsToCards("4C", "KS", "AH", "8H", "9C", "AC", "KD")
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

func BenchmarkOddsEngineCalculate(b *testing.B) {
	oe := OddsEngine{4, 1000000}
	state := pkg.GameState{
		Phase: pkg.PreFlop,
		Hole:  pkg.Hole{pkg.StringToCard("2H"), pkg.StringToCard("7S")},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := oe.Calculate(state); err != nil {
			b.Error(err)
		}
	}
}
