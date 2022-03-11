package main

import (
	"context"

	"github.com/PiotrKozimor/krkstops/test/mock"
)

func main() {
	println("listening")
	ctx := context.Background()
	go mock.Ttss(ctx)
	mock.Airly(ctx)
}
