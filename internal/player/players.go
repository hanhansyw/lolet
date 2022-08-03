package player

import (
	"math"
	"os"

	"github.com/gocarina/gocsv"
)

const (
	defaultPlayersDataFilePath = "data/players.csv"
	defaultTagGroupSize        = 3
)

var graphInstance *Graph

type tagIndex uint8

type Graph struct {
	tagToPlayersMapper    map[tagIndex][]*Player
	tagIndexToValueMapper map[tagIndex]string
	tagTypeToTagsMapper   map[TagType][]Tag
}

type Tag struct {
	Type  TagType `json:"type"`
	Key   int     `json:"key"`
	Value string  `json:"value"`
}

func (t Tag) index() tagIndex {
	return tagIndex(t.Type<<4 + t.Key)
}

type GetPlayersByTagsArg struct {
	TagType TagType
	TagKeys []int
}

type ComputeResult struct {
	Tags    []Tag    `json:"tags"`
	Players []Player `json:"players"`
}

func init() {
	newGraph, err := loadPlayerGraph(defaultPlayersDataFilePath)
	if err != nil {
		panic(err)
	}

	graphInstance = newGraph

}

func GetGraph() *Graph {
	return graphInstance
}

func loadPlayerGraph(csvFilePath string) (*Graph, error) {
	playersFile, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer playersFile.Close()

	players := make([]*Player, 0)
	if err = gocsv.UnmarshalFile(playersFile, &players); err != nil {
		return nil, err
	}

	tagValueToKeyMapper := make(map[TagType]map[string]int)
	tagToPlayersMapper := make(map[tagIndex][]*Player)
	tagIndexToValueMapper := make(map[tagIndex]string)
	tagTypeToTagsMapper := make(map[TagType][]Tag)
	for i, player := range players {
		player.tagIndexes = make(map[tagIndex]struct{})
		player.id = i + 1

		for _, tagType := range tagTypes {
			if _, ok := tagValueToKeyMapper[tagType]; !ok {
				tagValueToKeyMapper[tagType] = make(map[string]int)
			}

			tagValue := player.GetValueByTagType(tagType)
			tagKey, ok := tagValueToKeyMapper[tagType][tagValue]
			if !ok {
				tagKey = len(tagValueToKeyMapper[tagType]) + 1
				tagValueToKeyMapper[tagType][tagValue] = tagKey
			}

			tag := Tag{
				Type:  tagType,
				Key:   tagKey,
				Value: tagValue,
			}
			if !ok {
				tagTypeToTagsMapper[tagType] = append(tagTypeToTagsMapper[tagType], tag)
				tagIndexToValueMapper[tag.index()] = tagValue
			}

			players[i].tagIndexes[tag.index()] = struct{}{}
			tagToPlayersMapper[tag.index()] = append(tagToPlayersMapper[tag.index()], players[i])
		}
	}

	return &Graph{
		tagToPlayersMapper:    tagToPlayersMapper,
		tagIndexToValueMapper: tagIndexToValueMapper,
		tagTypeToTagsMapper:   tagTypeToTagsMapper,
	}, nil
}

func (g *Graph) ComputePlayersByTags(args []GetPlayersByTagsArg) []ComputeResult {
	tags := make([]Tag, 0)
	for _, arg := range args {
		for _, tagKey := range arg.TagKeys {
			tag := Tag{
				Type: arg.TagType,
				Key:  tagKey,
			}
			tag.Value = g.tagIndexToValueMapper[tag.index()]
			tags = append(tags, tag)
		}
	}

	computeResults := make([]ComputeResult, 0)
	var f func(i int, currentTags []Tag)
	f = func(i int, currentTags []Tag) {
		if len(currentTags) > 1 {
			currentTagSize := len(currentTags)
			if currentTags[currentTagSize-1].Type == currentTags[currentTagSize-2].Type {
				return
			}
		}

		if len(currentTags) == defaultTagGroupSize {
			targetPlayers := g.getPlayersByTags(currentTags)
			if len(targetPlayers) > 0 {
				computeResults = append(computeResults, ComputeResult{
					Tags:    currentTags,
					Players: targetPlayers,
				})
			}

			return
		}

		for j := i; j < len(tags); j++ {
			f(j+1, append(currentTags, tags[j]))
		}
	}
	f(0, nil)

	return computeResults
}

func (g *Graph) GetTagsMapper() map[TagType][]Tag {
	return g.tagTypeToTagsMapper
}

func (g *Graph) getPlayersByTags(tags []Tag) []Player {
	if len(tags) == 0 {
		return nil
	}

	minSizeTagIndex := -1
	minSize := math.MaxInt64
	for i, tag := range tags {
		tagPlayerSize := len(g.tagToPlayersMapper[tag.index()])
		if tagPlayerSize < minSize {
			minSize = tagPlayerSize
			minSizeTagIndex = i
		}
	}

	tags[0], tags[minSizeTagIndex] = tags[minSizeTagIndex], tags[0]
	basePlayers := g.tagToPlayersMapper[tags[0].index()]
	tags = tags[1:]

	players := make([]Player, 0)
	for _, player := range basePlayers {
		hasTags := true
		for _, tag := range tags {
			if _, ok := player.tagIndexes[tag.index()]; !ok {
				hasTags = false
				break
			}
		}

		if hasTags {
			players = append(players, *player)
		}
	}

	return players
}
