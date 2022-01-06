# shipwright

## About

I've had the idea in my head for several years now that continuous integration and delivery are built on incredibly flaky foundations.

Since the inception of Travis CI, configuration-based pipeline tooling has been prevalent all over open source software, and with good reason. It works in all languages, it's mostly platform agnostic, and it's flexible. However, in my experience, it does not scale with several people, and its problems are seen quickly.

- Development and testing strictly involves trial and error
- Automated testing is non-existent
- Development tools are limited to defining configuration schemas, which provide a shallow understanding of what each key and value will do.
- Configuration languages have nuanced syntax which is typically a lot different than standard programming languages.
- Dependencies / dependency management is often not supported as it's just not a typical configuration language construct.
- Lack of debugging means large / extensive pipelines are flaky and make it difficult to diagnose issues.
  - Problems like these are demoralizing and often lead to neglect; no one wants to address issues like these because they're so difficult to debug.

A **shipwright** is a person who builds boats. Since everything in the Kubernetes world, from Helm to Harbor follows a boat theme, I thought it appropriate to thematically follow.

The idea behind this application is that it is more of a library than an application. Users should, instead of defining this amalgamation of `yaml/json/toml/whatever` and `bash`, define their build, package, and release processes programmatically and deterministically. This opens up a whole world of possibilities, like:

- Writing unit and integration tests for your build pipeline.
- Reusing and sharing build, package, and deployment definitions.
- Leveraging existing tooling to make it easier to develop and debug pipelines.
- Improved visualization by allowing pipelines to define metrics and traces.

## Glossary

- **Pipeline**: A pipeline is a generic sequence of steps. A pipeline can be a set of steps to build an application, or it can define how to take an artifact, package it, and push it to a package repository.
  - **Action**: A pipeline action is a single reusable component in a pipeline. Actions can have arguments and define outputs.
  - **Source**: A pipeline source defines what causes a pipeline to begin. For typical continuous integration builds, this source would be a commit or a tag. For a release pipeline, this could be a NATS message, a Google Cloud Pub/Sub message that says an artifact is available, or it could be another pipeline.
- **Artifact**: The tangible, end-result of a pipeline. Not all pipelines produce artifacts.
  - An example of an artifact would be a compiled binary or a docker image.
- **Target**: A target is a software release destination. It is the final place that an artifact is sent before it is used to serve user requests.

## How does it work?

The main idea behind `shipwright` is that it defers what is typically considered server logic into the client / pipeline definitions and library.

### Clients

Shipwright clients are used in the pipelines themselves. All pipelines are programs, and thus have a `func main()`.

There are currently three planned Clients, which are selected using the `-mode` CLI argument.

1. `run` - Runs the pipeline locally, attempting to emulate what would normally be executed in a standard CI pipeline.
2. `drone` - Generates a Drone configuration that matches the specified pipeline.
3. `docker` - Similar to `run`, but instead runs the pipeline using the Docker CLI

The importance of the `shipwright` package can not be understated.

### Writing a Pipeline

Generally a pipeline is whatever you want it to be.

There are some helpful tools in the to improve your visualization / pipeline tracing, accept arguments, define outputs.

In this example our pipeline has a few steps:

1. Clone out project.
   - Because not all pipelines are general CI, we need to explicitely clone our project. Other systems typically include this step automatically, and make it tedious or impossible to accomplish the reverse.
2. Run two tasks simultaneously. Either of these steps could fail, however we want to be able to restart at one specific step, or we want to know whic failed if only one does.
   - Install the nodejs
   - Install go dependencies
3. Cache the results of these scripts, and only re-download them if the yarn.lock/go.sum are updated.
4. Run `make build`, followed by `make package`, and store the resulting `example.tar.gz`

```go=
{{ .demos.Basic }}
```

Once committed, this script can be treated like any other pipeline script and can be automatically ran when a new commit is made.

More interestingly though, you can run this pipeline by running:

```bash
shipwright -path=pipeline.go
```

If your pipeline defines any inputs it will prompt you for them, or optionally you can provide them as arguments by using the `-arg-{argument}` flag.

## Caveats

### Supported languages

- [ ] Go

## Package Design Principles

There are a few main principles behind designing the client library / packages. These princples should encourage writing libraries that make it easy to write pipelines that are not excessively verbose.

- todo

## Examples

To view examples of pipelines, visit the [demo](./demo) folder. These demos are used in our automated tests.

---

## Questions

- Transitionary phases
  - If I'm using [Starlark | Drone yaml] how can I transition step-by-step from that to Shipwright without fully committing?
