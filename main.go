package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/AstraBert/git-push-blog/commons"
	"github.com/AstraBert/git-push-blog/models"
	pagereader "github.com/AstraBert/git-push-blog/page_reader"
	"github.com/AstraBert/git-push-blog/templates"
	textsearch "github.com/AstraBert/git-push-blog/text_search"
	"github.com/a-h/templ"
)

func main() {
	mdFiles, errFls := pagereader.GetMarkdownFiles("./contents")
	if errFls != nil {
		fmt.Printf("Error while getting markdown files: %s\n", errFls.Error())
		return
	}
	blogs := []*models.BlogPost{}
	for _, fl := range mdFiles {
		post, errPost := pagereader.MarkdownToPost(fl)
		if errPost != nil {
			fmt.Printf("Error while getting the content for file %s: %s\n", fl, errPost.Error())
			continue
		}
		blogs = append(blogs, post)
	}
	if len(blogs) == 0 {
		fmt.Println("Could not retrieve the content of any post, exiting...")
		return
	}
	blogs = models.SortBlogPosts(blogs)
	if commons.PathExists("posts.bleve") {
		os.RemoveAll("posts.bleve")
	}
	index, err := textsearch.CreateIndex(blogs, "posts.bleve")
	if err != nil {
		fmt.Println("An error occurred while creating the index: " + err.Error())
		return
	}

	homeComponent := templates.Home()
	blogComponent := templates.BlogPage(blogs)

	http.Handle("/", templ.Handler(homeComponent))
	http.Handle("/blog", templ.Handler(blogComponent))
	http.HandleFunc("/blog/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		postId, errConv := strconv.Atoi(id)
		if errConv != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		if postId-1 < 0 || postId-1 >= len(blogs) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		var blogPost *models.BlogPost

		for _, post := range blogs {
			if post.Id == id {
				blogPost = post
				break
			}
		}

		templates.Post(blogPost).Render(r.Context(), w)
	})

	http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		searchText := r.FormValue("search")
		if searchText == "" {
			http.Error(w, "Search text is required", http.StatusBadRequest)
			return
		}

		results, err := textsearch.SearchText(searchText, index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		searchResults, errS := textsearch.ParseSearchResults(results)

		if errS != nil {
			http.Error(w, "No matches where found for the search", http.StatusInternalServerError)
			return
		}

		// Return HTML instead of JSON for HTMX
		w.Header().Set("Content-Type", "text/html")

		// Generate HTML for search results
		html := ""
		for _, post := range searchResults {
			html += fmt.Sprintf(`<li class="text-gray-600 font-sans text-xl">• <a href="/blog/%s"><span class="text-pink-500 underline">%s</span></a> posted by <span class="font-semibold text-pink-400">%s</span> on %s</li>`,
				post.Id, post.Title, post.Author, post.PublishingDate)
		}

		if html == "" {
			html = `<li class="text-gray-600 font-sans text-xl">No results found</li>`
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	fmt.Println("Server started on :8000")
	http.ListenAndServe(":8000", nil)
}
