BUILDDIR=$(CURDIR)/build
GOBIN=$(CURDIR)/bin

GOMOBILE=$(GOBIN)/gomobile
# Add GOBIN to $PATH so `gomobile` can find `gobind`.
GOBIND=env PATH="$(GOBIN):$(PATH)" "$(GOMOBILE)" bind
IMPORT_HOST=github.com
IMPORT_PATH=$(IMPORT_HOST)/tzmax/apple-gotun2socks-library

.PHONY: apple apple_future clean clean-all

all: apple apple_future

apple: $(BUILDDIR)/apple/Tun2socks.xcframework

$(BUILDDIR)/apple/Tun2socks.xcframework: $(GOMOBILE)
  # MACOSX_DEPLOYMENT_TARGET and -iosversion should match what outline-client supports.
  # TODO(fortuna): -s strips symbols and is obsolete. Why are we using it?
	export MACOSX_DEPLOYMENT_TARGET=10.14; $(GOBIND) -iosversion=11.0 -target=ios,iossimulator,macos -o $@ -ldflags '-s -w' -bundleid com.tzmax.tun2socks ./tun2socks/

apple_future: $(BUILDDIR)/apple_future/Tun2socks.xcframework

$(BUILDDIR)/apple_future/Tun2socks.xcframework: $(GOMOBILE)
	$(GOBIND) -iosversion=13.1 -target=ios,iossimulator,maccatalyst -o $@ -ldflags '-s -w' -bundleid com.tzmax.tun2socks ./tun2socks/

$(GOMOBILE): go.mod
	env GOBIN="$(GOBIN)" go install golang.org/x/mobile/cmd/gomobile
	env GOBIN="$(GOBIN)" $(GOMOBILE) init

$(XGO): go.mod
	env GOBIN="$(GOBIN)" go install github.com/crazy-max/xgo

go.mod: main.go
	go mod tidy
	touch go.mod

clean:
	rm -rf "$(BUILDDIR)"
	go clean

clean-all: clean
	rm -rf "$(GOBIN)"
