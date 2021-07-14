package app

import "github.com/knadh/koanf"

// Config hold global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var Config = koanf.New(".")
