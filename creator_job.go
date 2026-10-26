package main

import "fmt"

type CreatorJob struct {
	CreatorID    string
	AssetID      string
	SubscriberID string
	ContentID    string
}

type JobResult struct {
	Status  string
	Message string
}

func processCreatorJob(job CreatorJob, deliver func(CreatorJob) error, updateSubscriber func(CreatorJob) error, processContent func(CreatorJob) error, capture func(map[string]any) error) JobResult {
	steps := []struct {
		name string
		run  func() error
	}{
		{"digital-asset-delivery", func() error { return deliver(job) }},
		{"subscriber-update", func() error { return updateSubscriber(job) }},
		{"content-processing", func() error { return processContent(job) }},
	}
	for _, step := range steps {
		if err := step.run(); err != nil {
			captureErr := capture(map[string]any{
				"title":       fmt.Sprintf("creator job %s failed", step.name),
				"message":     err.Error(),
				"level":       "error",
				"fingerprint": []string{job.CreatorID, step.name},
				"exception":   fmt.Sprintf("%T: %v", err, err),
				"context": map[string]string{
					"creator_id":    job.CreatorID,
					"asset_id":      job.AssetID,
					"subscriber_id": job.SubscriberID,
					"content_id":    job.ContentID,
					"step":          step.name,
				},
			})
			if captureErr != nil {
				return JobResult{Status: "failed", Message: fmt.Sprintf("%s; capture: %v", err, captureErr)}
			}
			return JobResult{Status: "failed", Message: fmt.Sprintf("%s: %v", step.name, err)}
		}
	}
	return JobResult{Status: "delivered", Message: "asset delivered, subscriber updated, content processed"}
}
