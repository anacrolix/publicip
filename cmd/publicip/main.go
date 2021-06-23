package main

import (
	"context"

	"github.com/anacrolix/publicip"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(publicip.GetAll(context.TODO()))
}
