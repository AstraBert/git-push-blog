package pagereader

import (
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func TestAddStyleToHTML(t *testing.T) {
	testCases := []struct {
		starting string
		expected string
	}{
		{
			starting: "<h1>Hello World</h1>",
			expected: `<h1 class="font-comic text-3xl text-gray-700 mb-4">Hello World</h1>`,
		},
		{
			starting: "<h2>Subtitle</h2>",
			expected: `<h2 class="font-poppins text-2xl text-gray-700 mb-3">Subtitle</h2>`,
		},
		{
			starting: "<h3>Section Header</h3>",
			expected: `<h3 class="font-bold text-xl text-black mb-2 font-poppins">Section Header</h3>`,
		},
		{
			starting: "<h4>Subsection</h4>",
			expected: `<h4 class="font-semibold text-lg text-gray-600 mb-2 font-poppins">Subsection</h4>`,
		},
		{
			starting: "<p>This is a paragraph.</p>",
			expected: `<p class="text-gray-700/90 leading-relaxed mb-4 font-sans">This is a paragraph.</p>`,
		},
		{
			starting: "<code>fmt.Println</code>",
			expected: `<code class="bg-pink-100/70 text-pink-800 px-2 py-1 rounded font-mono text-sm border border-pink-200">fmt.Println</code>`,
		},
		{
			starting: "<pre>func main() {\\n    fmt.Println(\"Hello\")\\n}</pre>",
			expected: `<pre class="bg-gray-900/90 text-green-400 p-4 rounded-lg overflow-x-auto font-mono text-sm border border-gray-700 shadow-lg mb-4 backdrop-blur-sm">func main() {\n    fmt.Println("Hello")\n}</pre>`,
		},
		{
			starting: "<h1>Title</h1><p>Content with <code>inline code</code> here.</p>",
			expected: `<h1 class="font-comic text-3xl text-gray-700 mb-4">Title</h1><p class="text-gray-700/90 leading-relaxed mb-4 font-sans">Content with <code class="bg-pink-100/70 text-pink-800 px-2 py-1 rounded font-mono text-sm border border-pink-200">inline code</code> here.</p>`,
		},
		{
			starting: "<h2>Code Example</h2><pre>package main\\n\\nimport \"fmt\"</pre><p>This shows Go syntax.</p>",
			expected: `<h2 class="font-poppins text-2xl text-gray-700 mb-3">Code Example</h2><pre class="bg-gray-900/90 text-green-400 p-4 rounded-lg overflow-x-auto font-mono text-sm border border-gray-700 shadow-lg mb-4 backdrop-blur-sm">package main\n\nimport "fmt"</pre><p class="text-gray-700/90 leading-relaxed mb-4 font-sans">This shows Go syntax.</p>`,
		},
		{
			starting: "<h1>Main</h1><h2>Sub</h2><h3>Section</h3><h4>Detail</h4>",
			expected: `<h1 class="font-comic text-3xl text-gray-700 mb-4">Main</h1><h2 class="font-poppins text-2xl text-gray-700 mb-3">Sub</h2><h3 class="font-bold text-xl text-black mb-2 font-poppins">Section</h3><h4 class="font-semibold text-lg text-gray-600 mb-2 font-poppins">Detail</h4>`,
		},
		{
			starting: "<p>Use <code>go run</code> to execute <code>main.go</code> files.</p>",
			expected: `<p class="text-gray-700/90 leading-relaxed mb-4 font-sans">Use <code class="bg-pink-100/70 text-pink-800 px-2 py-1 rounded font-mono text-sm border border-pink-200">go run</code> to execute <code class="bg-pink-100/70 text-pink-800 px-2 py-1 rounded font-mono text-sm border border-pink-200">main.go</code> files.</p>`,
		},
		{
			starting: "<div>No styling for div</div>",
			expected: "<div>No styling for div</div>",
		},
		{
			starting: "",
			expected: "",
		},
		{
			starting: "Plain text with no HTML tags",
			expected: "Plain text with no HTML tags",
		},
	}
	for _, tc := range testCases {
		if tc.expected != AddStyleToHTML(tc.starting) {
			t.Errorf("Expecting %s, got %s", tc.expected, AddStyleToHTML(tc.starting))
		}
	}
}

func TestGetMarkdownFiles(t *testing.T) {
	files, err := GetMarkdownFiles("../test_files")
	if err != nil {
		t.Errorf("Expecting no error, got %s", err.Error())
	}
	if !slices.Equal(files, []string{"../test_files/1.md", "../test_files/2.md"}) {
		t.Errorf("Expecting the files list to be %v, got %v", []string{"./test_files/1.md", "./test_files/2.md"}, files)
	}
}

func TestMarkdownToPost(t *testing.T) {
	files, _ := GetMarkdownFiles("../test_files")
	for _, fl := range files {
		blog, err := MarkdownToPost(fl)
		if err != nil {
			t.Errorf("Expecting no error while processing %s, got %s", fl, err.Error())
		}

		postId, _ := strconv.Atoi(strings.ReplaceAll(path.Base(fl), ".md", ""))
		if blog.Id != postId {
			t.Errorf("Expected blog.Id to be %d, got %d", postId, blog.Id)
		}
		if blog.Author != "Clelia Astra Bertelli" {
			t.Errorf("Expected blog.Author to be 'Clelia Astra Bertelli', got '%s'", blog.Author)
		}
		if !strings.HasPrefix(blog.Title, "Getting Started with") {
			t.Errorf("Expected blog.Title to start with 'Getting Started with', got '%s'", blog.Title)
		}

		flContent, _ := os.ReadFile(fl)
		sanitazedMd := strings.SplitN(string(flContent), "---", 3)[2]
		extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse([]byte(sanitazedMd))

		// create HTML renderer with extensions
		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		postContent := string(markdown.Render(doc, renderer))

		if AddStyleToHTML(postContent) != blog.Content {
			t.Errorf("Expecting post content to be %s, got %s", postContent, blog.Content)
		}
	}
}
