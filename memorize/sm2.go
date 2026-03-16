package memorize

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Card struct {
	Word Word

	EF   float64
	Rep  int
	Intv int
}

type Deck struct {
	DeckQueue []Card `json:"queue"`
	cur       *Card
}

func NewDeck(list WordList) *Deck {
	d := &Deck{}

	cfg := GetConfig()
	efInit := 2.5
	intvInit := 1
	if cfg != nil {
		efInit = cfg.SM2.EFInitial
		intvInit = cfg.SM2.IntervalFirst
	}

	for _, w := range list.Words {
		d.DeckQueue = append(d.DeckQueue, Card{
			Word: w,
			EF:   efInit,
			Rep:  0,
			Intv: intvInit,
		})
	}

	d.Shuffle()

	return d
}

func (d *Deck) Next() *Word {
	if len(d.DeckQueue) == 0 {
		return nil
	}

	c := d.DeckQueue[0]
	d.DeckQueue = d.DeckQueue[1:]
	d.cur = &c

	return &c.Word
}

func (d *Deck) Queue() []Card {
	return d.DeckQueue
}

func (d *Deck) Shuffle() {
	for i := len(d.DeckQueue) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.DeckQueue[i], d.DeckQueue[j] = d.DeckQueue[j], d.DeckQueue[i]
	}
}

func (d *Deck) Review(q int) {
	if d.cur == nil {
		return
	}

	cfg := GetConfig()
	intvSecond := 6
	efMin := 1.3
	efCoefA := 0.1
	efCoefB := 0.08
	efCoefC := 0.02
	if cfg != nil {
		intvSecond = cfg.SM2.IntervalSecond
		efMin = cfg.SM2.EFMin
		efCoefA = cfg.SM2.EFCoefA
		efCoefB = cfg.SM2.EFCoefB
		efCoefC = cfg.SM2.EFCoefC
	}

	c := *d.cur

	if q < 3 {
		c.Rep = 0
		c.Intv = 1
	} else {
		c.Rep++

		if c.Rep == 1 {
			c.Intv = 1
		} else if c.Rep == 2 {
			c.Intv = intvSecond
		} else {
			c.Intv = int(math.Round(float64(c.Intv) * c.EF))
		}
	}

	c.EF += efCoefA - float64(5-q)*(efCoefB+float64(5-q)*efCoefC)
	if c.EF < efMin {
		c.EF = efMin
	}

	idx := c.Intv
	if idx > len(d.DeckQueue) {
		idx = len(d.DeckQueue)
	}

	d.DeckQueue = append(d.DeckQueue[:idx],
		append([]Card{c}, d.DeckQueue[idx:]...)...)

	d.cur = nil
}

func (d *Deck) SaveDeckCache(path string) error {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal deck: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write deck cache: %w", err)
	}

	return nil
}

func LoadDeckCache(path string) (*Deck, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read deck cache: %w", err)
	}

	var d Deck
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck cache: %w", err)
	}

	return &d, nil
}
