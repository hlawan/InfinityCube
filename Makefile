# Copyright (c) 2011-2012 Justus Winter <4winter@informatik.uni-hamburg.de>
#
# Permission to use, copy, modify, and distribute this software for any
# purpose with or without fee is hereby granted, provided that the above
# copyright notice and this permission notice appear in all copies.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

GOPATH         ?= $(shell pwd):$(VGOPATH)
GO             ?= go
GOFMT          ?= gofmt
GOFMT_OPTS     ?= -tabs=false -tabwidth=4
HIGHLIGHTFLAGS ?= --syntax go --line-numbers --anchors
JQUERYVERSION  ?= 2.1.0
SERIAL         ?= /dev/ttyUSB0

export GOPATH

.PHONY: all
all: infinitycube

.PHONY: infinitycube
infinitycube: content
	$(GO) get github.com/lucasb-eyer/go-colorful
	$(GO) get github.com/gordonklaus/portaudio
	$(GO) get github.com/mjibson/go-dsp/spectral
	$(GO) get github.com/kellydunn/go-opc
	$(GO) get github.com/fatih/structs
	$(GO) install infinitycube

.PHONY: format
format:
	$(GOFMT) -w=true $(GOFMT_OPTS) src/*/*.go

.PHONY: check-format
check-format:
	$(GOFMT) -d=true $(GOFMT_OPTS) src/*/*.go

static:
	mkdir -p static

static/%.js: frontend/%.coffee static
	coffee --compile --output static "$<"

SCRIPTS = static/infinitycube.js static/sound.js
STATIC  = frontend/*.html frontend/*.css frontend/*.js frontend/*.ico

.PHONY: content
content: static $(SCRIPTS)
	rsync -u $(STATIC) static

.PHONY: simulation
simulation: infinitycube content
#	simulation/gl_server -p 28144 -l simulation/cubeModel.json &
#	sleep 1
	bin/infinitycube -fcserver localhost:28144

.PHONY: deploy
deploy: web content
	rsync -av bin/web static "$(TARGET)"

.PHONY: clean
clean:
	rm -rf -- pkg bin static

.PHONY: cube
cube: infinitycube
	sudo bin/infinitycube -serial cube

.PHONY: test
test: infinitycube content
	$(GO) test infinitycube
