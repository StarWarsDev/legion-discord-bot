package search

import (
	"fmt"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/blevesearch/bleve"
	"github.com/bwmarrin/discordgo"
)

const (
	// IndexKey is the string name of the index
	IndexKey = "legion.data"
)

// Util is a struct used for searching LegionData
type Util struct {
	legionData *data.LegionData
	lookupUtil *lookup.Util
	index      bleve.Index
}

// NewUtil creates a new Util pointer
func NewUtil(legionData *data.LegionData, lookupUtil *lookup.Util) *Util {
	// index the data for searching
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(IndexKey, mapping)
	if err != nil {
		if err == bleve.ErrorIndexPathExists {
			fmt.Println("Found an existing index, reusing it...")
			index, err = bleve.Open(IndexKey)

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	util := &Util{
		legionData: legionData,
		lookupUtil: lookupUtil,
		index:      index,
	}
	util.indexCommandCards()
	util.indexUpgrades()
	util.indexUnits()

	return util
}

func (util *Util) indexCommandCards() {
	for _, card := range util.legionData.CommandCards {
		util.index.Index("commandcard."+card.LDF, card)
	}
}

func (util *Util) indexUpgrades() {
	for _, card := range util.legionData.Upgrades.Flattened() {
		util.index.Index("upgrade."+card.LDF, card)
	}
}

func (util *Util) indexUnits() {
	for _, card := range util.legionData.Units.Flattened() {
		util.index.Index("unit."+card.LDF, card)
	}
}

// FullSearch performs a full text search against all legion data
func (util *Util) FullSearch(text string) []*discordgo.MessageEmbed {
	results := []*discordgo.MessageEmbed{}
	query := bleve.NewMatchQuery(text)
	search := bleve.NewSearchRequest(query)
	searchResults, err := util.index.Search(search)
	if err != nil {
		fmt.Println(err)

		// return a message to the user that there was an error
		results = append(results, output.Error("No results", "There was an error searching for "+text))
		return results
	}

	for _, hit := range searchResults.Hits {
		id := strings.Split(hit.ID, ".")
		hitType := id[0]
		ldf := id[1]

		switch hitType {
		case "commandcard":
			card := util.lookupUtil.LookupCommandCardByLdf(ldf)
			results = append(results, output.CommandCard(card))
		case "unit":
			card := util.lookupUtil.LookupUnitByLdf(ldf)
			if len(card.CommandCards) > 0 {
				var commandCards []string
				for _, ldf := range card.CommandCards {
					card := util.lookupUtil.LookupCommandCardByLdf(ldf)
					if card != nil {
						commandCards = append(commandCards, card.Name)
					}
				}
				card.CommandCards = commandCards
			}
			results = append(results, output.Unit(card))
		case "upgrade":
			card := util.lookupUtil.LookupUpgradeByLdf(ldf)
			results = append(results, output.Upgrade(card))
		}
	}

	if searchResults.Total == 0 {
		results = append(
			results,
			output.Error(
				utils.WithTemplate("0 matches found for \"%s\"", text),
				utils.WithTemplate("Took %v", searchResults.Took),
			),
		)
	} else {
		results = append(
			results,
			output.Info(
				utils.WithTemplate("%d matches found for \"%s\"", len(results), text),
				utils.WithTemplate("Took %v", searchResults.Took),
			),
		)
	}

	return results
}
