// Stefan Nilsson 2013-03-13

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	star   = "Pythia"
	venue  = "Delphi"
	prompt = "> "
)

var (
	answers = [][]string{
		{
			"I am afraid you lost me entirely in the sea of perplexity.",
			"The moon of two worlds apart is the only matter sustaining balance in the system.",
			"Lana says she is pretty when she cries. One may wonder what that implies.",
		},
		{
			"That is definitely a question worth pursuing.",
			"I did not quite catch that.",
			"My uncle once said something similar. He was later diagnosed with hysteria.",
			"Yes, as a matter of fact, it is about the only logical thing around here.",
			"I do not believe so. Confabulating on such a matter as it might upset the King.",
			"That remains to be seen, however Bjorndir in Winterhold might know about it.",
		},
		{
			"You think? It sounds like you are in over your head.",
			"Well, I am nothing more than a " + venue + ", I would not know anything about that.",
			"That is definitely a delicately profound way of expressing it.",
			"That I can agree with, it may have that kind of effect.",
		},
	}
)

func main() {
	fmt.Printf("Welcome to %s, the oracle at %s.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")

	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%s heard: %s\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	questions := make(chan string)
	answers := make(chan string)

	// Answer questions.
	go func() {
		for q := range questions {
			go prophecy(q, answers)
		}
	}()

	// Make prophecies.
	go func() {
		for {
			dur := time.Duration(20+rand.Intn(20)) * time.Second
			time.Sleep(dur)
			prophecy("", answers)
		}
	}()

	// Print answers.
	go func() {
		for a := range answers {
			fmt.Printf("\n%s said: ", star)
			for _, c := range a {
				fmt.Printf("%c", c)
				time.Sleep(time.Millisecond * 25)
			}
			fmt.Printf("\n%s", prompt)
		}
	}()

	return questions
}

// This is the oracle's secret algorithm.
// It waits for a while and then sends a message on the answer channel.
func prophecy(question string, answer chan<- string) {
	// Keep them waiting. Pythia, the original oracle at Delphi,
	// only gave prophecies on the seventh day of each month.
	time.Sleep(time.Duration(2+rand.Intn(10)) * time.Second)

	// Match questions and statements
	patternQuestion := ".*\\?"
	patternStatement := ".*\\."
	collection := 0 // Default to generic prophecies
	var match bool
	var _ error

	// If the question is an actual question, use the question set
	match, _ = regexp.MatchString(patternQuestion, question)

	if match {
		collection = 1
	}

	// If the question is a statement, use the statement set
	match, _ = regexp.MatchString(patternStatement, question)
	if match {
		collection = 2
	}

	// Obtain the set of appropriate responses and pick a message
	set := answers[collection]
	msg := set[rand.Intn(len(set))]

	// Format and send the message
	answer <- msg
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}
