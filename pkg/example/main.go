package main

import (
	"context"
	"errors"

	golog "github.com/marcosstupnicki/go-log/pkg"
)

func main() {
	log, err := golog.New()
	if err != nil {
		panic(err)
	}

	log.Info(context.Background(), "foo",
		golog.WithField("mykey", mystruct{
			Foo: "foo",
			Bar: 1212,
		}),
		golog.WithField("err", errors.New("random error ocurred")))
}

type mystruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
