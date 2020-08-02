GUEST_BIN = uap
SERVER_BIN= uapd

PKG_URL := main
APP_VERSION = $(shell git describe --always --tags --abbrev=0 | tr -d "[v\r\n]")
DATE_VERSION := $(shell date +%Y%m%d-%H%M)
GIT_VERSION := $(shell git rev-parse --short HEAD)
GIT_DATE_VERSION := $(GIT_VERSION)-$(DATE_VERSION)
AUTO_VERSIONING := -X $(PKG_URL).Version=$(APP_VERSION) -X $(PKG_URL).BuildDate=$(DATE_VERSION) -X $(PKG_URL).CommitSHA=$(GIT_VERSION)

all: guest server

rsrc:
	command -v rsrc || go get github.com/akavel/rsrc

guest: rsrc
	rsrc -manifest ./cmd/guest/uap.exe.manifest -ico ./cmd/guest/app.ico -o ./cmd/guest/rsrc.syso
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w $(AUTO_VERSIONING) -H windowsgui" -o $(GUEST_BIN).exe ./cmd/guest/

server:
	GOOS=linux go build -ldflags "-s -w $(AUTO_VERSIONING)" -o $(SERVER_BIN) ./cmd/server/

install: server
	systemctl --user stop uapd.service
	sudo cp -v ./uapd /usr/local/bin/
	cp -v ./uapd.service ~/.config/systemd/user/uapd.service
	systemctl --user enable --now uapd.service
	systemctl --user status uapd

pack:
	tar cvjf $(SERVER_BIN)-$(APP_VERSION).tar.bz2 $(SERVER_BIN) $(SERVER_BIN).service Makefile
	zip -r $(GUEST_BIN)-$(APP_VERSION).zip $(GUEST_BIN).exe ./windows

lint:
	golangci-lint run -v

fmt:
	gofmt -w -s -d .
	go vet "./..."

clean:
	@rm -f ./cmd/guest/rsrc.syso
	@rm -f $(GUEST_BIN) $(GUEST_BIN).exe $(SERVER_BIN)
	@rm -f $(SERVER_BIN)*.tar.bz2
	@rm -f $(GUEST_BIN)*.zip