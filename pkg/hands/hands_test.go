package hands

import (
	"testing"

	"github.com/Denton-L/go-fish/pkg"
)

func TestPokerHandCompareTo(t *testing.T) {
	tests := []struct {
		cards1, cards2 []pkg.Card
		result         int
	}{
		{
			pkg.StringsToCards("8C", "2S", "3H", "AS", "KD", "7S", "5S"),
			pkg.StringsToCards("8C", "2S", "3H", "AS", "KD", "7S", "4S"),
			1,
		},
		{
			pkg.StringsToCards("8C", "2S", "3H", "AS", "KD", "7S", "4S"),
			pkg.StringsToCards("8C", "8D", "2S", "3H", "AS", "KD", "KS"),
			-1,
		},
		{
			pkg.StringsToCards("9H", "6S", "7S", "2D", "6H", "8C", "5H"),
			pkg.StringsToCards("8H", "2H", "5H", "7D", "9S", "6S", "6C"),
			0,
		},
	}

	for tn, test := range tests {
		hand1 := ComputeHand(test.cards1)
		hand2 := ComputeHand(test.cards2)

		comparison := hand1.CompareTo(hand2)
		if test.result^comparison < 0 {
			t.Error(tn, "incorrect CompareTo", "expected sign:", test.result, "actual:", comparison)
		}
	}
}

func TestComputeHand(t *testing.T) {
	tests := []struct {
		cards   []pkg.Card
		ranking Ranking
		ranks   [handSize]pkg.Rank
	}{
		{
			pkg.StringsToCards("8C", "2S", "3H", "AS", "KD", "7S", "4S"),
			HighCard,
			[...]pkg.Rank{pkg.Ace, pkg.King, 8, 7, 4},
		},
		{
			pkg.StringsToCards("8C", "8D", "3H", "AS", "KD", "7S", "4S"),
			OnePair,
			[...]pkg.Rank{8, 8, pkg.Ace, pkg.King, 7},
		},
		{
			pkg.StringsToCards("8C", "8D", "2S", "3H", "AS", "KD", "KS"),
			TwoPair,
			[...]pkg.Rank{pkg.King, pkg.King, 8, 8, pkg.Ace},
		},
		{
			pkg.StringsToCards("8C", "8D", "3S", "3H", "AS", "KD", "KS"),
			TwoPair,
			[...]pkg.Rank{pkg.King, pkg.King, 8, 8, pkg.Ace},
		},
		{
			pkg.StringsToCards("8C", "8D", "3S", "2H", "AS", "KD", "8S"),
			ThreeOfAKind,
			[...]pkg.Rank{8, 8, 8, pkg.Ace, pkg.King},
		},
		{
			pkg.StringsToCards("9H", "6S", "7S", "2D", "6H", "8C", "5H"),
			Straight,
			[...]pkg.Rank{9, 8, 7, 6, 5},
		},
		{
			pkg.StringsToCards("KH", "AS", "7S", "QD", "TH", "JC", "5H"),
			Straight,
			[...]pkg.Rank{pkg.Ace, pkg.King, pkg.Queen, pkg.Jack, 10},
		},
		{
			pkg.StringsToCards("4H", "2S", "AS", "2D", "3H", "8C", "5H"),
			Straight,
			[...]pkg.Rank{5, 4, 3, 2, pkg.LowAce},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "KS", "JC", "8S", "2S"),
			Flush,
			[...]pkg.Rank{pkg.King, pkg.Queen, 10, 8, 2},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "JC", "AS", "AC"),
			FullHouse,
			[...]pkg.Rank{pkg.Ace, pkg.Ace, pkg.Ace, 10, 10},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "TC", "AS", "AC"),
			FullHouse,
			[...]pkg.Rank{pkg.Ace, pkg.Ace, pkg.Ace, 10, 10},
		},
		{
			pkg.StringsToCards("AH", "TS", "QS", "TD", "TC", "AS", "TH"),
			FourOfAKind,
			[...]pkg.Rank{10, 10, 10, 10, pkg.Ace},
		},
		{
			pkg.StringsToCards("9H", "6S", "7H", "2D", "6H", "8H", "5H"),
			StraightFlush,
			[...]pkg.Rank{9, 8, 7, 6, 5},
		},
		{
			pkg.StringsToCards("KS", "AS", "7S", "QS", "TS", "JS", "5H"),
			StraightFlush,
			[...]pkg.Rank{pkg.Ace, pkg.King, pkg.Queen, pkg.Jack, 10},
		},
		{
			pkg.StringsToCards("4D", "2S", "AD", "2D", "3D", "8D", "5D"),
			StraightFlush,
			[...]pkg.Rank{5, 4, 3, 2, pkg.LowAce},
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
