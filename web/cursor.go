package web

type Cursor struct {
	current  int
	wildcard byte
	Points   []int
	Ruler    string
}

func NewCursor(ruler string, wildcard byte) *Cursor {
	ruler = TrimByte(ruler, wildcard)
	return &Cursor{Ruler: ruler, Points: make([]int, 0, len(ruler))}
}

func (c *Cursor) Reset() {
	c.current = 0
}

func (c *Cursor) End() {
	c.current = 0
}

func (c *Cursor) Last() {

}

func (c *Cursor) Next() {

}
