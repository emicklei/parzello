package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	selfdiagnose "github.com/emicklei/go-selfdiagnose"
	"github.com/emicklei/go-selfdiagnose/task"
)

func addSelfdiagnose() {
	// add http handlers for /internal/selfdiagnose.(html|json|xml)
	selfdiagnose.AddInternalHandlers()
	selfdiagnose.Register(task.ReportHttpRequest{})
	selfdiagnose.Register(task.ReportHostname{})
	selfdiagnose.Register(task.ReportCPU())

	if len(os.Getenv("GOOGLE_CLOUD_PROJECT")) > 0 {
		// https://cloud.google.com/appengine/docs/standard/nodejs/runtime
		m := map[string]interface{}{}
		for _, each := range []string{"GAE_APPLICATION", "GAE_DEPLOYMENT_ID", "GAE_ENV", "GAE_INSTANCE", "GAE_MEMORY_MB", "GAE_RUNTIME",
			"GAE_SERVICE", "GAE_VERSION", "GOOGLE_CLOUD_PROJECT", "PORT"} {
			m[each] = os.Getenv(each)
		}
		selfdiagnose.Register(task.ReportVariables{VariableMap: m, Description: "Google AppEngine Environment"})
	}
}

type CheckTopicTask struct {
	Topic  string
	Client *pubsub.Client
}

func (c *CheckTopicTask) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	_, err := c.Client.Topic(c.Topic).Exists(context.Background())
	if err != nil {
		result.Passed = false
		result.Reason = fmt.Sprintf("failed to check topic existence: %v", err)
		return
	}
	result.Reason = fmt.Sprintf("exists in project [%s]", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	result.Passed = true
}
func (c *CheckTopicTask) Comment() string {
	return "topic: " + c.Topic
}

type CheckSubscriptionTask struct {
	Subscription string
	Client       *pubsub.Client
}

func (c *CheckSubscriptionTask) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	_, err := c.Client.Subscription(c.Subscription).Exists(context.Background())
	if err != nil {
		result.Passed = false
		result.Reason = fmt.Sprintf("failed to check subscription existence: %v", err)
		return
	}
	result.Reason = fmt.Sprintf("exists in project [%s]", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	result.Passed = true
}
func (c *CheckSubscriptionTask) Comment() string {
	return "subscription: " + c.Subscription
}
