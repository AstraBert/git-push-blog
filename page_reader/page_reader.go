package pagereader

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/AstraBert/git-push-blog/models"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var HTMLStyleMap = map[string]string{
	"<h1>":   `<h1 class="font-comic text-3xl text-gray-700 mb-4">`,
	"<h2>":   `<h2 class="font-poppins text-2xl text-gray-700 mb-3">`,
	"<h3>":   `<h3 class="font-bold text-xl text-black mb-2 font-poppins">`,
	"<h4>":   `<h4 class="font-semibold text-lg text-gray-600 mb-2 font-poppins">`,
	"<p>":    `<p class="text-gray-700/90 leading-relaxed mb-4 font-sans">`,
	"<code>": `<code class="bg-pink-100/70 text-pink-800 px-2 py-1 rounded font-mono text-sm border border-pink-200">`,
	"<pre>":  `<pre class="bg-gray-900/90 text-green-400 p-4 rounded-lg overflow-x-auto font-mono text-sm border border-gray-700 shadow-lg mb-4 backdrop-blur-sm">`,
}

func AddStyleToHTML(content string) string {
	for k, v := range HTMLStyleMap {
		content = strings.ReplaceAll(content, k, v)
	}
	return content
}

func MarkdownToPost(filePath string) (*models.BlogPost, error) {
	postId := strings.ReplaceAll(path.Base(filePath), ".md", "")
	var postTitle string
	var postAuthor string
	var postDate string
	var postContent string
	md, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	header, sanitazedMd := strings.SplitN(string(md), "---", 3)[1], strings.SplitN(string(md), "---", 3)[2]
	for line := range strings.Lines(header) {
		switch {
		case strings.HasPrefix(line, "title:"):
			postTitle = strings.TrimSpace(strings.ReplaceAll(line, "title:", ""))
		case strings.HasPrefix(line, "author:"):
			postAuthor = strings.TrimSpace(strings.ReplaceAll(line, "author:", ""))
		case strings.HasPrefix(line, "publishing_date:"):
			postDate = strings.TrimSpace(strings.ReplaceAll(line, "publishing_date:", ""))
		default:
			continue
		}
	}
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(sanitazedMd))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	postContent = string(markdown.Render(doc, renderer))

	return models.NewBlogPost(postId, postTitle, postDate, postAuthor, AddStyleToHTML(postContent)), nil
}

func GetMarkdownFiles(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	mdFiles := []string{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			if _, errConv := strconv.Atoi(strings.ReplaceAll(file.Name(), ".md", "")); errConv != nil {
				continue
			}
			mdFiles = append(mdFiles, dirPath+"/"+file.Name())
		}
	}

	if len(mdFiles) == 0 {
		return nil, errors.New("no markdown file compliant with the format specifications has been found within the requested directory")
	}

	return mdFiles, nil
}
