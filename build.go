package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func CollectStatic(c *Config) error {
	var err error

	dest := filepath.Join(c.BuildDir, "static")

	for _, src := range c.StaticDirs {
		err = copydir(src, dest)
		if err != nil {
			return err
		}
	}

	return err
}

func Build(c *Config, g *Guide) error {
	var err error

	err = ensure(c.BuildDir, true)
	if err != nil {
		return err
	}

	pages, err := g.RenderAll()
	if err != nil {
		return err
	}

	for k, h := range pages {
		name := strings.ToLower(k)

		if name == "/" {
			name = "index"
		}

		path := filepath.Join(
			c.BuildDir,
			fmt.Sprintf("%s.html", strings.TrimRight(name, "/")),
		)

		err = ioutil.WriteFile(path, []byte(h), 0644)
		if err != nil {
			return err
		}
	}

	err = CollectStatic(c)
	if err != nil {
		return err
	}

	return err
}
