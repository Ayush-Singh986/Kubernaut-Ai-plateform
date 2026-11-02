package mocks

import (
	"context"

	"github.com/Ayush-Singh986/kubernaut/pkg/gateway/processing"
	"github.com/Ayush-Singh986/kubernaut/pkg/gateway/types"
)

// Compile-time interface compliance check.
var _ processing.RetryObserver = (*NoopRetryObserver)(nil)

// NoopRetryObserver is a shared no-op implementation of processing.RetryObserver.
// Use in any unit or integration test that needs a CRDCreator but doesn't care
// about retry observation (i.e., the test is not specifically validating retry audits).
type NoopRetryObserver struct{}

// OnRetryAttempt is a no-op — intentionally discards the notification.
func (*NoopRetryObserver) OnRetryAttempt(_ context.Context, _ *types.NormalizedSignal, _ int, _ error) {
}
