
all: build

.PHONY: test doc clean release rrun
test doc clean:
	cargo $@

build:
	cargo build --release

run:
	cargo run --release

debugbuild:
	cargo build

debugrun:
	cargo run

clearscreen:
	clear

clearscreen_debugbuild: clearscreen debugbuild

livedebugbuild:
	printf "%s\n" "+f\.rs$$" "-f" | clsync -q --have-recursive-sync -W . -R /dev/stdin -Mdirect -Smake -r1 --ignore-failures -t1 -w1 -- clearscreen_debugbuild

