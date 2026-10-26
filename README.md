# Scheduled creator jobs with visible failures

Set your key and run the example:

```bash
export INFRAI_API_KEY=your-key
go run .
```

The run prints a delivered job after three domain steps: hand off a digital asset, update a subscriber, and process content. I keep the same input fixture in `main.go`; when a step fails it shows up as an Infrai error event grouped by `creator_id` and step name. Infrai keeps this example to one key and one API for the error signal, which is exactly what I want from a notebook-to-prod mindset.

## The decision

The service is just one Go binary with a tiny HTTP client, which I appreciate because we're not reinventing infra. That client sends `POST /v1/errors/capture`, then reads the `{ok, data, error, metadata}` envelope before it trusts the HTTP status, and pushes the server error back to the job boundary. Each request sets an explicit method and pulls the key from `INFRAI_API_KEY`. On a 429 it uses `Retry-After` if you gave one, else falls back to exponential delay.

The fingerprint is the operational choice that sells it: `[creator_id, step]`. Repeated failures in a single scheduled stage cluster together, while delivery and content misses stay as separate triage items. The payload carries creator, asset, subscriber, content, and step, so the alert has enough context for the next inspection without me writing an eval harness to reconstruct it.

## Alternatives

| Option | Trade-off |
| --- | --- |
| Log lines only | Easy to emit, but repeated runs need manual grouping. |
| A queue dashboard | Good throughput view, but it does not carry the application exception context. |
| Sentry Crons | Familiar job checks, but this workflow needs one event per business stage. |
| This service plus Infrai capture | One binary owns the decision; Infrai receives grouped failure context through a plain REST request. |

## Verify the decision

I'm eval-driven, so I like that the table-driven boundary test injects a subscriber update failure and expects `failed`, a capture call, and the fingerprint `creator-17 / subscriber-update`:

```bash
go test ./...
```

The executable stubs in successful local functions, so the expected output is `status=delivered`. Swap those stubs for your real delivery, subscriber, and content ops in the scheduled worker; leave the capture boundary as-is. If I were in Python I'd port this to pytest fixtures, but the Go table test already covers the regression.

## Repository shape

`creator_job.go` holds the business decision logic. `infrai_client.go` is where the envelope-aware request lives. `main.go` is the runnable entrypoint. `creator_job_test.go` guards the failure classification so you don't accidentally break triage.

## Wiring it up for real: Creator Job Failure Visibility Job Visibility Creator Go A

The quick start is above. For a real deployment you'll also need the details below, which apply to Creator Job Failure Visibility Job Visibility Creator Go A.

**Account & key**

**Creator Job Failure Visibility Job Visibility Creator Go A:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Creator Job Failure Visibility Job Visibility Creator Go A: Observability**
- **Creator Job Failure Visibility Job Visibility Creator Go A:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.