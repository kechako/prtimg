package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

func isTmux() bool {
	return os.Getenv("TMUX") != ""
}

func printOSC() {
	if isTmux() {
		fmt.Print("\x1bPtmux;\x1b\x1b]1337;File=")
	} else {
		fmt.Print("\x1b]1337;File=")
	}
}

func printST() {
	fmt.Print("\a")
	if isTmux() {
		fmt.Print("\x1b\\")
	}
}

func printImage(r io.Reader, name string) error {
	printOSC()
	defer printST()

	// options
	fmt.Printf("name=%s;inline=1:", base64.StdEncoding.EncodeToString([]byte(name)))

	// encode image
	w := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	defer w.Close()

	// copy to stdout
	_, err := io.Copy(w, r)

	return err
}

func _main() (int, error) {
	args := os.Args[1:]

	fname := "Unknown Image"
	var r io.Reader
	if len(args) < 1 {
		r = os.Stdin
	} else {
		fname = args[0]
		_, err := os.Stat(fname)
		if err != nil {
			return 1, errors.New(fmt.Sprintf("%s is not found.", fname))
		}

		file, err := os.Open(fname)
		if err != nil {
			return 1, errors.Wrap(err, fmt.Sprintf("Could not open %s.", fname))
		}
		defer file.Close()
		r = file
	}

	err := printImage(r, fname)
	if err != nil {
		return 2, errors.Wrap(err, "Could not print image.")
	}

	fmt.Println()

	return 0, nil
}

func main() {
	code, err := _main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
	}
	os.Exit(code)
}
