package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/shurcooL/graphql"
)

// ArchivesClient handles communication to the Archives
type ArchivesClient struct {
	gqlClient *graphql.Client
}

// GetKeywords gets the keywords from the archives
func (client *ArchivesClient) GetKeywords(field, term string) []Keyword {
	var query struct {
		Keywords []Keyword `graphql:"keywords(query: $query)"`
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("%s:%s", field, term)),
	}

	fmt.Println(variables)

	err := client.gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err)
	}

	return query.Keywords
}

// GetCommandCards get the command cards from the archives
func (client *ArchivesClient) GetCommandCards(field, term string) []CommandCard {
	var query struct {
		Commands []CommandCard `graphql:"embed(query: $query)"`
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("%s:%s", field, term)),
	}

	err := client.gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		j, _ := json.Marshal(&query)
		log.Println(err, string(j), variables)
	}

	return query.Commands
}

// GetUnits gets the unit cards from the archives
func (client *ArchivesClient) GetUnits(field, term string) []Unit {
	var query struct {
		Units []Unit `graphql:"units(query: $query)"`
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("%s:%s", field, term)),
	}

	err := client.gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		j, _ := json.Marshal(&query)
		log.Println(err, string(j), variables)
	}

	return query.Units
}

// GetUpgrades get the upgrade cards from the archives
func (client *ArchivesClient) GetUpgrades(field, term string) []Upgrade {
	var query struct {
		Upgrades []Upgrade `graphql:"upgrades(query: $query)"`
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("%s:%s", field, term)),
	}

	err := client.gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		j, _ := json.Marshal(&query)
		log.Println(err, string(j), variables)
	}

	return query.Upgrades
}

// NewArchivesClient creates and returns a new ArchivesClient
func NewArchivesClient(url string) ArchivesClient {
	gqlClient := graphql.NewClient(url, nil)
	return ArchivesClient{
		gqlClient: gqlClient,
	}
}
