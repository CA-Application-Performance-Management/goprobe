/*Package goprobe ...*/
/*
This acts as a wrapper for every http request, users need to call this
wrapper for every incoming request to the server
*/
package goprobe

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

//seqGen type is used to store the seqNo and corId for every tid
type seqGen struct {
	seqNo int
	corid string
}

//HttpWrapper acts like a wrapper around the http request and will generate the following and route the handler
func HttpWrapper(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//uuid
		uuidvalue := uuid.Must(uuid.NewV4())
		tid := strconv.FormatInt(rand.Int63(), 10)
		uuidStr := uuidvalue.String()

		ctx := r.Context()
		ctx = context.WithValue(ctx, "tid", tid)
		ctx = context.WithValue(ctx, "corid", uuidStr)
		ctx = context.WithValue(ctx, "seqid", 0)
		ctx = context.WithValue(ctx, "httpMethod", r.Method)
		m.Lock()
		apmTidData[tid] = &seqGen{0, uuidStr}
		m.Unlock()
		f(w, r.WithContext(ctx))
	}
}
