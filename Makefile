BUILD_DIR := $(CURDIR)/build

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building project in $(BUILD_DIR)..."
	cmake -S '$(CURDIR)' -B '$(BUILD_DIR)' -G Ninja
	ninja -C '$(BUILD_DIR)'
	@echo

.PHONY: run
run: build
	@echo "Running project..."
	'$(BUILD_DIR)/minbox'
	@echo

.PHONY: format
format:
	clang-format --verbose -i -style=file src/*.cpp src/*.h