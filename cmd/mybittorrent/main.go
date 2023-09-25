package main

import (
	// Uncomment this line to pass the first stage
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"os"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func decode(b string, st int) (interface{}, int, error) {
	if st == len(b) {
		return nil, st, io.ErrUnexpectedEOF
	}

	switch {
	case b[st] == 'i':
		return decodeInt(b, st)
	case b[st] == 'l':
		return decodeList(b, st)
	case b[st] == 'd':
		return decodeDict(b, st)
	case b[st] >= '0' && b[st] <= '9':
		return decodeString(b, st)
	default:
		return nil, st, fmt.Errorf("unexpected value type: %q", b[st])
	}
}

func decodeInt(b string, st int) (int, int, error) {
	i := st + 1

	if i == len(b) {
		return 0, st, io.ErrUnexpectedEOF
	}

	neg := false

	if b[i] == '-' {
		neg = true
		i++
	}

	var x int

	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		x = x*10 + (int(b[i]) - '0')
		i++
	}

	if i == len(b) || b[i] != 'e' {
		return 0, st, fmt.Errorf("bad int")
	}

	i++

	if neg {
		x = -x
	}

	return x, i, nil
}

func decodeList(b string, st int) ([]interface{}, int, error) {
	i := st + 1

	var err error
	l := make([]interface{}, 0, 4)

	for {
		if i == len(b) {
			return nil, st, io.ErrUnexpectedEOF
		}

		if b[i] == 'e' {
			break
		}

		var x interface{}

		x, i, err = decode(b, i)
		if err != nil {
			return nil, i, err
		}

		l = append(l, x)
	}

	return l, i, nil
}

func decodeDict(b string, st int) (map[string]interface{}, int, error) {
	i := st + 1

	var err error
	d := make(map[string]interface{}, 4)

	for {
		if i == len(b) {
			return nil, st, io.ErrUnexpectedEOF
		}

		if b[i] == 'e' {
			break
		}

		var key, val interface{}

		key, i, err = decode(b, i)
		if err != nil {
			return nil, i, err
		}

		keys, ok := key.(string)
		if !ok {
			return nil, i, fmt.Errorf("dict key is not a string")
		}

		val, i, err = decode(b, i)
		if err != nil {
			return nil, i, err
		}

		d[keys] = val
	}

	return d, i, nil
}

func decodeString(b string, st int) (string, int, error) {
	var l int

	i := st
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		l = l*10 + (int(b[i]) - '0')
		i++
	}

	if i == len(b) || b[i] != ':' {
		return "", st, fmt.Errorf("bad string")
	}

	i++

	if i+l > len(b) {
		return "", st, fmt.Errorf("bad string: out of bounds")
	}

	x := b[i : i+l]
	i += l

	return x, i, nil
}

func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		x, _, err := decode(os.Args[2], 0)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		y, err := json.Marshal(x)
		if err != nil {
			fmt.Printf("error: marshal into json: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", y)
	default:
		fmt.Printf("Unknown command: %v\n", command)
		os.Exit(1)
	}
}
