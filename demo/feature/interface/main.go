package main

type GameModule interface {
	Play()
	Init()
}

type MyGame struct {
	name string
}

func (mygame *MyGame) Play() {
	println("play")
	println(mygame.name)
}

func (mygame *MyGame) Init() {
	println("init")
}

func runner(game GameModule) {
	game.Init()
	game.Play()
}

func main() {
	mygame := MyGame{name: "demo"}

	runner(&mygame)
}
