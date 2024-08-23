package main

import (
	sweetyml "sweet-common/yaml"
	appMain "sweet-src/main/golang"
)

func main() {
	sweetyml.Init()
	appMain.Main()
}
