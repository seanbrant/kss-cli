package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Creator struct {
	Dir     string
	Project []byte
}

func (c *Creator) Copy(name string, reader io.Reader) error {
	var err error

	if hasExt(name) == false {
		err = os.MkdirAll(name, 0775)
		if err != nil {
			return err
		}
	} else {
		out, err := os.Create(name)
		defer out.Close()

		_, err = io.Copy(out, reader)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *Creator) Create() error {
	var err error

	r := bytes.NewReader(DefaultProject())
	tr := tar.NewReader(r)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// Strip the project directory from the name
		name := filepath.Join(c.Dir, strings.TrimLeft(hdr.Name, "project"))

		// Exclude root directory and dot files
		if name == c.Dir || strings.HasPrefix(filepath.Base(name), ".") {
			continue
		}

		c.Copy(name, tr)
	}

	return err
}

func Create(dst string) error {
	var err error

	dst, err = expand(dst)
	if err != nil {
		return err
	}

	found, err := exists(dst)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("Destination directory '%s' already exist.", dst)
	}

	creator := &Creator{Dir: dst, Project: DefaultProject()}

	err = ensure(creator.Dir, false)
	if err != nil {
		return err
	}

	return creator.Create()
}
