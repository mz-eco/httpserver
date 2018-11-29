package httpserver

import (
	"fmt"
	"net/http"
	R "reflect"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/mz-eco/types"
)

const (
	metaAsk   = 0
	metaQuery = 1
	metaParam = 2
	metaAll   = 3
)

var (
	metaFields = []string{
		"Ask",
		"Query",
		"Param",
	}

	metaBinding = []Binding{
		&jsonBinding{},
		&queryBinding{},
		&paramsBinding{},
	}
)

type Meta struct {
	Func interface{}

	types []R.Type
	fn    R.Value
	ctx   R.Type
}

func (m *Meta) check(x string) {

	var (
		errFunc = func() {
			panic(
				fmt.Sprintf(
					"Meta for [%s] func must func (ctx <*struct implement httpserver.Context>)", x))
		}
		fn = types.Value(m.Func)
	)

	if !types.IsIn(fn, types.StructPtr) {
		errFunc()
	}

	if !types.IsOut(fn) {
		errFunc()
	}

	var (
		in = types.TypeIn(fn, 0)
	)

	if !types.HasField(in, TypeContext) {
		errFunc()
	}

	m.types = make([]R.Type, metaAll)

	for index := 0; index < metaAll; index++ {
		m.types[index] = types.FieldTypeByName(in, metaFields[index])
	}

	m.fn = fn
	m.ctx = in.Elem()
}

func (m *Meta) handler(context *gin.Context) {

	var (
		in  = R.New(m.ctx)
		ctx = in.Interface().(contextGetter).context()
	)

	ctx.ctx = context

	for index := 0; index < metaAll; index++ {
		if m.types[index] == nil {
			continue
		}

		x := R.Indirect(in).FieldByName(metaFields[index])

		err := metaBinding[index].Bind(R.Indirect(x).Addr().Interface(), context)

		if err != nil {
			ctx.Error(ErrBinding, errors.Wrapf(err, "binding %s fail.", metaFields[index]))
			return
		}
	}

	m.fn.Call([]R.Value{in})

	if ctx.ack == nil {
		ctx.Error(ErrAck, nil)
	}

	switch x := ctx.ack.(type) {
	case string:
		context.String(
			http.StatusOK,
			x,
		)
	case []byte:
		context.Data(
			http.StatusOK,
			"application/octet-stream",
			x,
		)
	default:

		if types.Kind(x) == R.Slice {
			context.JSON(
				http.StatusOK,
				map[string]interface{}{
					"Code":    OK,
					"Message": errorMessage(OK),
					"Data": map[string]interface{}{
						"Items": x,
					},
				},
			)
		} else {
			context.JSON(
				http.StatusOK,
				map[string]interface{}{
					"Code":    OK,
					"Message": errorMessage(OK),
					"Data":    x,
				},
			)
		}

	}

}
