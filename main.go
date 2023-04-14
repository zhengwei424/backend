package main

import (
	"backend/Router"
)

func main() {
	// 调用顺序 main.go -> router(Middlewares) -> Controllers -> Services(sessions) -> Models -> Databases
	Router.InitRouter()
}
