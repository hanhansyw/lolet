package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hkail/lolet/internal/player"
)

func main() {
	r := gin.Default()

	r.GET("/get_tags", func(c *gin.Context) {
		graph := player.GetGraph()
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"data": graph.GetTagsMapper(),
		})
	})

	r.GET("/compute", func(c *gin.Context) {
		type Req struct {
			LaneKeys          string `form:"lane_keys"`
			RegionKeys        string `form:"region_keys"`
			TeamsKeys         string `form:"team_keys"`
			FirstFeatureKeys  string `form:"first_feature_keys"`
			SecondFeatureKeys string `form:"second_feature_keys"`
			LevelKeys         string `form:"level_keys"`
		}
		var req Req
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Abort()
			return
		}

		keysStrSlice := []string{
			req.LaneKeys,
			req.RegionKeys,
			req.TeamsKeys,
			req.FirstFeatureKeys,
			req.SecondFeatureKeys,
			req.LevelKeys,
		}
		tagTypes := []player.TagType{
			player.TagTypeLane,
			player.TagTypeRegion,
			player.TagTypeTeam,
			player.TagTypeFirstFeature,
			player.TagTypeSecondFeature,
			player.TagTypeLevel,
		}
		args, err := generateGetPlayersByTagsArgs(keysStrSlice, tagTypes)
		if err != nil {
			c.Abort()
			return
		}

		graph := player.GetGraph()
		players := graph.ComputePlayersByTags(args)

		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"data": players,
		})
	})

	r.StaticFile("home", "./templates/index.html")

	log.Println(r.Run(":8888"))
}

func generateGetPlayersByTagsArgs(keysStrSlice []string, tagTypes []player.TagType) ([]player.GetPlayersByTagsArg, error) {
	var args []player.GetPlayersByTagsArg
	for i, keysStr := range keysStrSlice {
		if len(keysStr) == 0 {
			continue
		}

		keys, err := parseKeysStrToIntSlice(keysStr)
		if err != nil {
			return nil, err
		}

		args = append(args, player.GetPlayersByTagsArg{
			TagType: tagTypes[i],
			TagKeys: keys,
		})
	}

	return args, nil
}

func parseKeysStrToIntSlice(keysStr string) ([]int, error) {
	keyStrSlice := strings.Split(keysStr, ",")
	keys := make([]int, 0, len(keyStrSlice))
	for _, keyStr := range keyStrSlice {
		key, err := strconv.Atoi(keyStr)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}
