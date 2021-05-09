package hands

import (
	"testing"

	"github.com/Denton-L/go-fish/pkg"
)

func TestComputeHand(t *testing.T) {
	tests := []struct {
		cards   []pkg.Card
		ranking Ranking
		ranks   [handSize]int
	}{
		{
			pkg.StringsToCards("8C", "2S", "3H", "AS", "KD", "7S", "4S"),
			HighCard,
			[...]int{pkg.Ace, pkg.King, 8, 7, 4},
		},
		{
			pkg.StringsToCards("8C", "8D", "3H", "AS", "KD", "7S", "4S"),
			OnePair,
			[...]int{8, 8, pkg.Ace, pkg.King, 7},
		},
		{
			pkg.StringsToCards("8C", "8D", "2S", "3H", "AS", "KD", "KS"),
			TwoPair,
			[...]int{pkg.King, pkg.King, 8, 8, pkg.Ace},
		},
		{
			pkg.StringsToCards("8C", "8D", "3S", "3H", "AS", "KD", "KS"),
			TwoPair,
			[...]int{pkg.King, pkg.King, 8, 8, pkg.Ace},
		},
		{
			pkg.StringsToCards("8C", "8D", "3S", "2H", "AS", "KD", "8S"),
			ThreeOfAKind,
			[...]int{8, 8, 8, pkg.Ace, pkg.King},
		},
		{
			pkg.StringsToCards("9H", "6S", "7S", "2D", "6H", "8C", "5H"),
			Straight,
			[...]int{9, 8, 7, 6, 5},
		},
		{
			pkg.StringsToCards("KH", "AS", "7S", "QD", "TH", "JC", "5H"),
			Straight,
			[...]int{pkg.Ace, pkg.King, pkg.Queen, pkg.Jack, 10},
		},
		{
			pkg.StringsToCards("4H", "2S", "AS", "2D", "3H", "8C", "5H"),
			Straight,
			[...]int{5, 4, 3, 2, pkg.LowAce},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "KS", "JC", "8S", "2S"),
			Flush,
			[...]int{pkg.King, pkg.Queen, 10, 8, 2},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "JC", "AS", "AC"),
			FullHouse,
			[...]int{pkg.Ace, pkg.Ace, pkg.Ace, 10, 10},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "TC", "AS", "AC"),
			FullHouse,
			[...]int{pkg.Ace, pkg.Ace, pkg.Ace, 10, 10},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "TC", "AS", "TH"),
			FourOfAKind,
			[...]int{10, 10, 10, 10, pkg.Ace},
		},
		{
			pkg.StringsToCards("9H", "6S", "7H", "2D", "6H", "8H", "5H"),
			StraightFlush,
			[...]int{9, 8, 7, 6, 5},
		},
		{
			pkg.StringsToCards("KS", "AS", "7S", "QS", "TS", "JS", "5H"),
			StraightFlush,
			[...]int{pkg.Ace, pkg.King, pkg.Queen, pkg.Jack, 10},
		},
		{
			pkg.StringsToCards("4D", "2S", "AD", "2D", "3D", "8D", "5D"),
			StraightFlush,
			[...]int{5, 4, 3, 2, pkg.LowAce},
		},
	}

	for tn, test := range tests {
		hand := ComputeHand(test.cards)
		if hand.Ranking != test.ranking {
			t.Error(tn, "incorrect ranking", "expected:", test.ranking, "actual:", hand.Ranking)
		}

		for i := range hand.Hand {
			if hand.Hand[i].Rank != test.ranks[i] {
				t.Error(tn, "incorrect ranks", "expected:", test.ranks, "actual:", hand.Hand)
				break
			}
		}
	}
}
