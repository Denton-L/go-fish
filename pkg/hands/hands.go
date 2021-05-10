package hands

import (
	"sort"

	"github.com/Denton-L/go-fish/pkg"
)

type Ranking int

const (
	NoHand Ranking = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush

	handSize = 5
)

func (r Ranking) String() string {
	switch r {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	default:
		return "Invalid"
	}
}

type PokerHand struct {
	Ranking Ranking
	Hand    [handSize]pkg.Card
}

func (a PokerHand) CompareTo(b PokerHand) int {
	if a.Ranking != b.Ranking {
		return int(a.Ranking - b.Ranking)
	}
	for i := range a.Hand {
		if a.Hand[i].Rank != b.Hand[i].Rank {
			return int(a.Hand[i].Rank - b.Hand[i].Rank)
		}
	}
	return 0
}

func checkStraight(getLength func() int, getCard func(index int) pkg.Card) []pkg.Card {
	firstCard := getCard(0)

	length := getLength()
	runLength := 1
	lastCard := firstCard
	i := 0
	for i = 1; i < length; i++ {
		card := getCard(i)

		if lastCard.Rank == card.Rank+1 {
			runLength++

			if runLength == handSize {
				break
			}
		} else {
			runLength = 1
		}

		lastCard = card
	}

	if runLength == handSize || firstCard.Rank == pkg.Ace && lastCard.Rank == 2 && runLength == handSize-1 {
		straight := make([]pkg.Card, handSize)
		runIndex := i - handSize
		for j := range straight {
			runIndex++

			card := pkg.Card{}
			if runIndex < length {
				card = getCard(runIndex)
			} else {
				card = getCard(0)
				card.Rank = pkg.LowAce
			}

			straight[j] = card
		}
		return straight
	}

	return nil
}

func ComputeHand(cards []pkg.Card) PokerHand {
	if len(cards) < handSize {
		return PokerHand{}
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})

	ranks := [][]pkg.Card{}
	suits := [pkg.NumSuits + 1][]pkg.Card{}
	for _, card := range cards {
		lastRankIndex := len(ranks) - 1
		if lastRankIndex < 0 || ranks[lastRankIndex][0].Rank != card.Rank {
			ranks = append(ranks, nil)
		}
		ranks[len(ranks)-1] = append(ranks[len(ranks)-1], card)
		suits[card.Suit] = append(suits[card.Suit], card)
	}

	histogram := [pkg.NumRanks + 1][][]pkg.Card{}
	for _, rank := range ranks {
		histogram[len(rank)] = append(histogram[len(rank)], rank)
	}

	hand := PokerHand{}
	handBacking := hand.Hand[:0]

	fiveCardRanking := NoHand

	// technically, these doesn't work for 10 or more cards but that's irrelevant
	for _, suit := range suits {
		if len(suit) < handSize {
			continue
		}

		straightFlush := checkStraight(func() int {
			return len(suit)
		}, func(index int) pkg.Card {
			return suit[index]
		})

		handBacking = handBacking[:handSize]

		if straightFlush != nil {
			fiveCardRanking = StraightFlush
			copy(handBacking, straightFlush)
			break
		}

		fiveCardRanking = Flush
		copy(handBacking, suit)
		break
	}

	if fiveCardRanking == NoHand {
		straight := checkStraight(func() int {
			return len(ranks)
		}, func(index int) pkg.Card {
			return ranks[index][0]
		})

		if straight != nil {
			fiveCardRanking = Straight
			handBacking = handBacking[:handSize]
			copy(handBacking, straight)
		}
	}

	switch {
	case fiveCardRanking == StraightFlush:
		hand.Ranking = StraightFlush

	case len(histogram[4]) >= 1:
		hand.Ranking = FourOfAKind
		handBacking = handBacking[:4]
		copy(handBacking, histogram[4][0])

	case len(histogram[3]) >= 1 && len(histogram[2]) >= 1 || len(histogram[3]) >= 2:
		hand.Ranking = FullHouse
		handBacking = handBacking[:5]
		copy(handBacking[0:3], histogram[3][0])

		hist2Rank := pkg.Rank(pkg.NoRank)
		if len(histogram[2]) >= 1 {
			hist2Rank = histogram[2][0][0].Rank
		}

		hist3Rank := pkg.Rank(pkg.NoRank)
		if len(histogram[3]) >= 2 {
			hist3Rank = histogram[3][1][0].Rank
		}

		if hist2Rank > hist3Rank {
			copy(handBacking[3:], histogram[2][0])
		} else {
			copy(handBacking[3:], histogram[3][1])
		}

	case fiveCardRanking != NoHand: // Flush or Straight
		hand.Ranking = fiveCardRanking

	case len(histogram[3]) == 1:
		hand.Ranking = ThreeOfAKind
		handBacking = handBacking[:3]
		copy(handBacking, histogram[3][0])

	case len(histogram[2]) >= 2:
		hand.Ranking = TwoPair
		handBacking = handBacking[:4]
		copy(handBacking[0:2], histogram[2][0])
		copy(handBacking[2:4], histogram[2][1])

	case len(histogram[2]) == 1:
		hand.Ranking = OnePair
		handBacking = handBacking[:2]
		copy(handBacking, histogram[2][0])

	default:
		hand.Ranking = HighCard
	}

fillHand:
	for _, card := range cards {
		if len(handBacking) == cap(handBacking) {
			break
		}

		for _, hand := range handBacking {
			if card == hand {
				continue fillHand
			}
		}

		handBacking = handBacking[:len(handBacking)+1]
		handBacking[len(handBacking)-1] = card
	}

	return hand
}
