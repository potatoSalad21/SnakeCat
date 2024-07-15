build:
	@go build -o ./bin/snek github.com/potatoSalad21/SnakeCat/cmd

run: build
	@./bin/snek
