package helpful

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
)

func TraceMsg(ctx context.Context, msg string) {
	s := opentracing.SpanFromContext(ctx)
	if s != nil {
		s.LogKV("msg", msg)
	}
}

func TraceErr(ctx context.Context, err error) {
	s := opentracing.SpanFromContext(ctx)
	if s != nil {
		s.SetTag("error", "true")
		msg := fmt.Sprint(err)
		s.LogKV("error_message", msg)
	}
}
