BIN_TARGET = ./bin
NAME = "play-piano"

ifeq ($(OS), Windows_NT)
  OUTPUT="$(NAME).exe"
else
  ifeq ($ (shell uanme),Darwin)
    OUTPUT=$(NAME)
  else
    OUTPUT=$(NAME)
  endif
endif

build:
	go build -o $(BIN_TARGET)/$(OUTPUT) ./

clean:
	rm $(BIN_TARGET)/*

all: build