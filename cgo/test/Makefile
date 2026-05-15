# Makefile for h2.core C API test

ROOT_DIR := $(shell cd ../.. && pwd)
DIST_DIR := $(ROOT_DIR)/dist

# Detect OS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    LIB_EXT := dylib
    ARCH := $(shell uname -m)
    ifeq ($(ARCH),arm64)
        LIB_NAME := libh2core_darwin_arm64.dylib
    else
        LIB_NAME := libh2core_darwin_amd64.dylib
    endif
    LDFLAGS := -Wl,-rpath,$(DIST_DIR)
else
    LIB_EXT := so
    ARCH := $(shell uname -m)
    ifeq ($(ARCH),aarch64)
        LIB_NAME := libh2core_linux_arm64.so
    else
        LIB_NAME := libh2core_linux_amd64.so
    endif
    LDFLAGS := -Wl,-rpath,$(DIST_DIR)
endif

CC := gcc
CFLAGS := -Wall -I$(DIST_DIR)
TARGET := test_h2core

.PHONY: all clean run

all: $(TARGET)

$(TARGET): main.c $(DIST_DIR)/$(LIB_NAME)
	$(CC) $(CFLAGS) -o $@ main.c -L$(DIST_DIR) -lh2core_$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell go env GOARCH) $(LDFLAGS)

run: $(TARGET)
	./$(TARGET)

clean:
	rm -f $(TARGET)

# Ensure library is built
$(DIST_DIR)/$(LIB_NAME):
	cd $(ROOT_DIR) && ./build/cgo.sh
