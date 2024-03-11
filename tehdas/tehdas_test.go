package tehdas

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func makekv(K, V string) string {
	k := url.QueryEscape(K)
	v := url.QueryEscape(V)
	return fmt.Sprintf("%s::%s\n", k, v)
}

func TestTehdas_Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid input",
			input:   makekv("key1:key2:key3", "value1\nkey2::value2"),
			wantErr: false,
		},
		{
			name:    "valid escaped input",
			input:   makekv("key1", "value1\nkey2::value2"),
			wantErr: false,
		},
		{
			name:    "invalid input format",
			input:   "key1\nvalue:2::value1\nkey1::value2",
			wantErr: true,
			errMsg:  "kv pair on line 1 is invalid",
		},
		{
			name:    "duplicate key",
			input:   "key1::value1\nkey1::value2",
			wantErr: true,
			errMsg:  "duplicate key on line 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTehdas()
			err := tr.Decode(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Decode() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTehdas_Encode(t *testing.T) {
	tr := NewTehdas()
	tr.Add("key1", "value1")
	tr.Add("key2", "value2")

	buffer := &bytes.Buffer{}
	err := tr.Encode(buffer)
	if err != nil {
		t.Errorf("Encode() error = %v", err)
	}

	expected := "key1::value1\nkey2::value2\n"
	if got := buffer.String(); got != expected {
		t.Errorf("Encode() got = %v, want %v", got, expected)
	}
}

func TestTehdas_EncodeDecode_RoundTrip(t *testing.T) {
	input := "key1::value1\nkey2::value2"
	tr := NewTehdas()
	err := tr.Decode(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Decode() failed: %v", err)
	}

	buffer := &bytes.Buffer{}
	err = tr.Encode(buffer)
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	// Comparing the encoded output to the initial input after query escaping
	expected, _ := url.QueryUnescape(input)
	if got := buffer.String(); got != expected+"\n" {
		t.Errorf("Encode()/Decode() round trip got = %v, want %v", got, expected)
	}
}
