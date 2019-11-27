package service

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/flaw"
	"github.com/phogolabs/stride/contract"
)

// Validator validates a spec file
type Validator struct {
	Path     string
	Reporter contract.Reporter
}

// Validate validates the file
func (v *Validator) Validate() error {
	collector := flaw.ErrorCollector{}

	reporter := v.Reporter.With(contract.SeverityVeryHigh)
	reporter.Notice(" Validating spec...")

	var (
		ctx    = context.TODO()
		loader = openapi3.NewSwaggerLoader()
	)

	spec, err := loader.LoadSwaggerFromFile(v.Path)
	if err != nil {
		reporter.Error(" Validating spec fail: %v", err)
		// add the error to the collector
		collector.Wrap(err)
		return err
	}

	if components := &spec.Components; components != nil {
		reporter := v.Reporter.With(contract.SeverityHigh)
		reporter.Notice(" Validating components...")

		if err := spec.Components.Validate(ctx); err != nil {
			reporter.Error(" Validating components fail: %v", err)
			// add the error to the collector
			collector.Wrap(err)
		} else {
			reporter.Success(" Validating components successful")
		}
	}

	if security := spec.Security; security != nil {
		reporter := v.Reporter.With(contract.SeverityHigh)
		reporter.Info("廬Validating security...")

		if err := security.Validate(ctx); err != nil {
			reporter.Error("廬Validating security fail: %v", err)
			// add the error to the collector
			collector.Wrap(err)
		} else {
			reporter.Success("廬Validating security successful")
		}
	}

	if paths := spec.Paths; paths != nil {
		reporter := v.Reporter.With(contract.SeverityHigh)
		reporter.Info(" Validating paths...")

		if err := paths.Validate(ctx); err != nil {
			reporter.Error(" Validating paths fail: %v", err)
			// add the error to the collector
			collector.Wrap(err)
		} else {
			reporter.Success(" Validating paths successful")
		}
	}

	if servers := spec.Servers; servers != nil {
		reporter := v.Reporter.With(contract.SeverityHigh)
		reporter.Info("力Validating servers...")

		if err := servers.Validate(ctx); err != nil {
			reporter.Error("力Validating servers fail: %v", err)
			// add the error to the collector
			collector.Wrap(err)
		} else {
			reporter.Success("力Validating servers...")
		}
	}

	if len(collector) > 0 {
		reporter.Error(" Validating spec fail!")
		return flaw.Errorf("Please check the error log for more details")
	}

	reporter.Success(" Validating spec complete!")
	return nil
}
