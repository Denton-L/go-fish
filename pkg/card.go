package pkg

const (
	NumRanks = 13
	NumSuits = 4
	NumCards = NumRanks * NumSuits

	LowAce = 1
	Jack   = 11
	Queen  = 12
	King   = 13
	Ace    = 14
)

type Suit int

const (
	NoSuit Suit = iota
	Diamond
	Club
	Heart
	Spade
)

type Card struct {
	Rank int
	Suit Suit
}

func (c Card) String() string {
	rank := '?'
	switch c.Rank {
	case LowAce, Ace:
		rank = 'A'
	case King:
		rank = 'K'
	case Queen:
		rank = 'Q'
	case Jack:
		rank = 'J'
	case 10:
		rank = 'T'
	default:
		if 2 <= c.Rank && c.Rank <= 9 {
			rank = rune(c.Rank) + '0'
		}
	}

	suit := '?'
	switch c.Suit {
	case Diamond:
		suit = 'D'
	case Club:
		suit = 'C'
	case Heart:
		suit = 'H'
	case Spade:
		suit = 'S'
	}

	return string(rank) + string(suit)
}

func StringToCard(s string) Card {
	if len(s) != 2 {
		return Card{}
	}

	c := Card{}

	switch s[0] {
	case 'A':
		c.Rank = Ace
	case 'K':
		c.Rank = King
	case 'Q':
		c.Rank = Queen
	case 'J':
		c.Rank = Jack
	case 'T':
		c.Rank = 10
	default:
		if s[0] < '0' || s[0] > '9' {
			return Card{}
		}
		c.Rank = int(s[0] - '0')
	}

	switch s[1] {
	case 'D':
		c.Suit = Diamond
	case 'C':
		c.Suit = Club
	case 'H':
		c.Suit = Heart
	case 'S':
		c.Suit = Spade
	default:
		return Card{}
	}

	return c
}

func StringsToCards(s ...string) []Card {
	cards := make([]Card, len(s))
	for i := range s {
		cards[i] = StringToCard(s[i])
	}
	return cards
}
