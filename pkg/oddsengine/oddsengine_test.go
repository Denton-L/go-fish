package oddsengine

import (
	"math"
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
			0,
		}

		remainingCards := getRemainingCards(state)

		cards := make(map[pkg.Card]struct{})
		for _, c := range remainingCards {
			cards[c] = struct{}{}
		}

		if len(cards) != pkg.NumCards - test.offset {
			t.Error(tn, "wrong number of remaining cards", "expected:", test.offset, "actual:", len(cards))
		}

		for _, m := range missingCards[:test.offset] {
			if _, ok := cards[m]; ok {
				t.Error(tn, "expected missing card but present:", "expected missing:", m, "remainingCards:", remainingCards)
				break
			}
		}
	}
}

func makeOddsEngine() OddsEngine {
	return OddsEngine{4, 10000}
}

func TestOddsEngine(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tests := []struct {
		caseString string
		winPercentage float64
	}{
		// https://caniwin.com/texasholdem/preflop/heads-up.php
		{"AAo", 84.93},
		{"KKo", 82.11},
		{"QQo", 79.63},
		{"JJo", 77.15},
		{"TTo", 74.66},
		{"99o", 71.66},
		{"88o", 68.71},
		{"AKs", 66.21},
		{"77o", 65.72},
		{"AQs", 65.31},
		{"AJs", 64.39},
		{"AKo", 64.46},
		{"ATs", 63.48},
		{"AQo", 63.50},
		{"AJo", 62.53},
		{"KQs", 62.40},
		{"66o", 62.70},
		{"A9s", 61.50},
		{"ATo", 61.56},
		{"KJs", 61.47},
		{"A8s", 60.50},
		{"KTs", 60.58},
		{"KQo", 60.43},
		{"A7s", 59.38},
		{"A9o", 59.44},
		{"KJo", 59.44},
		{"55o", 59.64},
		{"QJs", 59.07},
		{"K9s", 58.63},
		{"A5s", 58.06},
		{"A6s", 58.17},
		{"A8o", 58.37},
		{"KTo", 58.49},
		{"QTs", 58.17},
		{"A4s", 57.13},
		{"A7o", 57.16},
		{"K8s", 56.79},
		{"A3s", 56.33},
		{"QJo", 56.90},
		{"K9o", 56.40},
		{"A5o", 55.74},
		{"A6o", 55.87},
		{"Q9s", 56.22},
		{"K7s", 55.84},
		{"JTs", 56.15},
		{"A2s", 55.50},
		{"QTo", 55.94},
		{"44o", 56.25},
		{"A4o", 54.73},
		{"K6s", 54.80},
		{"K8o", 54.43},
		{"Q8s", 54.41},
		{"A3o", 53.85},
		{"K5s", 53.83},
		{"J9s", 54.11},
		{"Q9o", 53.86},
		{"JTo", 53.82},
		{"K7o", 53.41},
		{"A2o", 52.94},
		{"K4s", 52.88},
		{"Q7s", 52.52},
		{"K6o", 52.29},
		{"K3s", 52.07},
		{"T9s", 52.37},
		{"J8s", 52.31},
		{"33o", 52.83},
		{"Q6s", 51.67},
		{"Q8o", 51.93},
		{"K5o", 51.25},
		{"J9o", 51.63},
		{"K2s", 51.23},
		{"Q5s", 50.71},
		{"T8s", 50.50},
		{"K4o", 50.22},
		{"J7s", 50.45},
		{"Q4s", 49.76},
		{"Q7o", 49.90},
		{"T9o", 49.81},
		{"J8o", 49.71},
		{"K3o", 49.33},
		{"Q6o", 48.99},
		{"Q3s", 48.93},
		{"98s", 48.85},
		{"T7s", 48.65},
		{"J6s", 48.57},
		{"K2o", 48.42},
		{"22o", 49.38},
		{"Q2s", 48.10},
		{"Q5o", 47.95},
		{"J5s", 47.82},
		{"T8o", 47.81},
		{"J7o", 47.72},
		{"Q4o", 46.92},
		{"97s", 46.99},
		{"J4s", 46.86},
		{"T6s", 46.80},
		{"J3s", 46.04},
		{"Q3o", 46.02},
		{"98o", 46.06},
		{"87s", 45.68},
		{"T7o", 45.82},
		{"J6o", 45.71},
		{"96s", 45.15},
		{"J2s", 45.20},
		{"Q2o", 45.10},
		{"T5s", 44.93},
		{"J5o", 44.90},
		{"T4s", 44.20},
		{"97o", 44.07},
		{"86s", 43.81},
		{"J4o", 43.86},
		{"T6o", 43.84},
		{"95s", 43.31},
		{"T3s", 43.37},
		{"76s", 42.82},
		{"J3o", 42.96},
		{"87o", 42.69},
		{"T2s", 42.54},
		{"85s", 41.99},
		{"96o", 42.10},
		{"J2o", 42.04},
		{"T5o", 41.85},
		{"94s", 41.40},
		{"75s", 40.97},
		{"T4o", 41.05},
		{"93s", 40.80},
		{"86o", 40.69},
		{"65s", 40.34},
		{"84s", 40.10},
		{"95o", 40.13},
		{"T3o", 40.15},
		{"92s", 39.97},
		{"76o", 39.65},
		{"74s", 39.10},
		{"T2o", 39.23},
		{"54s", 38.53},
		{"85o", 38.74},
		{"64s", 38.48},
		{"83s", 38.28},
		{"94o", 38.08},
		{"75o", 37.67},
		{"82s", 37.67},
		{"73s", 37.30},
		{"93o", 37.42},
		{"65o", 37.01},
		{"53s", 36.75},
		{"63s", 36.68},
		{"84o", 36.70},
		{"92o", 36.51},
		{"43s", 35.72},
		{"74o", 35.66},
		{"72s", 35.43},
		{"54o", 35.07},
		{"64o", 35.00},
		{"52s", 34.92},
		{"62s", 34.83},
		{"83o", 34.74},
		{"42s", 33.91},
		{"82o", 34.08},
		{"73o", 33.71},
		{"53o", 33.16},
		{"63o", 33.06},
		{"32s", 33.09},
		{"43o", 32.06},
		{"72o", 31.71},
		{"52o", 31.19},
		{"62o", 31.07},
		{"42o", 30.11},
		{"32o", 29.23},
	}

	for tn, test := range tests {
		card1 := pkg.Card{pkg.StringToRank(string(rune(test.caseString[0]))), pkg.Heart}
		card2 := pkg.Card{pkg.StringToRank(string(rune(test.caseString[1]))), pkg.Heart}
		if test.caseString[2] == 'o' {
			card2.Suit = pkg.Spade
		}
		state := pkg.GameState{
			Phase: pkg.PreFlop,
			Hole: pkg.Hole{card1, card2},
			Opponents: 1,
		}
		oe := makeOddsEngine()

		odds, err := oe.Calculate(state)
		if err != nil {
			t.Error(tn, "error returned:", err)
		}

		odds *= 100.0
		t.Log(tn, "odds:", odds)
		if math.Abs(odds - test.winPercentage) / (odds + test.winPercentage) * 2 >= 0.10 {
			t.Error(tn, "odds vary too much", "expected:", test.winPercentage, "actual:", odds)
		}
	}
}

func BenchmarkOddsEngineCalculate(b *testing.B) {
	oe := makeOddsEngine()
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
