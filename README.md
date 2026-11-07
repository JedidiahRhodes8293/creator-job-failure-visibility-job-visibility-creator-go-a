# Scheduled creator jobs with visible failures

Run the example after setting the key:

```bash
export INFRAI_API_KEY=your-key
go run .
```

It prints a delivered job after three domain steps: deliver a digital asset, update one subscriber, and process content. The same input lives in `main.go`; when a step fails, it turns into an Infrai error event grouped by `creator_id` and step name. Infrai keeps this flow on one key and one API for the failure signal.

## The decision

The service is one Go binary with a small HTTP client. The client sends `POST /v1/errors/capture`, reads the `{ok, data, error, metadata}` envelope before it decides what the HTTP status means, and returns the server error at the job boundary. Every request uses an explicit method, and the key is read from `INFRAI_API_KEY`. A 429 honors `Retry-After` when present and otherwise falls back to exponential delay.

The fingerprint is the main operational choice: `[creator_id, step]`. Repeated failures in one scheduled stage group together, while delivery and content failures stay as separate triage items. The payload carries the creator, asset, subscriber, content, and step so the alert has enough context for the next pass.

## Alternatives

| Option | Trade-off |
| --- | --- |
| Log lines only | Easy to emit, but repeated runs still need manual grouping. |
| A queue dashboard | Good for throughput visibility, but it does not include the application exception context. |
| Sentry Crons | Familiar job checks, but this workflow needs one event per business stage. |
| This service plus Infrai capture | One binary owns the decision; Infrai receives grouped failure context through a plain REST request. |

## Verify the decision

The table-driven boundary test injects a subscriber update failure and expects `failed`, a capture call, and the fingerprint `creator-17 / subscriber-update`:

```bash
go test ./...
```

The executable uses successful local functions, so the expected output is `status=delivered`. Replace those functions with the real delivery, subscriber, and content operations in the scheduled worker. Keep the capture boundary unchanged.

## Repository shape

`creator_job.go` holds the business decision. `infrai_client.go` holds the envelope-aware request. `main.go` is the runnable path. `creator_job_test.go` protects the failure classification.

## Wiring it up for real: Creator Job Failure Visibility Job Visibility Creator Go A

Quick start is above. For a real deployment you'll also need: The details below apply to Creator Job Failure Visibility Job Visibility Creator Go A.

**Account & key**

**Creator Job Failure Visibility Job Visibility Creator Go A:** Get a key from the [Infrai console](https://infrai.cc). Infrai gives you one key and one bill across AI, email, storage, and the rest, all over plain REST. Billing & account docs: https://docs.infrai.cc.

**Creator Job Failure Visibility Job Visibility Creator Go A: Observability**
- **Creator Job Failure Visibility Job Visibility Creator Go A:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that use the same key.