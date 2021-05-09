package pkg

const (
	NumRanks = 13
	NumSuits = 4
	NumCards = NumRanks * NumSuits
)

type Rank int

const (
	NoRank Rank = 0
	LowAce Rank = 1
	Jack   Rank = 11
	Queen  Rank = 12
	King   Rank = 13
	Ace    Rank = 14
)

func (r Rank) String() string {
	switch r {
	case LowAce, Ace:
		return "A"
	case King:
		return "K"
	case Queen:
		return "Q"
	case Jack:
		return "J"
	case 10:
		return "T"
	}

	if r < 2 || r > 9 {
		return "?"
	}
	return string(rune(r) + '0')
}

func StringToRank(s string) Rank {
	if len(s) != 1 {
		return NoRank
	}

	switch s[0] {
	case 'A':
		return Ace
	case 'K':
		return King
	case 'Q':
		return Queen
	case 'J':
		return Jack
	case 'T':
		return 10
	}

	if s[0] < '0' || s[0] > '9' {
		return NoRank
	}

	return Rank(s[0] - '0')
}

type Suit int

const (
	NoSuit Suit = iota
	Diamond
	Club
	Heart
	Spade
)

func (s Suit) String() string {
	switch s {
	case Diamond:
		return "D"
	case Club:
		return "C"
	case Heart:
		return "H"
	case Spade:
		return "S"
	default:
		return "?"
	}
}

func StringToSuit(s string) Suit {
	if len(s) != 1 {
		return NoSuit
	}

	switch s[0] {
	case 'D':
		return Diamond
	case 'C':
		return Club
	case 'H':
		return Heart
	case 'S':
		return Spade
	default:
		return NoSuit
	}
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

func StringToCard(s string) Card {
	if len(s) != 2 {
		return Card{}
	}
	return Card{StringToRank(string(s[0])), StringToSuit(string(s[1]))}
}

func StringsToCards(s ...string) []Card {
	cards := make([]Card, len(s))
	for i := range s {
		cards[i] = StringToCard(s[i])
	}
	return cards
}
