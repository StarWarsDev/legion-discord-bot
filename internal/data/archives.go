package data

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

// Keyword is a Legion keyword
type Keyword struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

// NewArchivesClient creates and returns a new ArchivesClient
func NewArchivesClient(url string) ArchivesClient {
	gqlClient := graphql.NewClient(url, nil)
	return ArchivesClient{
		gqlClient: gqlClient,
	}
}
