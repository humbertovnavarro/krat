prepare:
	go install mvdan.cc/garble@latest

build:
	garble build main.go -literals -tiny -seed=random
