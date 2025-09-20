package models

import (
	"cmp"
	"slices"
)

type BlogPost struct {
	Id             int
	Title          string
	PublishingDate string
	Author         string
	Content        string
}

func NewBlogPost(id int, title, pubDate, author, content string) *BlogPost {
	return &BlogPost{Id: id, Title: title, PublishingDate: pubDate, Author: author, Content: content}
}

func SortBlogPosts(blogs []*BlogPost) []*BlogPost {
	slices.SortFunc(blogs,
		func(a, b *BlogPost) int {
			return cmp.Compare(b.Id, a.Id)
		})
	return blogs
}
