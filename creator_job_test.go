package main

import (
	"errors"
	"testing"
)

func TestProcessCreatorJobCapturesTheFailedStep(t *testing.T) {
	cases := []struct {
		name string
		fail func(CreatorJob) error
		step string
	}{
		{"delivery failure", func(CreatorJob) error { return errors.New("asset delivery rejected") }, "digital-asset-delivery"},
		{"subscriber failure", func(CreatorJob) error { return errors.New("subscriber store rejected update") }, "subscriber-update"},
		{"content failure", func(CreatorJob) error { return errors.New("content processor rejected input") }, "content-processing"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured map[string]any
			job := CreatorJob{CreatorID: "creator-17", AssetID: "pack-2026-10", SubscriberID: "subscriber-08", ContentID: "post-441"}
			deliver := func(CreatorJob) error { return nil }
			update := func(CreatorJob) error { return nil }
			process := func(CreatorJob) error { return nil }
			switch tc.step {
			case "digital-asset-delivery":
				deliver = tc.fail
			case "subscriber-update":
				update = tc.fail
			case "content-processing":
				process = tc.fail
			}
			result := processCreatorJob(job, deliver, update, process, func(payload map[string]any) error { captured = payload; return nil })
			if result.Status != "failed" {
				t.Fatalf("status = %q, want failed", result.Status)
			}
			if captured["title"] != "creator job "+tc.step+" failed" {
				t.Fatalf("title = %v", captured["title"])
			}
			fingerprint := captured["fingerprint"].([]string)
			if fingerprint[0] != job.CreatorID || fingerprint[1] != tc.step {
				t.Fatalf("fingerprint = %v", fingerprint)
			}
		})
	}
}
