.PHONY: test live-test recapture fmt vet

# Hermetic suite: schema conformance + behavior tests against committed
# captures. No network, runs in CI.
test:
	go test -race ./...

# Live conformance: walks testdata/live_ids.json against the real FPL API
# and validates our models via the same conformance rules. Never blocks PRs.
live-test:
	FPL_LIVE_TEST=1 go test -race -count=1 -run TestLiveConformance ./endpoints/ -v

# Refresh the hermetic captures in testdata/ from the live API. The hermetic
# suite must stay green afterwards unless the schema actually changed.
recapture:
	FPL_LIVE_TEST=1 FPL_RECAPTURE=1 go test -count=1 -run TestLiveConformance ./endpoints/ -v

fmt:
	gofmt -l .

vet:
	go vet ./...
