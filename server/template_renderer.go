package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/labstack/echo/v4"
)

// Recursively find an parse template files in a given directory.
// Thanks to: https://stackoverflow.com/a/50581032
func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

// A custom html/template renderer for the Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func timeAgo(t time.Time) string {
	now := time.Now()
	if t.After(now) {
		return t.Format("2006-01-02")
	}

	since := now.Sub(t)
	hours := math.Floor(since.Hours())
	if hours >= 24*7 {
		return t.Format("2006-01-02")
	}

	if hours >= 24 {
		days := math.Floor(hours / 24)
		s := ""
		if days > 1 {
			s = "s"
		}
		return fmt.Sprintf("%.f day%s ago", days, s)
	}

	minutes := math.Floor(since.Minutes())
	if minutes >= 60 {
		s := ""
		if hours > 1 {
			s = "s"
		}
		return fmt.Sprintf("%.f hour%s ago", hours, s)
	}

	seconds := since.Seconds()
	if seconds >= 60 {
		s := ""
		if minutes > 1 {
			s = "s"
		}
		return fmt.Sprintf("%.f minute%s ago", minutes, s)
	}

	return "now"
}

// Truncate a string to a given length.
// Thanks to: https://stackoverflow.com/a/59955803
func truncateString(str string, max int) string {
	if len(str) == max {
		return str
	}

	lastSpaceIx := -1
	len := 0
	for i, r := range str {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len >= max {
			if lastSpaceIx != -1 {
				return str[:lastSpaceIx] + "..."
			}
			// If here, string is longer than max, but has no spaces
		}
	}
	// If here, string is shorter than max

	return str
}

// Extracts out the base scheme and host name of a URL.
func urlBase(str string) string {
	u, err := url.Parse(str)
	if err != nil || u.Host == "" {
		return str
	}

	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Host)

	return fmt.Sprintf("%s://%s", scheme, host)
}

// Extracts out the base host name of a URL.
func urlHost(str string) string {
	u, err := url.Parse(str)
	if err != nil || u.Host == "" {
		return str
	}

	base := strings.ToLower(u.Host)
	if strings.HasPrefix(base, "www.") {
		return base[4:]
	}

	return base
}

// Setup the template renderer, adding it to the Echo server
func setupTemplateRenderer(e *echo.Echo) error {
	// Setup template functions
	funcMap := template.FuncMap{
		"dec": func(i int) int {
			return i - 1
		},
		"inc": func(i int) int {
			return i + 1
		},
		"timeago":  timeAgo,
		"truncate": truncateString,
		"urlbase":  urlBase,
		"urlhost":  urlHost,
	}

	// Find and parse template files
	rootTemplate, err := findAndParseTemplates("templates", funcMap)
	if err != nil {
		return err
	}

	// Setup renderer
	renderer := &TemplateRenderer{
		templates: rootTemplate,
	}
	e.Renderer = renderer

	return nil
}
