package main

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
)

func Serve(g *Guide, w http.ResponseWriter, r *http.Request) error {
	var err error

	path := r.URL.Path
	slashPath := appendSlash(path)

	if path != slashPath {
		http.Redirect(w, r, slashPath, http.StatusMovedPermanently)
		return err
	}

	html, err := g.Render(path)
	if err != nil {
		return err
	}

	ct := mime.TypeByExtension(filepath.Ext(r.URL.Path))
	w.Header().Set("Content-Type", ct)

	fmt.Fprint(w, html)

	return err
}