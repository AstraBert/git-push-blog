package models

import (
	"cmp"
	"slices"
	"strconv"
)

type BlogPost struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	PublishingDate string `json:"publishing_date"`
	Author         string `json:"author"`
	Content        string `json:"content"`
}

func NewBlogPost(id string, title, pubDate, author, content string) *BlogPost {
	return &BlogPost{Id: id, Title: title, PublishingDate: pubDate, Author: author, Content: content}
}

func SortBlogPosts(blogs []*BlogPost) []*BlogPost {
	slices.SortFunc(blogs,
		func(a, b *BlogPost) int {
			bId, _ := strconv.Atoi(b.Id)
			aId, _ := strconv.Atoi(a.Id)
			return cmp.Compare(bId, aId)
		})
	return blogs
}
