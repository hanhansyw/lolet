package main

import (
	"fmt"

	"github.com/hkail/lolet/internal/player"
)

func main() {
	fmt.Println("hello world")

	graph := player.GetGraph()
	players := graph.GetPlayersByTags([]player.GetPlayersByTagsArg{
		{
			TagType: 1,
			TagKeys: []int{1},
		},
		{
			TagType: 2,
			TagKeys: []int{1, 2},
		},
		{
			TagType: 3,
			TagKeys: []int{1, 2, 3, 4},
		},
		{
			TagType: 4,
			TagKeys: []int{1, 2},
		},
		{
			TagType: 5,
			TagKeys: []int{1},
		},
	})

	for _, p := range players {
		fmt.Printf("%#v\n", p)
	}
	fmt.Println(len(players))
}
