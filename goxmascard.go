package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

type Card struct {
	from string
	to string
	subject string
	message string
}

type Friend struct {
	name string
	email string
}

func postman(name string, in <-chan Card) chan string {
	comm := make(chan string) // communication chan
	go func() {
		for card := range in{
			comm <- fmt.Sprintf("POSTMAN %s, CARD: %+v", name, card)
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
		}
		close(comm)
	}()

	return comm
}

func sendCards(friends []Friend) chan Card {
	del := make(chan Card) // delivery chan

	go func() {
		for i := range friends {
			f := friends[i]
			card := Card{
				"Best Friend <hello@yourbestfriend.co.uk>",
				fmt.Sprintf("%s <%s>", f.name, f.email),
				"hohoho",
				"Have a merry little christmas!",
			}
			del <- card
		}
		close(del)
	}()

	return del
}

func merge(cs ...<-chan string) <-chan string{
	out := make(chan string)
	var wg sync.WaitGroup

	output := func(c <-chan string) {
		for m := range c{
			out <- m
		}
		wg.Done()
	}
	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {

	friends := []Friend{
		Friend{"Adrian", "adrian@email.com"},
		Friend{"Leo", "leo@dicaprio.com"},
		Friend{"Tom", "tom@m.com"},
		Friend{"Simon", "simon@example.com"},
		Friend{"Boris", "contact@boris.com"},
		Friend{"Donald", "donald@duck.com"},
	}

	delChan := sendCards(friends)

	c1 := postman("Pat", delChan)
	c2 := postman("Bob", delChan)

	for msg := range merge(c1, c2) {
		fmt.Println(msg)
	}
}
