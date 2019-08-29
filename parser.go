package main

// See https://www.json.org/ for the details of JSON spec.

import (
	"errors"
	"fmt"
	"strconv"
)

var escapeChars map[byte]bool

func init() {
	chrs := []byte{'"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u'}
	escapeChars = make(map[byte]bool, len(chrs))
	for _, ch := range chrs {
		escapeChars[ch] = true
	}
}

func Parse(l *Lexer) (interface{}, error) {
	return parseValue(l)
}

func parseValue(l *Lexer) (v interface{}, err error) {
	l.SkipWhitespaces()

	ch := l.PeekChar()
	switch ch {
	case 0:
		break
	case '"':
		v, err = parseString(l)
	case '[':
		v, err = parseArray(l)
	case '{':
		v, err = parseObject(l)
	case 't':
		v, err = parseTrue(l)
	case 'f':
		v, err = parseFalse(l)
	case 'n':
		v, err = parseNull(l)
	default:
		if ch == '-' || isDigit(ch) {
			v, err = parseNumber(l)
		} else {
			err = fmt.Errorf("unexpected character %s\n", string(ch))
		}
	}

	l.SkipWhitespaces()
	return v, err
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func parseObject(l *Lexer) (obj map[string]interface{}, err error) {
	l.ReadChar() // Left brace

	obj = make(map[string]interface{}, 0)

	for {
		l.SkipWhitespaces()
		key, err := parseString(l)
		if err != nil {
			return nil, err
		}

		l.SkipWhitespaces()

		ch := l.PeekChar()
		if ch != ':' {
			return nil, fmt.Errorf("invalid object: expected :, got %s\n", string(ch))
		}

		l.ReadChar()
		v, err := parseValue(l)
		if err != nil {
			return nil, err
		}

		obj[key] = v

		ch = l.PeekChar()
		switch ch {
		case ',':
			l.ReadChar()
			continue
		case '}':
			l.ReadChar()
			return obj, nil
		default:
			return nil, fmt.Errorf("invalid object: expected }, got %s\n", string(ch))
		}
	}

}

func parseArray(l *Lexer) (vs []interface{}, err error) {
	l.ReadChar() // Left bracket

	vs = make([]interface{}, 0)
	for {
		v, err := parseValue(l)
		if err != nil {
			return vs, err
		}

		vs = append(vs, v)

		ch := l.PeekChar()
		switch ch {
		case ',':
			l.ReadChar()
			continue
		case ']':
			l.ReadChar()
			return vs, nil
		default:
			return nil, fmt.Errorf("invalid array: expected ], got %s\n", string(ch))
		}
	}
}

func parseString(l *Lexer) (string, error) {
	l.ReadChar() // Quote

	start := l.Position()
	end := start
	closed := true
	for ch := l.PeekChar(); ch != '"'; {
		if ch == 0 {
			closed = false
			break
		}
		if ch == '\\' {
			esc := l.ReadChar()
			end += 1
			if !escapeChars[esc] {
				return "", fmt.Errorf("invalid string: unknown escape sequence %v\n", esc)
			}
			if esc == 'u' {
				return "", errors.New("invalid string: Sorry, unicode escape sequence does not be supported yet")
			}
		}
		ch = l.ReadChar()
		end += 1
	}

	if !closed {
		return "", errors.New("invalid string: no closing quote")
	}
	l.ReadChar() // Quote

	return l.Range(start, end), nil
}

func parseTrue(l *Lexer) (b bool, err error) {
	token := string(l.ReadMany(4))
	if token != "true" {
		return b, fmt.Errorf("expected true, got %s", token)
	}
	return true, nil
}

func parseFalse(l *Lexer) (b bool, err error) {
	token := string(l.ReadMany(5))
	if token != "false" {
		return b, fmt.Errorf("expected false, got %s", token)
	}
	return false, nil
}

func parseNull(l *Lexer) (v interface{}, err error) {
	token := string(l.ReadMany(4))
	if token != "null" {
		return v, fmt.Errorf("expected null, got %s", token)
	}
	return v, nil
}

func parseNumber(l *Lexer) (v interface{}, err error) {
	start := l.Position()
	end := start + 1
	fraction := false
	exponent := false
	for ; ; end++ {
		ch := l.ReadChar()

		if isDigit(ch) {
			continue
		}
		if ch == '.' {
			if fraction {
				return v, errors.New("invalid number: double dot")
			}
			fraction = true
			next := l.ReadChar()
			end += 1
			if isDigit(next) {
				continue
			} else {
				return v, errors.New("invalid number: dot must be followed by digit")
			}
		}
		if ch == 'e' || ch == 'E' {
			if !fraction {
				return v, errors.New("invalid number: exponent before fraction")
			}
			if exponent {
				return v, errors.New("invalid number: double exponent")
			}

			exponent = true
			sign := l.ReadChar()
			end += 1
			if sign != '-' && sign != '+' {
				return v, errors.New("invalid number: e/E must be followed by +/-")
			}

			digit := l.ReadChar()
			end += 1
			if !isDigit(digit) {
				return v, errors.New("invalid number: exponent must be followed by digit")
			}
			continue
		}

		break
	}

	num := l.Range(start, end)

	head := 0
	if num[0] == '-' {
		head += 1
	}
	if len(num) > head+1 {
		if num[head] == '0' && isDigit(num[head+1]) {
			return v, errors.New("invalid number: 0 start")
		}
	}

	if fraction {
		return strconv.ParseFloat(num, 64)
	} else {
		return strconv.Atoi(num)
	}
}
