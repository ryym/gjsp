package main

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input, 0}
}

func (l *Lexer) Position() int {
	return l.pos
}

func (l *Lexer) Range(start, end int) string {
	return l.input[start:end]
}

func (l *Lexer) PeekChar() byte {
	if l.pos >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.pos]
}

func (l *Lexer) ReadChar() byte {
	if l.pos >= len(l.input) {
		panic("cannot read after EOF")
	}
	l.pos += 1
	return l.PeekChar()
}

func (l *Lexer) ReadMany(n int) []byte {
	bs := make([]byte, n)
	start := l.pos
	for p := start; p < start+n; p++ {
		bs[p-start] = l.PeekChar()
		l.ReadChar()
	}
	return bs
}

func (l *Lexer) SkipWhitespaces() {
	ch := l.PeekChar()
	for ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		ch = l.ReadChar()
	}
}
