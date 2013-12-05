package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"path/filepath"
)

type Config struct {
	BuildDir     string   `build_dir`
	ExampleDirs  []string `example_dirs`
	RootDir      string
	SourceDirs   []string `source_dirs`
	StaticDirs   []string `static_dirs`
	StaticRoot   string   `static_root`
	StaticUrl    string   `static_url`
	TemplateDirs []string `template_dirs`
	TemplateExt  string   `template_ext`
	PageExt      string   `page_ext`
}

func NewConfig(name string) (c *Config, err error) {
	c = &Config{}

	name, err = expand(name)
	if err != nil {
		return
	}

	found, err := exists(name)
	if found == false {
		err = fmt.Errorf("configuration file not found at path '%s'", name)
		return
	}

	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}

	c.RootDir = filepath.Dir(name)
	c.StaticRoot = "static"
	c.StaticUrl = "/static"
	c.TemplateExt = ".mustache"
	c.PageExt = ""

	goyaml.Unmarshal(bytes, &c)

	c.BuildDir = filepath.Join(c.RootDir, c.BuildDir)
	c.ExampleDirs = join(c.RootDir, c.ExampleDirs)
	c.SourceDirs = join(c.RootDir, c.SourceDirs)
	c.StaticDirs = join(c.RootDir, c.StaticDirs)
	c.TemplateDirs = join(c.RootDir, c.TemplateDirs)

	return
}
