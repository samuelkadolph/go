GO = go

REPO = github.com/samuelkadolph/go
PACKAGES = campfire httpstream mpg123 nullable phidgets phidgets/raw

FORMATS = $(addprefix fmt/,$(PACKAGES))
SYMLINK = $(GOPATH)/src/$(REPO)
TESTS = $(addprefix test/,$(PACKAGES))

all: fmt test

fmt: $(FORMATS)

fmt/%: ln
	$(GO) fmt $(REPO)/$*

ln: $(SYMLINK)

test: $(TESTS)

test/%: %
	$(GO) test $(REPO)/$*

$(SYMLINK):
	mkdir -p "$(dir $@)"
	ln -fs "$(CURDIR)" "$@"

.PHONY: all fmt ln test
