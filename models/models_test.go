package models

import (
	"slices"
	"testing"
)

func TestNewBlogPost(t *testing.T) {
	testCases := []struct {
		id          string
		title       string
		author      string
		publishDate string
		content     string
	}{
		{"1", "hello", "world", "2025-09-20", "hello world"},
		{"2", "test", "test-author", "2025-09-20", "this is a test"},
		{"3", "test-1", "test-author-1", "2025-09-20", "this is a test and even more..."},
	}
	for _, tc := range testCases {
		post := NewBlogPost(tc.id, tc.title, tc.publishDate, tc.author, tc.content)
		if tc.author != post.Author || tc.id != post.Id || tc.content != post.Content || tc.publishDate != post.PublishingDate {
			t.Errorf("Expecting BlogPost{Id: %s, Author: %s, Content: %s, PublishingDate: %s, Title: %s}, got BlogPost{Id: %s, Author: %s, Content: %s, PublishingDate: %s, Title: %s}", tc.id, tc.author, tc.content, tc.publishDate, tc.title, post.Id, post.Author, post.Content, post.PublishingDate, post.Title)
		}
	}
}

func TestSortBlogPosts(t *testing.T) {
	testCases := []struct {
		startingList []*BlogPost
		expected     []string
	}{
		{[]*BlogPost{NewBlogPost("1", "", "", "", ""), NewBlogPost("4", "", "", "", ""), NewBlogPost("2", "", "", "", ""), NewBlogPost("3", "", "", "", "")}, []string{"4", "3", "2", "1"}},
		{[]*BlogPost{NewBlogPost("5", "", "", "", ""), NewBlogPost("3", "", "", "", ""), NewBlogPost("2", "", "", "", ""), NewBlogPost("1", "", "", "", "")}, []string{"5", "3", "2", "1"}},
		{[]*BlogPost{NewBlogPost("5", "", "", "", ""), NewBlogPost("3", "", "", "", ""), NewBlogPost("1", "", "", "", ""), NewBlogPost("1", "", "", "", "")}, []string{"5", "3", "1", "1"}},
	}
	for _, tc := range testCases {
		blogs := SortBlogPosts(tc.startingList)
		blogIds := make([]string, len(blogs))
		for i, blog := range blogs {
			blogIds[i] = blog.Id
		}
		if !slices.Equal(blogIds, tc.expected) {
			t.Errorf("Expecting the blog IDs to be in the following order: %v, got %v", tc.expected, blogIds)
		}
	}
}
