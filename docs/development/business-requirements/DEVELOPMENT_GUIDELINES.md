# Development Guidelines: Patterns and Standards

This document provides authoritative guidance on **development patterns** and **coding standards** for the kubernaut system.

## 🎯 **Decision Framework**

### Quick Decision Tree
```
📝 QUESTION: What pattern should I use?

├─ 🔄 "Updating Kubernetes resource status?"
│  └─ Use retry.RetryOnConflict ─────────────► STATUS UPDATE PATTERN
│
├─ 🔐 "Logging sensitive data?"
│  └─ Use pkg/shared/sanitization ───────────► SANITIZATION PATTERN
│
├─ 📝 "Writing audit events?"
│  └─ Use pkg/audit (BufferedAuditStore) ────► AUDIT PATTERN
│
└─ 🔧 "Creating shared utilities?"
   └─ Use pkg/shared/ ───────────────────────► SHARED PACKAGE PATTERN
```

---

## 🔄 **Kubernetes Status Update Pattern - MANDATORY**

### The Problem: Optimistic Locking Conflicts

Kubernetes uses **optimistic concurrency control** via `resourceVersion`. When multiple reconcilers update the same resource simultaneously, conflicts occur:

```
Error: the object has been modified; please apply your changes to the latest version and try again
```

### The Solution: `retry.RetryOnConflict`

**ALL Kubernetes status updates MUST use `retry.RetryOnConflict`** from `k8s.io/client-go/util/retry`.

#### ✅ CORRECT Pattern

```go
import "k8s.io/client-go/util/retry"

// Standard retry pattern for status updates
func (r *MyReconciler) updateStatus(ctx context.Context, obj *myv1.MyResource) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // 1. ALWAYS refetch to get latest resourceVersion
        if err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
            return err
        }

        // 2. Apply status changes
        obj.Status.Phase = "Ready"
        obj.Status.Message = "Processing complete"

        // 3. Update status subresource
        return r.Status().Update(ctx, obj)
    })
}
```

#### ❌ WRONG Patterns

```go
// ❌ WRONG: No retry - lost updates on conflict
if err := r.Status().Update(ctx, obj); err != nil {
    log.Error(err, "failed to update status")
    return err
}

// ❌ WRONG: Custom retry logic (reinventing the wheel)
for i := 0; i < 3; i++ {
    err := r.Status().Update(ctx, obj)
    if err == nil {
        return nil
    }
    // Custom backoff...
}

// ❌ WRONG: No refetch before update
return retry.RetryOnConflict(retry.DefaultRetry, func() error {
    // Missing: r.Get() to refresh resourceVersion
    obj.Status.Phase = "Ready"
    return r.Status().Update(ctx, obj)  // Will fail with stale resourceVersion
})
```

### Services Using This Pattern

| Service | File | Implementation |
|---------|------|----------------|
| **Gateway** | `pkg/gateway/processing/status_updater.go` | ✅ `retry.RetryOnConflict` |
| **RemediationOrchestrator** | `pkg/remediationorchestrator/controller/reconciler.go` | ✅ `retry.RetryOnConflict` |
| **RemediationOrchestrator** | `pkg/remediationorchestrator/handler/*.go` | ✅ `retry.RetryOnConflict` |
| **SignalProcessing** | Documented in implementation plans | ✅ Pattern defined |
| **Notification** | `internal/controller/notification/` | ⚠️ Migration in progress |

### Custom Backoff Configuration

For latency-sensitive services, customize the backoff:

```go
import "k8s.io/client-go/util/retry"

// Custom backoff for Gateway (P95 <50ms requirement)
var GatewayRetryBackoff = wait.Backoff{
    Steps:    5,
    Duration: 10 * time.Millisecond,  // Start small
    Factor:   2.0,
    Jitter:   0.1,
}

// Usage
return retry.RetryOnConflict(GatewayRetryBackoff, func() error {
    // ...
})
```

### Why This Matters

| Issue | Without Retry | With `retry.RetryOnConflict` |
|-------|--------------|------------------------------|
| Concurrent reconcilers | ❌ Lost updates | ✅ Automatic retry |
| Status consistency | ❌ Stale data | ✅ Fresh resourceVersion |
| Test reliability | ❌ Flaky tests | ✅ Deterministic |
| Production stability | ❌ Race conditions | ✅ Optimistic locking |

---

## 🔐 **Log Sanitization Pattern - MANDATORY**

### The Problem: Sensitive Data in Logs

Logging sensitive data (passwords, tokens, API keys) creates security vulnerabilities (CVSS 5.3).

### The Solution: Shared Sanitization Library

**ALL services MUST use `pkg/shared/sanitization/`** for log sanitization.

#### ✅ CORRECT Pattern

```go
import "github.com/Ayush-Singh986/kubernaut/pkg/shared/sanitization"

// Simple string sanitization
sanitized := sanitization.SanitizeForLog(sensitiveData)
logger.Info("Processing request", "payload", sanitized)

// With fallback for graceful degradation (BR-NOT-055)
sanitizer := sanitization.NewSanitizer()
clean, err := sanitizer.SanitizeWithFallback(content)
if err != nil {
    logger.Error(err, "Sanitization used fallback mode")
}
```

#### ❌ WRONG Patterns

```go
// ❌ WRONG: Logging sensitive data directly
logger.Info("Connection string", "url", dbURL)  // May contain password!

// ❌ WRONG: Custom sanitization logic
func sanitize(s string) string {
    return strings.ReplaceAll(s, "password", "[REDACTED]")  // Incomplete!
}

// ❌ WRONG: Service-specific sanitization package
import "github.com/Ayush-Singh986/kubernaut/pkg/myservice/sanitization"  // Use shared!
```

### Patterns Covered

The shared library covers all DD-005 required patterns:
- Passwords (`password`, `passwd`, `pwd`)
- API Keys (OpenAI `sk-*`, AWS keys)
- Tokens (Bearer, GitHub `ghp_*`)
- Database URLs (PostgreSQL, MySQL, MongoDB)
- Certificates (PEM format)
- Kubernetes Secrets (base64 data)

### Migration Status

| Service | Status |
|---------|--------|
| **Gateway** | 🟡 Pending migration |
| **Notification** | ✅ **MIGRATED** (V1.0) |
| **Data Storage** | ✅ Compliant (structured logging) |

---

## 📝 **Audit Trail Pattern - MANDATORY**

### The Problem: Compliance and Traceability

All services must emit audit events for compliance (DD-AUDIT-003, ADR-034).

### The Solution: Shared Audit Library

**ALL CRD controllers MUST use `pkg/audit/`** for audit event emission.

#### ✅ CORRECT Pattern

```go
import "github.com/Ayush-Singh986/kubernaut/pkg/audit"

// Initialize in main.go
auditClient := audit.NewHTTPDataStorageClient(dataStorageURL, &http.Client{})
auditStore := audit.NewBufferedAuditStore(auditClient, audit.BufferedAuditStoreConfig{
    BufferSize:    10000,
    BatchSize:     100,
    FlushInterval: 5 * time.Second,
})
defer auditStore.Close()

// Emit audit events (fire-and-forget, non-blocking)
event := &audit.AuditEvent{
    EventType:     "notification.message.sent",
    CorrelationID: notification.Spec.RemediationID,
    ActorID:       "notification-controller",
    // ... other fields
}
auditStore.Write(ctx, event)  // Non-blocking
```

### Key Principles

1. **Fire-and-forget**: Audit writes MUST NOT block business logic (BR-NOT-063)
2. **Correlation**: All events MUST include `correlation_id` for tracing (BR-NOT-064)
3. **Graceful degradation**: Audit failures don't fail the reconciliation
4. **Buffered writes**: Use `BufferedAuditStore` for efficiency (ADR-038)

---

## 📦 **Shared Package Pattern**

### When to Create Shared Code

| Scenario | Action |
|----------|--------|
| Used by 2+ services | ✅ Create in `pkg/shared/` |
| Used by 1 service | ❌ Keep in service package |
| Cross-cutting concern (logging, errors) | ✅ Create in `pkg/shared/` |
| Service-specific logic | ❌ Keep in service package |

### Package Structure

```
pkg/shared/
├── sanitization/     # Log sanitization (DD-005)
├── types/            # Shared type definitions
├── errors/           # Structured error types
└── k8sutil/          # Kubernetes utilities (future)
```

### ✅ CORRECT: Shared Package Creation

```go
// pkg/shared/myutil/myutil.go
package myutil

// Document the business requirement
// BR-XXX-XXX: Description of requirement

func MySharedFunction() {
    // Implementation used by multiple services
}
```

### ❌ WRONG: Duplicate Code

```go
// pkg/service1/util.go
func sanitize(s string) string { /* ... */ }

// pkg/service2/util.go
func sanitize(s string) string { /* ... */ }  // DUPLICATE!
```

---

## 🔧 **Error Handling Pattern - MANDATORY**

### The Rule: ALWAYS Handle Errors

**Errors MUST NEVER be ignored.** Every error must be either:
1. Returned to the caller
2. Logged with context
3. Both

#### ✅ CORRECT Pattern

```go
result, err := doSomething()
if err != nil {
    log.Error(err, "Failed to do something", "context", value)
    return fmt.Errorf("doSomething failed: %w", err)  // Wrap with context
}
```

#### ❌ WRONG Patterns

```go
// ❌ WRONG: Ignoring error
result, _ := doSomething()

// ❌ WRONG: Empty error handling
if err != nil {
    // Do nothing
}

// ❌ WRONG: Log without return (in functions that should propagate)
if err != nil {
    log.Error(err, "Failed")
    // Missing: return err
}
```

---

## 🏗️ **CRD Controller Pattern**

### Standard Controller Structure

```go
type MyReconciler struct {
    client.Client
    Scheme     *runtime.Scheme

    // Dependencies
    AuditStore audit.AuditStore
    Sanitizer  *sanitization.Sanitizer
}

func (r *MyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 1. Fetch resource
    obj := &myv1.MyResource{}
    if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil  // Resource deleted
        }
        return ctrl.Result{}, err
    }

    // 2. Business logic
    // ...

    // 3. Update status with retry
    return ctrl.Result{}, r.updateStatus(ctx, obj)
}

func (r *MyReconciler) updateStatus(ctx context.Context, obj *myv1.MyResource) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        if err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
            return err
        }
        // Apply status changes
        return r.Status().Update(ctx, obj)
    })
}
```

---

## 📋 **Quick Reference**

| Pattern | Import | Usage |
|---------|--------|-------|
| Status Update | `k8s.io/client-go/util/retry` | `retry.RetryOnConflict(...)` |
| Sanitization | `pkg/shared/sanitization` | `sanitization.SanitizeForLog(...)` |
| Audit | `pkg/audit` | `auditStore.Write(ctx, event)` |
| Errors | Standard Go | Always handle, never ignore |

---

## 🔗 **Related Documents**

| Document | Purpose |
|----------|---------|
| [TESTING_GUIDELINES.md](./TESTING_GUIDELINES.md) | Testing patterns and standards |
| [NAMING_CONVENTIONS.md](./NAMING_CONVENTIONS.md) | Naming patterns |
| [DD-005-OBSERVABILITY-STANDARDS.md](../../architecture/decisions/DD-005-OBSERVABILITY-STANDARDS.md) | Observability requirements |
| [DD-GATEWAY-011](../../architecture/decisions/DD-GATEWAY-011-shared-status-deduplication.md) | Status update pattern origin |

---

**Document Version**: 1.0
**Created**: December 9, 2025
**Last Updated**: December 9, 2025
**Authority**: Authoritative development patterns for kubernaut

