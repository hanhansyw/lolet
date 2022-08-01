package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type Player struct {
	Name   string `csv:"name"`
	Lane   string `csv:"lane"`
	Region string `csv:"region"`
	Team   string `csv:"team"`
	Tag1   string `csv:"tag1"`
	Tag2   string `csv:"tag2"`
	Level  string `csv:"level"`
}

func main() {
	fmt.Println("hello world")
	players, err := loadPlayers()
	if err != nil {
		panic(err)
	}

	for _, player := range players {
		fmt.Printf("%#v\n", player)
	}
}

func loadPlayers() ([]Player, error) {
	playersFile, err := os.OpenFile("data/players.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer playersFile.Close()

	players := make([]Player, 0)
	if err = gocsv.UnmarshalFile(playersFile, &players); err != nil {
		return nil, err
	}

	return players, nil
}
