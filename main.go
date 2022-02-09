package main

import (
	"es/application"
	"es/data"
)

func main() {
	data.GetInfo()
	application.Listen()
}
