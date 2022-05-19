package drone

import (
	"fmt"

	"github.com/drone/drone-yaml/yaml"
	"github.com/grafana/shipwright/plumbing/pipeline"
)

// TODO: I'm lazy at the moment and haven't implemented reverse filters (exlude).
func addEvent(c yaml.Conditions, e pipeline.Event) (yaml.Conditions, error) {
	if branch, ok := e.Filters["branch"]; ok {
		c.Event.Include = append(c.Event.Include, "branch")
		c.Branch.Include = append(c.Branch.Include, branch.String())
	}

	if tag, ok := e.Filters["tag"]; ok {
		c.Event.Include = append(c.Event.Include, "tag")
		if tag != nil {
			c.Ref.Include = append(c.Ref.Include, fmt.Sprintf("refs/tags/%s", tag.String()))
		}
	}

	return c, nil
}

func (c *Client) Events(events []pipeline.Event) (yaml.Conditions, error) {
	conditions := yaml.Conditions{}
	for _, event := range events {
		c, err := addEvent(conditions, event)
		if err != nil {
			return yaml.Conditions{}, err
		}

		conditions = c
	}

	return conditions, nil
}
