package wrappers

import (
	"context"

	"github.com/grafana/shipwright/plumbing/pipeline"
	"github.com/grafana/shipwright/plumbing/plog"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type TraceWrapper struct {
	Opts   pipeline.CommonOpts
	Tracer opentracing.Tracer
}

func (l *TraceWrapper) Fields(ctx context.Context, step pipeline.Step[pipeline.Action]) logrus.Fields {
	fields := plog.DefaultFields(ctx, step, l.Opts)

	return fields
}

func TagSpan(span opentracing.Span, opts pipeline.CommonOpts, step pipeline.Step[pipeline.Action]) {
	span.SetTag("job", "shipwright")
	span.SetTag("build_id", opts.Args.BuildID)
}

func (l *TraceWrapper) WrapStep(steps ...pipeline.Step[pipeline.Action]) []pipeline.Step[pipeline.Action] {
	for i := range steps {
		// Steps that provide a nil action should continue to provide a nil action.
		// There is nothing for us to trace in the execution of this action anyways, though there is an implication that
		// this step may execute something that is not defined in the pipeline.
		if steps[i].Content == nil {
			continue
		}

		step := steps[i]
		action := step.Content
		steps[i].Content = func(ctx context.Context, opts pipeline.ActionOpts) error {
			parent := opentracing.SpanFromContext(ctx)
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, l.Tracer, step.Name, opentracing.ChildOf(parent.Context()))
			TagSpan(span, l.Opts, step)
			defer span.Finish()

			if err := action(ctx, opts); err != nil {
				span.SetTag("error", err)
				return err
			}

			return nil
		}
	}

	return steps
}

func (l *TraceWrapper) Wrap(wf pipeline.StepWalkFunc) pipeline.StepWalkFunc {
	return func(ctx context.Context, step ...pipeline.Step[pipeline.Action]) error {
		steps := l.WrapStep(step...)

		if err := wf(ctx, steps...); err != nil {
			return err
		}
		return nil
	}
}
