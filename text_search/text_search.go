package textsearch

// import (
// 	"fmt"
// 	"regexp"
// 	"slices"
// 	"strings"

// 	"github.com/AstraBert/git-push-blog/models"
// 	"github.com/crawlab-team/bm25"
// )

// func createCorpus(posts []*models.BlogPost) []string {
// 	corpus := make([]string, len(posts)*2)
// 	for _, post := range posts {
// 		corpus = append(corpus, []string{post.Title, post.Content}...)
// 	}
// 	return corpus
// }

// func getTokenizer() func(s string) []string {
// 	re := regexp.MustCompile(`\n{2,}`)

// 	return func(s string) []string {
// 		return re.Split(s, -1)
// 	}
// }

// func GetTextSearcher(posts []*models.BlogPost) (*bm25.BM25L, error) {
// 	corpus := createCorpus(posts)

// 	tokenizer := getTokenizer()

// 	bm25, err := bm25.NewBM25L(corpus, tokenizer, 1.2, 0.75, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bm25, err
// }

// func SearchText(text string, searcher *bm25.BM25L, posts []*models.BlogPost) ([]*models.BlogPost, error) {
// 	tokenizedText := getTokenizer()(text)
// 	fmt.Println(tokenizedText)
// 	maxSimDocs, err := searcher.GetTopN(tokenizedText, 1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(maxSimDocs)
// 	matchedPosts := []*models.BlogPost{}
// 	for _, text := range maxSimDocs {
// 		for _, post := range posts {
// 			if strings.Contains(post.Content, text) || strings.Contains(post.Title, text) {
// 				if !slices.Contains(matchedPosts, post) {
// 					matchedPosts = append(matchedPosts, post)
// 					break
// 				}
// 			}
// 		}
// 	}
// 	return matchedPosts, nil
// }

import (
	"github.com/AstraBert/git-push-blog/models"
	"github.com/blevesearch/bleve"
)

func CreateIndex(posts []*models.BlogPost) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("posts.bleve", mapping)
	if err != nil {
		return nil, err
	}
	batch := index.NewBatch()
	for _, post := range posts {
		err := batch.Index(post.Title, post)
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
	req.Fields = []string{"title", "author", "content"}
	searchResult, err := index.Search(req)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}
