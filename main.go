package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AstraBert/git-push-blog/models"
	pagereader "github.com/AstraBert/git-push-blog/page_reader"
	"github.com/AstraBert/git-push-blog/templates"
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
			if post.Id == postId {
				blogPost = post
				break
			}
		}

		templates.Post(blogPost).Render(r.Context(), w)
	})

	fmt.Println("Server started on :8000")
	http.ListenAndServe(":8000", nil)
}
