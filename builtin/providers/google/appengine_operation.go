package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	appengine "google.golang.org/api/appengine/v1"
)

type AppEngineOperationWaiter struct {
	Service *appengine.AppsService
	Op      *appengine.Operation
	AppsID  string
}

func (w *AppEngineOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func (w *AppEngineOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := w.Service.Get(w.AppsID).Do()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Progress of operation %q: %q", w.Op.Name, resp.ServingStatus)

		return resp, resp.ServingStatus, err
	}
}

func appEngineOperationWait(config *Config, op *appengine.Operation, appsID, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	w := &AppEngineOperationWaiter{
		Service: config.clientAppEngine.Apps,
		Op:      op,
		AppsID:  appsID,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}
