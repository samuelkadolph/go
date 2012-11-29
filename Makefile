REPO = github.com/samuelkadolph/go
PACKAGES = campfire

SYMLINK = $(GOPATH)/src/$(REPO)
FORMATS = $(addprefix fmt/$(REPO)/,$(PACKAGES))

all: ls

fmt: $(FORMATS)

fmt/%: ln
	go fmt $*

ln: $(SYMLINK)

$(SYMLINK):
	mkdir -p "$(dir $@)"
	ln -fs "$(CURDIR)" "$@"
