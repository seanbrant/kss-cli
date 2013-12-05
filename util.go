package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/knieriem/markdown"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func expand(path string) (string, error) {
	var err error

	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return path, err
		}

		path = filepath.Join(usr.HomeDir, strings.TrimPrefix(path, "~/"))
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return path, err
	}

	return path, err
}

func join(root string, paths []string) []string {
	for i, path := range paths {
		paths[i] = filepath.Join(root, path)
	}

	return paths
}

func appendSlash(path string) string {
	if hasExt(path) == false && strings.HasSuffix(path, "/") == false {
		return fmt.Sprintf("%s/", path)
	}

	return path
}

func splitExt(path string) (string, string) {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)], ext
}

func hasExt(path string) bool {
	return len(filepath.Ext(path)) != 0
}

func find(paths []string, name string) (string, error) {
	for _, path := range paths {
		fullPath := filepath.Join(path, name)

		found, err := exists(fullPath)
		if err != nil {
			return "", err
		}

		if found {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("No file named '%s' found", name)
}

func copyfile(src string, dest string) error {
	found, err := exists(dest)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.Chmod(dest, info.Mode())
	if err != nil {
		return err
	}

	return nil
}

func copydir(src string, dest string) error {
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if p == src {
			return nil
		}

		name := filepath.Join(dest, strings.Replace(p, src, "", 1))

		found, err := exists(name)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if found == false {
				err = os.MkdirAll(name, info.Mode())
				if err != nil {
					return err
				}
			}
		} else {
			err = copyfile(p, name)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func ensure(dir string, remove bool) error {
	var err error

	if remove {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}

	return err
}

func markdownify(data string) string {
	w := new(bytes.Buffer)
	p := markdown.NewParser(&markdown.Extensions{Smart: true})
	p.Markdown(strings.NewReader(data), markdown.ToHTML(w))
	return w.String()
}

func template(dirs []string, name string) (*mustache.Template, error) {
	path, err := find(dirs, name)
	if err != nil {
		return nil, err
	}

	return mustache.ParseFile(path)
}

func render(dirs []string, lname string, cname string, ctx ...interface{}) (string, error) {
	var err error

	l, err := template(dirs, lname)
	if err != nil {
		return "", err
	}

	c, err := template(dirs, cname)
	if err != nil {
		return "", err
	}

	context := map[string]interface{}{}

	// Note: We marshal and unmarshal JSON because I can't figure out
	// a better/easier/cleaner way to make the mustache renderer use
	// lower-case variables without using interfaces for everything.
	for _, v := range ctx {
		d, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		json.Unmarshal(d, &context)
	}

	content := c.Render(context)
	data := l.Render(context, map[string]string{"content": content})

	return data, nil
}
