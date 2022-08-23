package structures

import "encoding/json"

type CardPile struct {
	head   *Card
	tail   *Card
	length int
}

func (p *CardPile) Len() int {
	return p.length
}

func (p *CardPile) Add(values ...string) {
	for _, v := range values {
		c := Card{value: v}

		if p.Len() == 0 {
			p.head = &c
			p.tail = &c
		} else {
			p.tail.next = &c
			p.tail = &c
		}
		p.length++
	}
}

func (p *CardPile) Remove() *Card {
	c := p.head
	if c != nil {
		p.head = p.head.next
		p.length--

		if c == p.tail {
			p.tail = nil
		}
	}

	return c
}

func (p *CardPile) Peek() *Card {
	return p.head
}

func (p *CardPile) ToSlice() []string {
	list := make([]string, 0, p.Len())

	n := p.head

	for n != nil {
		list = append(list, n.Value())
		n = n.next
	}

	return list
}

func (p *CardPile) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToSlice())
}

type Card struct {
	value string
	next  *Card
}

func (c Card) Value() string {
	return c.value
}
