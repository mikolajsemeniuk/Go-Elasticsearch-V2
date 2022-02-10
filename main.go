package main

import (
	"es/application"
	"es/data"
	"es/extensions"
)

func main() {
	data.GetInfo()
	extensions.Info("[main.go] application starts...")
	application.Listen()
}
