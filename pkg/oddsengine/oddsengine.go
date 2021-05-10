package oddsengine

import (
	"math/rand"

	"github.com/Denton-L/go-fish/pkg"
	"github.com/Denton-L/go-fish/pkg/hands"
)

func getRemainingCards(state pkg.GameState) pkg.Deck {
	deck := pkg.Deck{}
	for rank := pkg.Rank(2); rank <= pkg.Ace; rank++ {
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

type result int

const (
	loss result = iota
	win
	draw
)

func determineResult(state pkg.GameState, deck pkg.Deck) result {
	deckIndex := 0
	drawCard := func() pkg.Card {
		card := deck[deckIndex]
		deckIndex++
		return card

	}

	opponentHoles := make([]pkg.Hole, state.Opponents)
	for i := range opponentHoles {
		opponentHoles[i] = pkg.Hole{drawCard(), drawCard()}
	}

	switch state.Phase {
	case pkg.PreFlop:
		for i := range state.Flop {
			state.Flop[i] = drawCard()
		}
		fallthrough
	case pkg.Flop:
		state.Turn = drawCard()
		fallthrough
	case pkg.Turn:
		state.River = drawCard()
	}

	computeHand := func(hole pkg.Hole) hands.PokerHand {
		cards := make([]pkg.Card, 7)
		copy(cards[0:2], hole[:])
		copy(cards[2:5], state.Flop[:])
		cards[5] = state.Turn
		cards[6] = state.River
		return hands.ComputeHand(cards)
	}

	hand := computeHand(state.Hole)
	worstResult := win

compareOpponents:
	for i := range opponentHoles {
		opponentHand := computeHand(opponentHoles[i])
		comparison := hand.CompareTo(opponentHand)

		switch {
		case comparison < 0:
			worstResult = loss
			break compareOpponents
		case comparison == 0:
			worstResult = draw
		}
	}

	return worstResult
}

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
