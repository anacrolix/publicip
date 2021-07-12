package main

import (
	"context"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/anacrolix/publicip"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var args struct {
		Ip4 bool `arg:"-4"`
		Ip6 bool `arg:"-6"`
	}
	arg.MustParse(&args)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if args.Ip4 {
		spew.Dump(publicip.Get4(ctx))
	}
	if args.Ip6 {
		spew.Dump(publicip.Get6(ctx))
	}
	if !args.Ip4 && !args.Ip6 {
		spew.Dump(publicip.GetAll(ctx))
	}
}
