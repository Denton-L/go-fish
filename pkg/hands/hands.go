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

func checkStraight(getLength func() int, getRank func(index int) pkg.Rank) (bool, int) {
	firstRank := getRank(0)

	length := getLength()
	runLength := 1
	lastRank := firstRank
	for i := 1; i < length; i++ {
		rank := getRank(i)

		if lastRank == rank+1 {
			runLength++

			if runLength == handSize {
				return true, i - runLength + 1
			}
		} else {
			runLength = 1
		}

		lastRank = rank
	}

	if firstRank == pkg.Ace && lastRank == 2 && runLength == handSize-1 {
		return true, -1
	}

	return false, 0
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

		hasStraightFlush, runStart := checkStraight(func() int {
			return len(suit)
		}, func(index int) pkg.Rank {
			return suit[index].Rank
		})

		if !hasStraightFlush {
			fiveCardRanking = Flush
			handBacking = handBacking[:handSize]
			copy(handBacking, suit)
			break
		}

		fiveCardRanking = StraightFlush
		handBacking = handBacking[:handSize]

		if runStart < 0 {
			copy(handBacking, suit[len(suit)-handSize+1:])
			lastCard := suit[0]
			lastCard.Rank = pkg.LowAce
			handBacking[len(handBacking)-1] = lastCard
			break
		}

		copy(handBacking, suit[runStart:])
		break
	}

	if fiveCardRanking == NoHand {
		hasStraight, runIndex := checkStraight(func() int {
			return len(ranks)
		}, func(index int) pkg.Rank {
			return ranks[index][0].Rank
		})

		if hasStraight {
			fiveCardRanking = Straight
			handBacking = handBacking[:handSize]

			lowStraight := runIndex < 0
			if runIndex < 0 {
				runIndex = len(ranks) - handSize + 1
			}

			for i := range handBacking {
				if runIndex >= len(ranks) {
					break
				}
				handBacking[i] = ranks[runIndex][0]
				runIndex++
			}

			if lowStraight {
				lastCard := ranks[0][0]
				lastCard.Rank = pkg.LowAce
				handBacking[len(handBacking)-1] = lastCard
			}
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
