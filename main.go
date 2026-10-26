package main

import "fmt"

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	job := CreatorJob{CreatorID: "creator-17", AssetID: "pack-2026-10", SubscriberID: "subscriber-08", ContentID: "post-441"}
	result := processCreatorJob(job,
		func(CreatorJob) error { return nil },
		func(CreatorJob) error { return nil },
		func(CreatorJob) error { return nil },
		client.Capture,
	)
	fmt.Printf("status=%s message=%s\n", result.Status, result.Message)
}
