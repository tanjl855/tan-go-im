package main

import app "gin_chat/router"

func main() {
	r := app.Router()
	r.Run(":8081")
}
