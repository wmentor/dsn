package dsn

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type DSN map[string]string

var (
	mapper    map[rune]rune
	mapperRev map[rune]string

	ErrParse            error
	ErrKeyAlreadyExists error
)

func init() {

	ErrParse = errors.New("parse error")
	ErrKeyAlreadyExists = errors.New("key already exists")

	mapper = map[rune]rune{
		'\\': '\\',
		'\'': '\'',
		'"':  '"',
		's':  ' ',
		't':  '\t',
		'r':  '\r',
		'n':  '\n',
		' ':  ' ',
		'=':  '=',
	}

	mapperRev = map[rune]string{
		' ':  "\\s",
		'\t': "\\t",
		'\n': "\\n",
		'\r': "\\r",
		'"':  "\\\"",
		'\'': "\\'",
		'\\': "\\\\",
		'=':  "\\=",
	}
}

func encode(str string) string {

	if str == "" {
		return str
	}

	found := false

	for _, r := range str {

		if _, h := mapperRev[r]; h {
			found = true
			break
		}

	}

	if !found {
		return str
	}

	maker := strings.Builder{}

	for _, r := range str {

		if nv, h := mapperRev[r]; h {
			maker.WriteString(nv)
		} else {
			maker.WriteRune(r)
		}

	}

	return maker.String()
}

func New(str string) (DSN, error) {

	res := make(map[string]string)

	k := strings.Builder{}
	v := strings.Builder{}

	mode := 0

	for _, r := range str {

		switch mode {

		case 0:

			if unicode.IsSpace(r) {

			} else if r == '\\' {
				mode = 1
				k.Reset()
			} else if r == '=' {
				k.Reset()
				mode = 3
			} else {
				k.Reset()
				k.WriteRune(r)
				mode = 2
			}

		case 1:

			if c, h := mapper[r]; h {
				k.WriteRune(c)
				mode = 2
			} else {
				return nil, ErrParse
			}

		case 2:

			if r == '=' {
				mode = 3
			} else if r == '\\' {
				mode = 1
			} else if unicode.IsSpace(r) {
				return nil, ErrParse
			} else {
				k.WriteRune(r)
			}

		case 3:

			v.Reset()

			if unicode.IsSpace(r) {
				key := k.String()
				if _, h := res[key]; h {
					return nil, ErrKeyAlreadyExists
				}
				res[key] = ""
				mode = 0
			} else if r == '=' {
				return nil, ErrParse
			} else if r == '\\' {
				mode = 4
			} else {
				v.WriteRune(r)
				mode = 5
			}

		case 4:

			if c, h := mapper[r]; h {
				v.WriteRune(c)
				mode = 5
			} else {
				return nil, ErrParse
			}

		case 5:

			if unicode.IsSpace(r) {
				key := k.String()
				if _, h := res[key]; h {
					return nil, ErrKeyAlreadyExists
				}
				res[key] = v.String()
				mode = 0
			} else if r == '=' {
				return nil, ErrParse
			} else if r == '\\' {
				mode = 4
			} else {
				v.WriteRune(r)

			}

		}

	}

	switch mode {

	case 1, 2, 4:
		return nil, ErrParse

	case 3, 5:

		if mode == 3 {
			v.Reset()
		}

		key := k.String()
		if _, h := res[key]; h {
			return nil, ErrKeyAlreadyExists
		}
		res[key] = v.String()

	}

	return res, nil
}

func (dsn DSN) GetString(key string, defval string) string {

	if v, h := dsn[key]; h {
		return v
	}

	return defval
}

func (dsn DSN) SetString(key string, val string) {
	dsn[key] = val
}

func (dsn DSN) GetInt64(key string, defval int64) int64 {

	if v, h := dsn[key]; h {

		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
	}

	return defval
}

func (dsn DSN) SetInt64(key string, val int64) {
	dsn[key] = strconv.FormatInt(val, 10)
}

func (dsn DSN) GetInt(key string, defval int) int {

	if v, h := dsn[key]; h {

		if val, err := strconv.Atoi(v); err == nil {
			return val
		}
	}

	return defval
}

func (dsn DSN) SetInt(key string, val int) {
	dsn[key] = strconv.Itoa(val)
}

func (dsn DSN) GetBool(key string, defval bool) bool {

	if v, h := dsn[key]; h {

		if val, err := strconv.ParseBool(v); err == nil {
			return val
		}
	}

	return defval
}

func (dsn DSN) SetBool(key string, val bool) {
	dsn[key] = strconv.FormatBool(val)
}

func (dsn DSN) GetFloat(key string, defval float64) float64 {

	if v, h := dsn[key]; h {

		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
	}

	return defval
}

func (dsn DSN) SetFloat(key string, val float64) {
	dsn[key] = fmt.Sprintf("%.6f", val)
}

func (dsn DSN) String() string {

	notFirst := false

	maker := strings.Builder{}

	for k, v := range dsn {

		if notFirst {
			maker.WriteRune(' ')
		}

		maker.WriteString(encode(k))
		maker.WriteRune('=')
		maker.WriteString(encode(v))

		notFirst = true
	}

	return maker.String()
}
