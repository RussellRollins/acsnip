package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if err := inner(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func inner() error {
	if len(os.Args) < 3 {
		return errors.New("insufficient command line arguments, please specify the file to be converted and time to start snipping")
	}
	if len(os.Args) > 3 {
		return fmt.Errorf("unexpected command line arguments [%s]", strings.Join(os.Args[3:], " "))
	}

	name := os.Args[1]
	input, err := os.Open(name)
	if err != nil {
		return err
	}
	defer input.Close()

	noExtension := strings.TrimSuffix(input.Name(), filepath.Ext(input.Name()))
	outName := fmt.Sprintf("%s.snip.cast", noExtension)

	output, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer output.Close()

	stampStr := os.Args[2]
	first, err := strconv.ParseFloat(stampStr, 64)
	if err != nil {
		return err
	}

	stampRe := regexp.MustCompile(`(?m)^\[(?P<stamp>[\d\.]*)`)
	replaceRe := regexp.MustCompile(`(?m)^\[[\d\.]*`)

	scanner := bufio.NewScanner(input)
	found := false
	last := float64(0)
	diff := float64(0)
	for scanner.Scan() {
		text := scanner.Text()
		match := stampRe.FindStringSubmatch(text)
		// Just send anything not matching the timestamp format to the writer.
		if len(match) == 0 {
			output.WriteString(fmt.Sprintf("%s\n", text))
			continue
		}
		result := make(map[string]string)
		for i, name := range stampRe.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		this, err := strconv.ParseFloat(result["stamp"], 64)
		if err != nil {
			return err
		}

		if first == this {
			found = true
			diff = this - last
		}

		if found {
			now := this - diff

			text = replaceRe.ReplaceAllString(text, fmt.Sprintf("%f", now))

			// This is an even grosser hack than the rest of this, maybe the regex should just find the number?
			text = fmt.Sprintf("[%s", text)
		}

		output.WriteString(fmt.Sprintf("%s\n", text))
		last = this
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("did not find an entry at timestamp %s", stampStr)
	}

	return nil
}
