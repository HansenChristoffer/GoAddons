# This file is part of GoAddons, which is licensed under the GNU General Public License v3.0.
# You should have received a copy of the GNU General Public License along with this program.
# If not, see <https://www.gnu.org/licenses/>.

# Makefile for building GoAddons

# The name of the binary to be created
BINARY_NAME=goaddons

# Output directory for binaries
BIN_DIR=./bin

# Full path for the binary
BINARY_PATH=$(BIN_DIR)/$(BINARY_NAME)

# Retrieves the current version from version.go
VERSION=$(shell grep 'Version   =' version/Version.go | awk '{ print $$3 }' | tr -d '"')

# Retrieves the latest commit hash of the current branch
COMMIT_HASH=$(shell git rev-parse HEAD)

# Retrieves the current date and time
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=gofmt

# Ensure the bin directory exists
.PHONY: prepare
prepare:
	mkdir -p $(BIN_DIR)

# Check go fmt
.PHONY: fmt-check
fmt-check:
	@echo "Running gofmt"
	@test -z $($(GOFMT) -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

# Vet
.PHONY: vet
vet:
	@echo "Running go vet"
	@$(GOVET) ./... || (echo "[ERROR] Fix the issues reported by 'go vet' above" && exit 1)

# Build the project
.PHONY: release
release: prepare fmt-check vet
	$(GOBUILD) -ldflags "-X goaddons/version.Version=$(VERSION) -X goaddons/version.Commit=$(COMMIT_HASH) -X goaddons/version.BuildDate=$(BUILD_DATE)" -o $(BINARY_PATH) -v

# Clean up binaries
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...
