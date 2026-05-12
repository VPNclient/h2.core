.PHONY: build-unix clean

build-unix:
	@bash ./build/unix.sh

clean:
	@rm -rf ./dist
