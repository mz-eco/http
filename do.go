package http

import "time"

type DoStmt interface {
	GetRequest() (*Request, bool)
}

func Run(x *Group, requests ...DoStmt) {

	for _, stmt := range requests {
		r, ok := stmt.GetRequest()

		if !ok {
			continue
		}

		Do(x, r)
	}
}

func Do(x *Group, r *Request) *Translate {

	x.check()

	var (
		ctx = x.context()
	)

	tx := &Translate{
		Name:    "",
		Request: r,
		Created: time.Now(),
	}

	done := func() {
		x.dispatch(ctx, AcDone, tx)
	}

	tx.Error = x.dispatch(ctx, AcPacker, r)

	if tx.Error != nil {
		done()
	}

	tx.Error = x.dispatch(ctx, AcRequest, r)

	if tx.Error != nil {
		done()
	}

	response, err := r.Do(c)

	if err != nil {
		tx.Error = err
		x.dispatch(ctx, AcDone, tx)
	} else {

		tx.Error = x.dispatch(ctx, AcUnPacker, response)

		if tx.Error != nil {
			done()
		}

		tx.Error = x.dispatch(ctx, AcResponse, response)

		if err != nil {
			done()
		}
	}

	tx.Response = response
	x.dispatch(ctx, AcDone, tx)

	return tx

}
