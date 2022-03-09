package wrappers

import (
	"context"

	"github.com/grafana/shipwright/plumbing/pipeline"
	"github.com/grafana/shipwright/plumbing/plog"
	"github.com/sirupsen/logrus"
)

type LogWrapper struct {
	Opts pipeline.CommonOpts
	Log  *logrus.Logger
}

func (l *LogWrapper) Fields(ctx context.Context, step pipeline.Step) logrus.Fields {
	fields := plog.DefaultFields(ctx, step, l.Opts)

	return fields
}

func (l *LogWrapper) WrapStep(steps ...pipeline.Step) []pipeline.Step {
	for i := range steps {
		step := steps[i]
		action := steps[i].Action
		steps[i].Action = func(ctx context.Context, opts pipeline.ActionOpts) error {
			l.Log.WithFields(l.Fields(ctx, step)).Infoln("starting step'")

			stdoutFields := l.Fields(ctx, step)
			stdoutFields["stream"] = "stdout"

			stderrFields := l.Fields(ctx, step)
			stderrFields["stream"] = "stderr"

			opts.Stdout = l.Log.WithFields(stdoutFields).Writer()
			opts.Stderr = l.Log.WithFields(stderrFields).Writer()

			if err := action(ctx, opts); err != nil {
				l.Log.WithFields(l.Fields(ctx, step)).Infoln("encountered error", err.Error())
				return err
			}

			l.Log.WithFields(l.Fields(ctx, step)).Infoln("done running step without error")
			return nil
		}
	}

	return steps
}

func (l *LogWrapper) Wrap(wf pipeline.WalkFunc) pipeline.WalkFunc {
	return func(ctx context.Context, step ...pipeline.Step) error {
		steps := l.WrapStep(step...)

		if err := wf(ctx, steps...); err != nil {
			return err
		}
		return nil
	}
}
