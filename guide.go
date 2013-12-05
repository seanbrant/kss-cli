package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/seanbrant/kss-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const indexFilename = "index"

type Modifier struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ClassName   string `json:"class_name"`
	Example     string `json:"example"`
}

func (m *Modifier) AddExample(t *mustache.Template) {
	m.Example = t.Render(map[string]string{
		"modifier_class": fmt.Sprintf(" %s", m.ClassName),
	})
}

type Section struct {
	Filename       string      `json:"filename"`
	Reference      string      `json:"reference"`
	Description    string      `json:"description"`
	Modifiers      []*Modifier `json:"modifiers"`
	ModifiersCount int         `json:"modifiers_count"`
	Example        string      `json:"example"`
}

func (s *Section) AddExample(t *mustache.Template) {
	s.Example = t.Render(map[string]string{})

	for _, m := range s.Modifiers {
		m.AddExample(t)
	}
}

func NewSection(s *kss.Section) *Section {
	modifiers := make([]*Modifier, 0, len(s.Modifiers))

	for _, m := range s.Modifiers {
		modifiers = append(modifiers, &Modifier{
			Name:        m.Name,
			Description: m.Description,
			ClassName:   m.ClassName(),
		})
	}

	return &Section{
		Filename:       s.Filename,
		Reference:      s.Reference,
		Description:    markdownify(s.Description),
		Modifiers:      modifiers,
		ModifiersCount: len(modifiers),
	}
}

type Page struct {
	Config        *Config    `json:-`
	Filename      string     `json:"filename"`
	Name          string     `json:"name"`
	Sections      []*Section `json:"sections"`
	SectionsCount int        `json:"sections_count"`
	template      string
}

func (p *Page) Url() string {
	if len(p.Config.PageExt) > 0 {
		return fmt.Sprintf("/%s%s", p.Filename, p.Config.PageExt)
	}

	if p.Filename == indexFilename {
		return "/"
	} else {
		return fmt.Sprintf("/%s/", p.Filename)
	}
}

func (p *Page) AddSection(section *kss.Section) {
	p.Sections = append(p.Sections, NewSection(section))
}

func (p *Page) Render(ctx ...interface{}) (string, error) {
	lname := fmt.Sprintf("layout%s", p.Config.TemplateExt)
	cname := fmt.Sprintf("%s%s", p.template, p.Config.TemplateExt)

	ctx = append(ctx, map[string]string{
		"static_url": p.Config.StaticUrl,
	})
	ctx = append(ctx, p)

	return render(p.Config.TemplateDirs, lname, cname, ctx...)
}

type Guide struct {
	Config *Config          `json:-`
	Pages  map[string]*Page `json:"pages"`
}

func (g *Guide) AddPageSection(name string, section *kss.Section) *Page {
	p := g.AddPage(name, name, "styleguide")
	p.AddSection(section)
	return p
}

func (g *Guide) AddPage(filename, name, template string) *Page {
	filename = strings.ToLower(filename)

	if g.Pages[filename] == nil {
		g.Pages[filename] = &Page{
			Config:   g.Config,
			Filename: filename,
			Name:     name,
			template: template,
		}
	}

	return g.Pages[filename]
}

func (g *Guide) LoadExamples(paths []string) error {
	examples := make(map[string]*mustache.Template)

	for _, parent := range paths {
		filepath.Walk(parent, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			t, _ := mustache.ParseFile(p)
			name, _ := splitExt(info.Name())
			examples[strings.ToLower(name)] = t

			return nil
		})
	}

	for _, p := range g.Pages {
		for _, s := range p.Sections {
			key := strings.ToLower(s.Reference)

			if example, ok := examples[key]; ok {
				s.AddExample(example)
			}
		}
	}

	return nil
}

func (g *Guide) Nav(active *Page) map[string]interface{} {
	number := 0
	items := make([]map[string]interface{}, 0, len(g.Pages)+1)

	for _, p := range g.Pages {
		items = append(items, map[string]interface{}{
			"name":      p.Name,
			"url":       p.Url(),
			"number":    number,
			"is_active": (p == active),
		})

		number += 1
	}

	return map[string]interface{}{
		"nav": items,
	}
}

func (g *Guide) RenderPage(path string) (string, error) {
	path = strings.TrimRight(path, g.Config.PageExt)

	p := g.Pages[strings.ToLower(strings.Trim(path, "/"))]
	if p == nil {
		return "", fmt.Errorf("page '%s' not found", path)
	}

	return p.Render(g.Nav(p))
}

func (g *Guide) RenderStatic(path string) (string, error) {
	path = strings.Replace(path, g.Config.StaticUrl, "", 1)
	path, err := find(g.Config.StaticDirs, path)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (g *Guide) Render(path string) (string, error) {
	if strings.HasPrefix(path, g.Config.StaticUrl) {
		return g.RenderStatic(path)
	} else {
		return g.RenderPage(path)
	}
}

func (g *Guide) RenderAll() (map[string]string, error) {
	var err error
	rendered := map[string]string{}

	for k, _ := range g.Pages {
		rendered[k], err = g.Render(k)
		if err != nil {
			return rendered, err
		}
	}

	return rendered, err
}

func NewGuide(c *Config) (g *Guide, err error) {
	g = &Guide{
		Config: c,
		Pages:  make(map[string]*Page),
	}

	g.AddPage(indexFilename, "Overview", "index")

	for reference, section := range kss.Parser(c.SourceDirs...) {
		g.AddPageSection(strings.Split(reference, ".")[0], section)
	}

	g.LoadExamples(c.ExampleDirs)

	return
}
