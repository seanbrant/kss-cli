build:
	go build -o _build/kss


compile-project:
	go get github.com/jteeuwen/go-bindata
	tar -cvf project.tar project
	go-bindata -out "project.go" -func DefaultProject project.tar
	go fmt project.go
	rm project.tar
