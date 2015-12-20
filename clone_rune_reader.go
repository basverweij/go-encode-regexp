package encre

import (
	"io"
)

type CloningRuneReader struct {
	b []rune
	r io.RuneReader
}

func NewCloningRuneReader(r io.RuneReader) *CloningRuneReader {
	return &CloningRuneReader{b: make([]rune, 0), r: r}
}

func (c *CloningRuneReader) ReadRune() (r rune, size int, err error) {
	r, size, err = c.r.ReadRune()

	c.b = append(c.b, r)

	return
}

func (c *CloningRuneReader) Slice(fromIncl, toExcl int) string {
	return string(c.b[fromIncl:toExcl])
}
