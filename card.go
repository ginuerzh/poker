package poker

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	HgihCard = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// rank
const (
	Deuce = iota
	Trey
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace

	RankMask = 0xF
)

// suit
const (
	Club    = 0x8000
	Diamond = 0x4000
	Heart   = 0x2000
	Spade   = 0x1000

	SuitMask = 0xF000
)

const (
	NumCard = 52
)

type Card uint32

const NilCard = Card(0)

func ParseCard(c string) Card {
	if len(c) != 2 {
		return NilCard
	}
	b := []byte(c)

	rank := 0
	suit := 0

	switch b[0] {
	case 'c', 'C':
		suit = Club
	case 'd', 'D':
		suit = Diamond
	case 'h', 'H':
		suit = Heart
	case 's', 'S':
		suit = Spade
	}

	switch b[1] {
	case '2':
		rank = Deuce
	case '3':
		rank = Trey
	case '4':
		rank = Four
	case '5':
		rank = Five
	case '6':
		rank = Six
	case '7':
		rank = Seven
	case '8':
		rank = Eight
	case '9':
		rank = Nine
	case 't', 'T':
		rank = Ten
	case 'j', 'J':
		rank = Jack
	case 'q', 'Q':
		rank = Queen
	case 'k', 'K':
		rank = King
	case 'a', 'A':
		rank = Ace
	default:
		rank = -1
	}

	if suit == 0 || rank < 0 {
		return NilCard
	}

	card := (1 << uint32(16+rank)) | suit | (rank << 8) | primes[rank]

	return Card(card)
}

func (card Card) Rank() int {
	return int((card >> 8) & RankMask)
}

func (card Card) Suit() int {
	return int(card & SuitMask)
}

func (card Card) prime() int {
	return int(card & 0xFF)
}

func (card Card) MarshalJSON() ([]byte, error) {
	return []byte("\"" + card.String() + "\""), nil
}

func (card *Card) UnmarshalJSON(b []byte) error {
	*card = ParseCard(strings.Trim(string(b), "\""))
	return nil
}

func (card Card) String() string {
	b := make([]byte, 2)

	switch card.Suit() {
	case Club:
		b[0] = 'C'
	case Diamond:
		b[0] = 'D'
	case Heart:
		b[0] = 'H'
	case Spade:
		b[0] = 'S'
	default:
		return ""
	}

	switch card.Rank() {
	case Deuce:
		b[1] = '2'
	case Trey:
		b[1] = '3'
	case Four:
		b[1] = '4'
	case Five:
		b[1] = '5'
	case Six:
		b[1] = '6'
	case Seven:
		b[1] = '7'
	case Eight:
		b[1] = '8'
	case Nine:
		b[1] = '9'
	case Ten:
		b[1] = 'T'
	case Jack:
		b[1] = 'J'
	case Queen:
		b[1] = 'Q'
	case King:
		b[1] = 'K'
	case Ace:
		b[1] = 'A'
	default:
		return ""
	}

	return string(b)
}

type Deck struct {
	cards [NumCard]Card
	pos   int
}

func NewDeck() *Deck {
	deck := new(Deck)
	deck.Init()
	return deck
}

//
//   This routine initializes the deck.  A deck of cards is
//   simply an integer array of length 52 (no jokers).  This
//   array is populated with each card, using the following
//   scheme:
//
//   An integer is made up of four bytes.  The high-order
//   bytes are used to hold the rank bit pattern, whereas
//   the low-order bytes hold the suit/rank/prime value
//   of the card.
//
//   +--------+--------+--------+--------+
//   |xxxbbbbb|bbbbbbbb|cdhsrrrr|xxpppppp|
//   +--------+--------+--------+--------+
//
//   p = prime number of rank (deuce=2,trey=3,four=5,five=7,...,ace=41)
//   r = rank of card (deuce=0,trey=1,four=2,five=3,...,ace=12)
//   cdhs = suit of card
//   b = bit turned on depending on rank of card
//
func (deck *Deck) Init() {
	n := 0
	suit := 0x8000

	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			deck.cards[n] = Card(primes[j] | (j << 8) | suit | (1 << uint32(16+j)))
			n++
		}
		suit >>= 1
	}

	deck.pos = 0
}

func (deck *Deck) Find(rank, suit int) Card {
	for _, card := range deck.cards {
		if card.Rank() == rank && card.Suit() == suit {
			return card
		}
	}

	return NilCard
}

func (deck *Deck) Shuffle() {
	deck.pos = 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := r.Perm(52)
	for i, v := range a {
		deck.cards[i], deck.cards[v] = deck.cards[v], deck.cards[i]
	}
}

func (deck *Deck) Take() Card {
	if deck.pos >= NumCard {
		return NilCard
	}

	card := deck.cards[deck.pos]
	deck.pos++
	return card
}

func HandRank(val int) int {
	if val > 6185 { // 1277 high card
		return HgihCard
	}
	if val > 3325 { // 2860 one pair
		return OnePair
	}
	if val > 2467 { //  858 two pair
		return TwoPair
	}
	if val > 1609 { //  858 three-kind
		return ThreeOfAKind
	}
	if val > 1599 { //   10 straights
		return Straight
	}
	if val > 322 { // 1277 flushes
		return Flush
	}
	if val > 166 { //  156 full house
		return FullHouse
	}
	if val > 10 { //  156 four-kind
		return FourOfAKind
	}
	if val == 1 {
		return RoyalFlush
	}
	return StraightFlush //   10 straight-flushes
}

func eva5cards(cards [5]Card) int {
	q := int(cards[0]|cards[1]|cards[2]|cards[3]|cards[4]) >> 16

	// check for Flushes and StraightFlushes
	if cards[0]&cards[1]&cards[2]&cards[3]&cards[4]&SuitMask != 0 {
		return flushes[q]
	}

	// check for Straights and HighCard hands
	s := unique5[q]
	if s != 0 {
		return (s)
	}

	// let's do it the hard way
	q = cards[0].prime() * cards[1].prime() * cards[2].prime() * cards[3].prime() * cards[4].prime()
	q = find(q)

	return values[q]
}

// perform a binary search on a pre-sorted array
func find(key int) int {
	var low, high, mid int
	high = 4887

	for low <= high {
		mid = (high + low) >> 1 // divide by two
		if key < products[mid] {
			high = mid - 1
		} else if key > products[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}

	fmt.Fprintf(os.Stderr, "ERROR:  no match found; key = %d\n", key)
	return -1
}

var perm7 = [][5]int{
	{0, 1, 2, 3, 4},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 6},
	{0, 1, 2, 4, 5},
	{0, 1, 2, 4, 6},
	{0, 1, 2, 5, 6},
	{0, 1, 3, 4, 5},
	{0, 1, 3, 4, 6},
	{0, 1, 3, 5, 6},
	{0, 1, 4, 5, 6},
	{0, 2, 3, 4, 5},
	{0, 2, 3, 4, 6},
	{0, 2, 3, 5, 6},
	{0, 2, 4, 5, 6},
	{0, 3, 4, 5, 6},
	{1, 2, 3, 4, 5},
	{1, 2, 3, 4, 6},
	{1, 2, 3, 5, 6},
	{1, 2, 4, 5, 6},
	{1, 3, 4, 5, 6},
	{2, 3, 4, 5, 6},
}

func Eva7Hand(cards [7]Card) int {
	var hand [5]Card

	best := 0xFFFF

	for i, _ := range perm7 {
		for j, _ := range hand {
			hand[j] = cards[perm7[i][j]]
		}
		v := eva5cards(hand)
		if v < best {
			best = v
		}
	}

	return best
}
