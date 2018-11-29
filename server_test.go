package httpserver

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	type X struct {
		Context
		Query struct {
			Id int
		}
		Param struct {
			ID int
		}
	}

	GET("/:id", func(ctx *X) {
		ctx.Done(ctx.Param)
	})

	fmt.Println(ListenAndServe(":2244"))
}
