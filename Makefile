VERSION = 1.0.3

APP      := tcp-time
PACKAGES := $(shell go list -f {{.Dir}} ./...)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: clean release docker docker-latest

# go get -u github.com/github/hub
release: zip
	git push
	hub release delete $(VERSION) || true
	hub release create $(VERSION) -m "$(VERSION)" -a release/$(APP)_$(VERSION)_osx_x86_64.zip -a release/$(APP)_$(VERSION)_windows_x86_64.zip -a release/$(APP)_$(VERSION)_linux_x86_64.zip

docker: binaries/linux_x86_64/$(APP)
	docker build -t quay.io/sergey_grebenshchikov/$(APP):v$(VERSION) .
	docker push quay.io/sergey_grebenshchikov/$(APP):v$(VERSION)

docker-latest: docker
	docker tag quay.io/sergey_grebenshchikov/$(APP):v$(VERSION) quay.io/sergey_grebenshchikov/$(APP):latest
	docker push quay.io/sergey_grebenshchikov/$(APP):latest

zip: release/$(APP)_$(VERSION)_osx_x86_64.zip release/$(APP)_$(VERSION)_windows_x86_64.zip release/$(APP)_$(VERSION)_linux_x86_64.zip

binaries: binaries/osx_x86_64/$(APP) binaries/windows_x86_64/$(APP).exe binaries/linux_x86_64/$(APP)

clean: 
	rm -rf binaries/
	rm -rf release/

release/$(APP)_$(VERSION)_osx_x86_64.zip: binaries/osx_x86_64/$(APP)
	mkdir -p release
	cd ./binaries/osx_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_osx_x86_64.zip $(APP)
	
binaries/osx_x86_64/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/osx_x86_64/$(APP) ./cmd/$(APP)

release/$(APP)_$(VERSION)_windows_x86_64.zip: binaries/windows_x86_64/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_64.zip $(APP).exe
	
binaries/windows_x86_64/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/windows_x86_64/$(APP).exe ./cmd/$(APP)

release/$(APP)_$(VERSION)_linux_x86_64.zip: binaries/linux_x86_64/$(APP)
	mkdir -p release
	cd ./binaries/linux_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_linux_x86_64.zip $(APP)
	
binaries/linux_x86_64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_x86_64/$(APP) ./cmd/$(APP)