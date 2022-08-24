package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.TODO()

	Test(&ctx)

	log.Print(ctx.Value("KEY"))
}

func Test(ctx *context.Context) {
	*ctx = context.WithValue(*ctx, "KEY", "VALUE")
}
