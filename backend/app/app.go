package main

import (
	"dichterdev/dichter.dev/api/internal/handlers/edf"
	"dichterdev/dichter.dev/api/internal/handlers/quotes"
)

func main() {
	edf.Start()
	quotes.Start()
}
