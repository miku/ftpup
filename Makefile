SHELL := /bin/bash
TARGETS := ftpup
PKGNAME := ftpup
VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOLDFLAGS += -w -s
GOFLAGS = -ldflags "$(GOLDFLAGS)"

.PHONY: all
all: $(TARGETS)

%: %.go
	go build -o $@ -ldflags "$(GOLDFLAGS)" $<

.PHONY: clean
clean:
	rm -f $(TARGETS)
	rm -f $(PKGNAME)*deb
	rm -fr packaging/deb/$(PKGNAME)/usr
	rm -fr packaging/deb/$(PKGNAME)/etc

.PHONY: deb
deb: all
	mkdir -p packaging/deb/$(PKGNAME)/usr/local/bin
	mkdir -p packaging/deb/$(PKGNAME)/etc/systemd/system
	cp $(TARGETS) packaging/deb/$(PKGNAME)/usr/local/bin
	cp ftpup.service packaging/deb/$(PKGNAME)/etc/systemd/system
	cd packaging/deb && fakeroot dpkg-deb --build $(PKGNAME) .
	mv packaging/deb/$(PKGNAME)_*.deb .
