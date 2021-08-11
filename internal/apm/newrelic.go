package apm

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/proviant-io/core/internal/config"
	"log"
	"net/http"
)

type Apm interface {
	WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request))
}

func NewApm(cfg config.APM) Apm{
	if cfg.Vendor == "newrelic"{
		return newNewRelic(cfg.ApplicationName, cfg.LicenseKey)
	}
	return newNoop()
}

type NoopApm struct {
}

func (n *NoopApm) WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return pattern, handler
}

func newNoop() Apm {
	return &NoopApm{}
}

type NewRelicApm struct {
	app *newrelic.Application
}

func (n *NewRelicApm) WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return newrelic.WrapHandleFunc(n.app, pattern, handler)
}

func newNewRelic(name, key string) Apm {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense(key),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		log.Printf("error in apm setup: %s", err.Error())
		return nil
	}

	return &NewRelicApm{
		app: app,
	}
}
