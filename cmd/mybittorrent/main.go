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

func decode(b string) (interface{}, error) {
	if len(b) == 0 {
		return nil, io.ErrUnexpectedEOF
	}

	switch {
	case b[0] >= '0' && b[0] <= '9':
		return decodeString(b)
	default:
		return nil, fmt.Errorf("unexpected value type: %q", b[0])
	}
}

func decodeString(b string) (string, error) {
	var l int

	i := 0
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		l = l*10 + (int(b[i]) - '0')
		i++
	}

	if i == len(b) || b[i] != ':' {
		return "", fmt.Errorf("bad string")
	}

	i++

	if i+l > len(b) {
		return "", fmt.Errorf("bad string: out of bounds")
	}

	x := b[i : i+l]
	//	i += l

	return x, nil
}

func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		x, err := decode(os.Args[2])
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
