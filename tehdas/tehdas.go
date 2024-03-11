package tehdas

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"
)

var Inst = NewTehdas()

const (
	key = iota
	val
)

func NewTehdas() *Tehdas {
	return &Tehdas{
		Values: make(map[string]string),
		order:  make(map[string]int),
	}
}

type Tehdas struct {
	Values map[string]string
	order  map[string]int
}

func (t *Tehdas) Add(k, v string) {
	t.Values[k] = v
	t.order[k] = len(t.Values) - 1
}
func (t *Tehdas) Del(k string) {
	delete(t.Values, k)
	delete(t.order, k)
}
func (t *Tehdas) MustGet(k string) string {
	v, ok := t.Values[k]
	if !ok {
		panic(fmt.Sprintf("key %s not found", k))
	}
	return v
}

func (t *Tehdas) Decode(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var i int
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.ReplaceAll(text, "\n", "")
		i++
		sar := strings.Split(text, "::")
		if len(sar) != 2 {
			return fmt.Errorf("kv pair on line %v is invalid", i)
		}
		if _, ok := t.Values[sar[key]]; ok {
			return fmt.Errorf("duplicate key on line %v", i)
		}
		k, err := url.QueryUnescape(sar[key])
		if err != nil {
			return fmt.Errorf("couldn't unescape key on line %v;	err: %s", i, err)
		}
		v, err := url.QueryUnescape(sar[val])
		if err != nil {
			return fmt.Errorf("couldn't unescape value on line %v;	err: %s", i, err)
		}
		t.Values[k] = v
		t.order[k] = i - 1
	}
	if i == 0 {
		return &ErrEmpty{}
	}
	return nil
}

type ErrEmpty struct{}

func (e *ErrEmpty) Error() string {
	return "no values found"
}

func IsEmptyErr(err error) bool {
	_, ok := err.(*ErrEmpty)
	return ok
}

func (t *Tehdas) Encode(w io.Writer) error {
	vals := make([]string, len(t.Values))
	for k, i := range t.order {
		vals[i] = fmt.Sprintf("%s::%s\n", url.QueryEscape(k), url.QueryEscape(t.Values[k]))
	}

	builder := strings.Builder{}
	for _, kv := range vals {
		builder.WriteString(kv)
	}
	_, err := w.Write([]byte(builder.String()))
	if err != nil {
		return fmt.Errorf("couldn't write encoded values to writer, err: %s", err)
	}
	return nil
}
