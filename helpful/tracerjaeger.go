package helpful

import (
	"fmt"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func NewJaegerFromConfig(cfg Config) (opentracing.Tracer, io.Closer, error) {
	sn, err := cfg.GetString("SERVICE_NAME")
	if err != nil {
		return nil, nil, err
	}
	jep, err := cfg.GetString("JAEGER_ENDPOINT")
	if err != nil {
		return nil, nil, err
	}
	jd, err := cfg.GetString("JAEGER_DISABLED")
	if err != nil {
		return nil, nil, err
	}
	jst, err := cfg.GetString("JAEGER_SAMPLER_TYPE")
	if err != nil {
		return nil, nil, err
	}
	jsp, err := cfg.GetString("JAEGER_SAMPLER_PARAM")
	if err != nil {
		return nil, nil, err
	}
	jrls, err := cfg.GetString("JAEGER_REPORTER_LOG_SPANS")
	if err != nil {
		return nil, nil, err
	}
	jrmqs, err := cfg.GetString("JAEGER_REPORTER_MAX_QUEUE_SIZE")
	if err != nil{
		return nil, nil, err
	}

	err = os.Setenv("JAEGER_ENDPOINT", jep)
	if err != nil {
		return nil, nil, err
	}
	err = os.Setenv("JAEGER_DISABLED", jd)
	if err != nil {
		return nil, nil, err
	}
	err = os.Setenv("JAEGER_SAMPLER_TYPE", jst)
	if err != nil {
		return nil, nil, err
	}
	err = os.Setenv("JAEGER_SAMPLER_PARAM", jsp)
	if err != nil {
		return nil, nil, err
	}
	err = os.Setenv("JAEGER_REPORTER_LOG_SPANS", jrls)
	if err != nil {
		return nil, nil, err
	}
	err = os.Setenv("JAEGER_REPORTER_MAX_QUEUE_SIZE", jrmqs)
	if err != nil{
		return nil, nil, err
	}
	return NewJaegerFromEnv(sn)
}

func NewJaegerFromEnv(serviceName string) (opentracing.Tracer, io.Closer, error) {
	c, err := config.FromEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("init jaeger from env error: %v", err)
	}

	c.ServiceName = serviceName

	tracer, closer, err := c.NewTracer()
	if err != nil {
		return nil, nil, fmt.Errorf("create jaeger tracing error: %v", err)
	}
	return tracer, closer, nil
}
