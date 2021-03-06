package main

import (
	"log"

	loco "github.com/lucasmdrs/go-loco/v2"
)

func main() {
	g := loco.Init()

	if err := g.AddProject("KEY_1", "project1"); err != nil {
		log.Fatalln(err.Error())
	}

	if err := g.AddProject("KEY_2", "project2"); err != nil {
		log.Fatalln(err.Error())
	}

	ch, err := g.StartPoller()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		select {
		case <-ch:
			log.Println("changed!")
		}
	}

}
