package http

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	R "reflect"
	"strings"
	"time"

	"github.com/mz-eco/memoir"

	"github.com/mz-eco/types"
	"github.com/pkg/errors"
)

type argType int

const (
	argNil argType = iota
	argString
	argBytes
	argStruct
)

type invoker struct {
	fn R.Value
	x  R.Type
	at argType
}

func (m invoker) Invoke(ctx R.Value, v interface{}) error {

	var (
		in   = make([]R.Value, 3)
		size = 3
		read = func() []byte {
			return v.(BodyReader).Bytes()
		}
	)

	in[0] = ctx
	in[1] = types.Value(v)

	switch {

	}

	switch m.at {
	case argNil:
		size = 2
		break
	case argBytes:
		in[2] = types.Value(read())
	case argString:
		in[2] = types.Value(string(read()))
	case argStruct:
		var (
			x = R.New(m.x)
		)

		err := json.Unmarshal(read(), x.Elem().Interface())

		if err != nil {
			return err
		}

		in[2] = x.Elem()
	}

	o := m.fn.Call(in[:size])

	if o[0].IsNil() {
		return nil
	}

	return o[0].Interface().(error)

}

type Group struct {
	Domain  string
	Hooker  Hooker
	Writer  io.Writer
	Context interface{}

	ctx     R.Type
	fn      map[Action]map[string]invoker
	checked int
}

func (m *Group) check() {

	if m.checked == 0 {
		var (
			errCtx = func() {
				panic("Group.Context must be a <*struct>")
			}
		)

		if m.Context == nil {
			errCtx()
		}

		ctx := types.Type(m.Context)

		switch {
		case types.Is(ctx, types.StructPtr):
			m.ctx = ctx
		case types.Is(ctx, R.Ptr):
			if types.Is(ctx.Elem(), types.StructPtr) {
				m.ctx = ctx.Elem()
			} else {
				errCtx()
			}
		default:
			errCtx()
		}
	}

	if m.Writer == nil {
		m.Writer = os.Stdout
	}

	m.Hooker.Hook(m)

	m.checked = 1

}

var (
	typeRequest  = R.TypeOf((*Request)(nil))
	typeResponse = R.TypeOf((*Response)(nil))
	typeError    = R.TypeOf((*error)(nil)).Elem()
	typeBytes    = R.TypeOf(([]byte)(nil))
)

func (m *Group) context() R.Value {
	return types.New(m.ctx).Elem()
}

func (m *Group) find(action Action, name string) (invoker, bool) {

	if m.fn == nil {
		return invoker{}, false
	}

	ac, ok := m.fn[action]

	if !ok {
		return invoker{}, false
	}

	iv, ok := ac[name]

	if !ok {
		iv, ok = ac["*"]

		if !ok {
			return invoker{}, false
		}

		return iv, ok
	}

	return iv, true
}

func (m *Group) dispatch(ctx R.Value, action Action, v interface{}) error {

	switch action {
	case AcDone:

		var (
			x = v.(*Translate)
		)

		x.Used = time.Now().Sub(x.Created)

		memoir.Format(m.Writer, x)
	case AcRequest, AcPacker:

		var (
			x = v.(*Request)
		)

		if !strings.HasPrefix(x.URL, "http") {
			x.URL = fmt.Sprintf("%s%s", m.Domain, x.URL)
		}

		iv, ok := m.find(action, x.URL)

		if ok {
			return iv.Invoke(ctx, x)
		}

	case AcResponse, AcUnPacker:
		var (
			x = v.(*Response)
		)

		iv, ok := m.find(action, x.Request.URL)

		if ok {
			return iv.Invoke(ctx, x)
		}
	}

	return nil
}

func (m *Group) set(ac Action, name string, iv invoker) {

	if m.fn == nil {
		m.fn = make(map[Action]map[string]invoker)
	}

	ivs, ok := m.fn[ac]

	if !ok {
		ivs = make(map[string]invoker)
		m.fn[ac] = ivs
	}

	ivs[name] = iv

}

func (m *Group) invoker(ac Action, x R.Type, fn interface{}) invoker {

	var (
		at    = argNil
		throw = func() {
			panic(
				fmt.Sprintf("%s fn must be func(%s,%s, [*struct|string|bytes]) error",
					ac, m.ctx, x),
			)
		}
	)

	if !types.IsOut(fn, typeError) {
		throw()
	}

	switch {
	case types.IsIn(fn, m.ctx, x):
		at = argNil
	case types.IsIn(fn, m.ctx, x, types.StructPtr):
		at = argStruct
	case types.IsIn(fn, m.ctx, x, R.String):
		at = argString
	case types.IsIn(fn, m.ctx, x, typeBytes):
		at = argBytes
	default:
		throw()
	}

	return invoker{
		fn: types.Value(fn),
		x:  types.TypeIn(fn, 1),
		at: at,
	}
}

func (m *Group) On(name string, ac Action, fn interface{}) {

	if !types.IsFunc(fn) {
		panic(
			errors.Errorf(
				"argument error, #1 must want a func, but give a %s", types.Type(fn)).Error(),
		)
	}

	switch ac {
	case AcRequest, AcPacker:
		m.set(
			ac, name, m.invoker(ac, typeRequest, fn))
	case AcResponse, AcUnPacker:
		m.set(
			ac, name, m.invoker(ac, typeResponse, fn),
		)
	}
}

func (m *Group) Do(names ...string) {

	m.check()

	for _, name := range names {
		x := types.CallByName(m.Hooker, name)
		Run(m, x[0].(DoStmt))
	}

}
