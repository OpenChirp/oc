# Makes the distributable binaries
# The idea was to make this runnable by anyone, even without the GOPATH setup
#
# Craig Hesling
# January 29, 2019

# The name of the output binary basename
BINARY=oc
# Directory to place builds in
BUILDS=builds

SOURCES=$(wildcard *.go)

# Sorta need this obscure mkdir line in all targets because GNU Make seems
# to check the directory access/modification time and constantly forces a rebuild
# of the enclosed build targets.
MKDIR_LINE=mkdir -p $(BUILDS)
BUILD_LINE=GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ $(SOURCES)
# We need to can go get for each arch/os, since some libraries depend on
# different libraries when using compiling for different arch/os
GOGET_LINE=GOOS=$(GOOS) GOARCH=$(GOARCH) go get -v .

.PHONY: all clean

# Build binary for all platforms
all: $(addprefix $(BUILDS)/, $(BINARY) $(BINARY).arm $(BINARY).osx $(BINARY).exe)

$(BUILDS)/$(BINARY): GOOS=linux
$(BUILDS)/$(BINARY): GOARCH=amd64
$(BUILDS)/$(BINARY): BINNAME=$(BINARY)
$(BUILDS)/$(BINARY): $(SOURCES) Makefile
	$(GOGET_LINE)
	$(MKDIR_LINE)
	$(BUILD_LINE)

$(BUILDS)/$(BINARY).arm: GOOS=linux
$(BUILDS)/$(BINARY).arm: GOARCH=arm64
$(BUILDS)/$(BINARY).arm: BINNAME=$(BINARY)
$(BUILDS)/$(BINARY).arm: $(SOURCES) Makefile
	$(GOGET_LINE)
	$(MKDIR_LINE)
	$(BUILD_LINE)

$(BUILDS)/$(BINARY).osx: GOOS=darwin
$(BUILDS)/$(BINARY).osx: GOARCH=amd64
$(BUILDS)/$(BINARY).osx: BINNAME=$(BINARY).osx
$(BUILDS)/$(BINARY).osx: $(SOURCES) Makefile
	$(GOGET_LINE)
	$(MKDIR_LINE)
	$(BUILD_LINE)

$(BUILDS)/$(BINARY).exe: GOOS=windows
$(BUILDS)/$(BINARY).exe: GOARCH=amd64
$(BUILDS)/$(BINARY).exe: BINNAME=$(BINARY).exe
$(BUILDS)/$(BINARY).exe: $(SOURCES) Makefile
	$(GOGET_LINE)
	$(MKDIR_LINE)
	$(BUILD_LINE)

clean:
	$(RM) -r $(BUILDS)
