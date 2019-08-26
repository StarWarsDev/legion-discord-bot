package lookup

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
)

// Util exposes lookup functions
type Util struct {
	legionData *data.LegionData
}

// NewUtil creates a new Util pointer
func NewUtil(legionData *data.LegionData) Util {
	return Util{
		legionData: legionData,
	}
}
