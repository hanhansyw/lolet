package player

type TagType = int

const (
	TagTypeUnknown TagType = iota
	TagTypeLane
	TagTypeRegion
	TagTypeTeam
	TagTypeFirstFeature
	TagTypeSecondFeature
	TagTypeLevel
)

var tagTypes = []TagType{
	TagTypeLane,
	TagTypeRegion,
	TagTypeTeam,
	TagTypeFirstFeature,
	TagTypeSecondFeature,
	TagTypeLevel,
}
