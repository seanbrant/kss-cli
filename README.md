kss-cli
=======

Commands
--------
* kss create path/to/styleguide                         - Creates a default style guide project
* kss serve path/to/styleguide/config.yml [addr[:port]] - Serves your style guide project with a built in web server
* kss build path/to/styleguide/config.yml               - Builds a static version of your style guide


Config
------
* build_dir     - Output location for the build command. Relative to config file.
* example_dirs  - Array of dictionaries containing your example html files. Relative to config file.
* source_dirs   - Array of dictionaries containing your source (css/less/scss/sass) files. Relative to config file.
* static_dirs   - Array of dictionaries containing your static assets required to render the style guide. Relative to config file.
* static_root   - Optional location of where to copy your static directories. Relative to config file. Default static, relative to build_dir.
* static_url    - Optional url to serve the static assets from outputted in the rendered html. Default static.
* template_dirs - Array of dictionaries containing your mustache templates used to render the html.
* template_ext  - Optional extension for your mustache templates. Default .mustache
* page_ext      - Optional extension for the style guide pages in the rendered html. If empty paths will end in a trailing stash like (/buttons/). Default "".
