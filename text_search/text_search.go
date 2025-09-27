package textsearch

import (
	"errors"

	"github.com/AstraBert/git-push-blog/models"
	"github.com/blevesearch/bleve/v2"
)

func CreateIndex(posts []*models.BlogPost, path string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}
	batch := index.NewBatch()
	for _, post := range posts {
		err := batch.Index(post.Id, post)
		if err != nil {
			continue
		}
	}
	errBatch := index.Batch(batch)

	if errBatch != nil {
		return nil, errBatch
	}

	return index, nil
}

func SearchText(text string, index bleve.Index) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(text)
	req := bleve.NewSearchRequest(query)
	req.Explain = false
	req.Fields = []string{"title", "author", "publishing_date"}
	searchResult, err := index.Search(req)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func ParseSearchResults(searchResults *bleve.SearchResult) ([]*models.BlogPost, error) {
	var posts []*models.BlogPost

	for _, hit := range searchResults.Hits {
		post := &models.BlogPost{}

		// Access fields directly from the hit
		if title, ok := hit.Fields["title"].(string); ok {
			post.Title = title
		} else {
			continue
		}

		if author, ok := hit.Fields["author"].(string); ok {
			post.Author = author
		} else {
			continue
		}

		if publishingDate, ok := hit.Fields["publishing_date"].(string); ok {
			post.PublishingDate = publishingDate
		} else {
			continue
		}

		// You can also get the document ID if you stored it
		post.Id = hit.ID

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, errors.New("no posts matched the search")
	}

	return posts, nil
}
