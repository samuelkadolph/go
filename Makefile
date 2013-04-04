PACKAGES = campfire httpstream nullable phidgets/raw phidgets
FORMATS = $(addprefix fmt/,$(PACKAGES))
TESTS = $(addprefix test/,$(PACKAGES))

all: fmt test

fmt: $(FORMATS)

test: $(TESTS)

.PHONY: all fmt test

include Makefile.common
