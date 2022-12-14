package player

type Player struct {
	Name          string `csv:"name" json:"name"`
	Lane          string `csv:"lane" json:"lane"`
	Region        string `csv:"region" json:"region"`
	Team          string `csv:"team" json:"team"`
	FirstFeature  string `csv:"first_feature" json:"first_feature"`
	SecondFeature string `csv:"second_feature" json:"second_feature"`
	Level         string `csv:"level" json:"level"`

	id         int
	tagIndexes map[tagIndex]struct{}
}

func (p Player) GetValueByTagType(tagType TagType) string {
	switch tagType {
	case TagTypeLane:
		return p.Lane
	case TagTypeRegion:
		return p.Region
	case TagTypeTeam:
		return p.Team
	case TagTypeFirstFeature:
		return p.FirstFeature
	case TagTypeSecondFeature:
		return p.SecondFeature
	case TagTypeLevel:
		return p.Level
	default:
		return ""
	}
}
