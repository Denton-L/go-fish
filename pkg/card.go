package pkg

import (
	"fmt"
)

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

	return fmt.Sprint(rank, c.Suit)
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

	c.Suit = StringToSuit(string(s[1]))

	return c
}

func StringsToCards(s ...string) []Card {
	cards := make([]Card, len(s))
	for i := range s {
		cards[i] = StringToCard(s[i])
	}
	return cards
}
