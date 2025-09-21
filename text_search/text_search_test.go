package textsearch

import (
	"os"
	"slices"
	"testing"

	"github.com/AstraBert/git-push-blog/commons"
	"github.com/AstraBert/git-push-blog/models"
	pagereader "github.com/AstraBert/git-push-blog/page_reader"
)

func getTestPosts() []*models.BlogPost {
	files, _ := pagereader.GetMarkdownFiles("../test_files")
	posts := make([]*models.BlogPost, len(files))
	for i, file := range files {
		posts[i], _ = pagereader.MarkdownToPost(file)
	}
	return posts
}

func TestCreateIndex(t *testing.T) {
	posts := getTestPosts()
	if commons.PathExists("tests.bleve") {
		os.RemoveAll("tests.bleve")
	}
	_, err := CreateIndex(posts, "tests.bleve")
	if err != nil {
		t.Errorf("Expecting the index creation to yield no error, but it produced: %s", err.Error())
	}
	if !commons.PathExists("tests.bleve") {
		t.Errorf("Expecting index creation to generate a tests.bleve folder, but it does not exists")
	}
}

func TestSearchText(t *testing.T) {
	posts := getTestPosts()
	if commons.PathExists("tests.bleve") {
		os.RemoveAll("tests.bleve")
	}
	index, _ := CreateIndex(posts, "tests.bleve")
	searchResult, err := SearchText("level two header", index)
	if err != nil {
		t.Errorf("Expecting the search operation to not yield any error, but it produced: %s", err.Error())
	}
	if searchResult.Hits.Len() == 0 {
		t.Error("Expecting the search operation to retrieve at least one match, but got 0")
	}
	ids := make([]string, searchResult.Hits.Len())
	for _, hit := range searchResult.Hits {
		ids = append(ids, hit.ID)
	}
	if !slices.Contains(ids, "1") {
		t.Error("Expecting the search operation to retrieve document 1.md as match, but it did not.")
	}
}

func TestParseSearchResults(t *testing.T) {
	posts := getTestPosts()
	if commons.PathExists("tests.bleve") {
		os.RemoveAll("tests.bleve")
	}
	index, _ := CreateIndex(posts, "tests.bleve")
	searchResult, _ := SearchText("Clelia Astra Bertelli", index)
	blogs, err := ParseSearchResults(searchResult)
	if err != nil {
		t.Errorf("Expecting the search operation to not yield any error, but it produced: %s", err.Error())
	}
	if len(blogs) < 2 {
		t.Errorf("Expecting the search operation to produce 2 matches, but it produced: %d", len(blogs))
	}
	ids := make([]string, len(blogs))
	for _, post := range blogs {
		ids = append(ids, post.Id)
	}
	if !slices.Contains(ids, "1") || !slices.Contains(ids, "2") {
		t.Error("Expecting the search operation to retrieve documents 1.md and 2.md as match, but it did not.")
	}
	searchResult1, _ := SearchText("rgargargamdbdns", index)
	_, err1 := ParseSearchResults(searchResult1)
	if err1 == nil {
		t.Error("Expecting the parsing operation to throw an error, but it did not")
	}
	if err1.Error() != "no posts matched the search" {
		t.Errorf("Expecting the parsing error message to be %s, but got %s", "no posts matched the search", err1.Error())
	}
}
