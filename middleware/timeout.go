package middleware

import (
	"context"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	Next http.Handler
}

func (tm *TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//get the context and create a new one with a timeout of 10 seconds
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//ensure resources related with the context are released from memory once
	//the middleware operations are complete
	defer cancel()

	//give the request our newly created context
	r.WithContext(ctx)

	//create a channel to receive on when the request finishes
	chRequestFinished := make(chan struct{})

	//pass on to next chain on middleware (or the DefaultServeMux)
	//do this in an inline goroutine so we can check for completion or timeout on both
	//ctx.Done channel or our requestFinsihed channel
	go func() {
		tm.Next.ServeHTTP(w, r)
		chRequestFinished <- struct{}{}
	}()

	//now use a select construct to execute on whichever channel recevies signal 1st
	select {
	case <-chRequestFinished:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

//constructor for our timeout middleware
func NewTimeoutMiddleware(handler http.Handler) *TimeoutMiddleware {
	return &TimeoutMiddleware{Next: handler}
}
