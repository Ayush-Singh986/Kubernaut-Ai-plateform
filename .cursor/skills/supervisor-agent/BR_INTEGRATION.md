# Business Requirements Integration

## 🔗 How Supervisor Agent Uses Business Requirements

The supervisor agent reads your BR documents and orchestrates their implementation.

---

## 📁 BR Document Structure (Your Format)

Based on `/Users/jgil/go/src/github.com/Ayush-Singh986/cursor-swarm-dev/docs/business-requirements/`:

```
BR-001-CURSOR-SWARM-VISION.md
BR-002-MULTI-TEAM-PARALLEL-WORK.md
BR-003-STICKY-MESSAGE-UX.md
BR-004-PERMISSION-MANAGEMENT.md
```

### BR Document Format

```markdown
# BR-XXX: [Feature Name]

**Status:** ✅ Approved | 🔄 In Progress | ❌ Rejected
**Date:** [Date]
**Author:** [Name]
**Priority:** P0 | P1 | P2

## 1. Business Problem
[Current pain points]

## 2. Business Objective
[What we want to achieve]

## 3. Acceptance Criteria
[Specific testable criteria]

## 4. Technical Requirements
[Implementation details]

## 5. Constraints
[Limitations and considerations]
```

---

## 🤖 Supervisor Agent Parsing

### Phase 1: Read & Understand

```
Supervisor reads BR document and extracts:

1. Business Objective → Overall goal for worker
2. Acceptance Criteria → Validation checkpoints
3. Technical Requirements → Implementation guidance
4. Constraints → Guardrails for worker
```

### Example: BR-WORKFLOW-197 Parsing

**Input Document**:
```markdown
# BR-WORKFLOW-197: Workflow Label Validation

## Business Objective
Prevent invalid workflows from being created by validating mandatory labels.

## Acceptance Criteria
- AC-1: System rejects workflows missing mandatory labels
- AC-2: System validates label format (key=value)
- AC-3: System provides clear error messages
- AC-4: Validation happens before persistence

## Technical Requirements
- Implement validator in pkg/datastorage/
- Integrate with workflow creation handler
- Unit test coverage 70%+
- Use existing validation patterns

## Constraints
- Must not break existing workflow creation
- Performance: validation <10ms per workflow
- Must work with all workflow types
```

**Supervisor Extracts**:
```
Goal: Prevent invalid workflows via label validation

Success Criteria:
✓ Rejects workflows without mandatory labels
✓ Validates label format
✓ Clear error messages
✓ Pre-persistence validation

Implementation Plan:
- Location: pkg/datastorage/
- Integration: workflow creation handler
- Testing: 70%+ coverage, Ginkgo/Gomega
- Patterns: Follow existing validators

Constraints:
- No breaking changes
- <10ms validation time
- All workflow types supported
```

---

## 📋 Task Decomposition from BR

### Supervisor's Task Breakdown Process

**From BR-WORKFLOW-197**, supervisor creates:

```
Phase 1 - RED (Tests):
  Task 1.1: Write test for AC-1 (reject missing labels)
    File: test/unit/datastorage/validator_test.go
    Focus: Behavior - workflow should be rejected
    
  Task 1.2: Write test for AC-2 (validate format)
    File: test/unit/datastorage/validator_test.go
    Focus: Behavior - invalid format should error
    
  Task 1.3: Write test for AC-3 (error messages)
    File: test/unit/datastorage/validator_test.go
    Focus: Behavior - messages should be clear and actionable
    
  Task 1.4: Write test for AC-4 (pre-persistence)
    File: test/integration/datastorage/workflow_creation_test.go
    Focus: Behavior - validation occurs before DB write

Phase 2 - GREEN (Implementation):
  Task 2.1: Implement validator logic
    File: pkg/datastorage/validator.go
    Requirement: Minimal logic to pass tests
    
  Task 2.2: Integrate in workflow handler
    File: pkg/datastorage/server/workflow_handlers.go
    Requirement: Call validator before persistence
    
  Task 2.3: Wire in main application
    File: cmd/datastorage/main.go
    Requirement: MANDATORY - no orphaned code
    
  Task 2.4: Verify tests pass
    Command: make test-unit-datastorage

Phase 3 - REFACTOR (Enhancement):
  Task 3.1: Extract validation rules to config
    File: pkg/datastorage/config/validation_rules.yaml
    Enhancement: Configurable rules
    
  Task 3.2: Optimize performance
    File: pkg/datastorage/validator.go
    Enhancement: Ensure <10ms constraint met
    
  Task 3.3: Add detailed error context
    File: pkg/datastorage/validator.go
    Enhancement: Rich error details for debugging

Checkpoints:
✓ After Task 1.4: Validate test quality
✓ After Task 2.4: Validate integration + passing tests
✓ After Task 3.3: Final comprehensive review
```

---

## ✅ Checkpoint Validation Mapped to AC

### Checkpoint 1: RED Phase

**Supervisor validates tests against Acceptance Criteria**:

```
AC-1: "System rejects workflows missing mandatory labels"
→ Test: "should reject workflow when mandatory label missing"
✓ Validates behavior (rejection) not implementation
✓ Clear expected outcome

AC-2: "System validates label format (key=value)"
→ Test: "should reject workflow when label format invalid"
✓ Validates behavior (format validation) not implementation
✓ Covers edge cases (no =, multiple =, empty value)

AC-3: "System provides clear error messages"
→ Test: "should return descriptive error for missing label"
✓ Validates behavior (error message quality) not implementation
✓ Checks message includes label name

AC-4: "Validation happens before persistence"
→ Integration Test: "should validate before database write"
✓ Validates behavior (order of operations) not implementation
✓ Verifies DB call not made when validation fails
```

**Decision**:
- ✅ All AC covered with behavior tests
- ✅ No implementation-logic tests found
- ✅ Proceed to GREEN

---

### Checkpoint 2: GREEN Phase

**Supervisor validates implementation against Technical Requirements**:

```
Requirement: "Implement validator in pkg/datastorage/"
→ Check: pkg/datastorage/validator.go exists
✓ File found in correct location

Requirement: "Integrate with workflow creation handler"
→ Check: grep "validator" pkg/datastorage/server/workflow_handlers.go
✓ Handler calls validator.Validate()

Requirement: "Integrate in main application"
→ Check: grep "validator" cmd/datastorage/main.go
✓ Validator initialized and passed to handlers

Requirement: "Unit test coverage 70%+"
→ Check: make test-coverage-datastorage
✓ Coverage: 85% (exceeds requirement)

Requirement: "Use existing validation patterns"
→ Check: Compare with pkg/datastorage/schema_validator.go
✓ Follows same error return pattern
```

**Verify Constraints**:
```
Constraint: "Must not break existing workflow creation"
→ Check: make test-integration-datastorage
✓ All existing tests pass

Constraint: "Performance: validation <10ms per workflow"
→ Note: Will validate in REFACTOR phase benchmarks
```

**Decision**:
- ✅ All technical requirements met
- ✅ Integration verified
- ✅ Constraints respected
- ✅ Proceed to REFACTOR

---

### Checkpoint 3: REFACTOR Phase (Final)

**Supervisor performs comprehensive review**:

```
Business Objective: "Prevent invalid workflows via label validation"
→ Validation: End-to-end test with invalid workflow
✓ Workflow rejected as expected
✓ Objective achieved

All Acceptance Criteria:
✓ AC-1: Rejects missing labels (verified)
✓ AC-2: Validates format (verified)
✓ AC-3: Clear error messages (verified)
✓ AC-4: Pre-persistence validation (verified)

All Technical Requirements:
✓ Location correct (pkg/datastorage/)
✓ Integration complete (handler + main app)
✓ Coverage adequate (85% > 70%)
✓ Patterns followed (matches existing code)

All Constraints:
✓ No breaking changes (all tests pass)
✓ Performance <10ms (benchmark: 3ms avg)
✓ All workflow types (tested with 5 types)

TDD Compliance:
✓ RED-GREEN-REFACTOR sequence followed
✓ Tests validate behavior not implementation
✓ No Skip() calls
✓ Integration in GREEN (not deferred)

Code Quality:
✓ No lint errors
✓ Errors handled properly
✓ No orphaned code
```

**Decision**: ✅ **APPROVED** - Ready for commit

---

## 🔄 Feedback Loop with BR

### When Worker Needs BR Clarification

```
Worker: "BR says 'clear error messages' but doesn't specify format.
         Should error include:
         A) Just label name
         B) Label name + expected format
         C) Label name + format + example"

Supervisor → User: "⏸️ PAUSED - Need BR clarification:

BR-WORKFLOW-197 AC-3 says 'clear error messages' but format not specified.

Options:
A) Minimal: 'Missing required label: environment'
B) Detailed: 'Missing required label: environment (format: key=value)'  
C) With example: 'Missing required label: environment (format: key=value, example: environment=production)'

Which level of detail does the business requirement intend?"

User → Supervisor: "Use Option B - format explanation without example"

Supervisor → Worker: "Use format: 'Missing required label: [name] (format: key=value)'"
```

---

## 📊 BR Completion Report

When supervisor approves work:

```
# BR-WORKFLOW-197 Implementation Complete

## Business Objective
✅ Prevent invalid workflows via label validation - ACHIEVED

## Acceptance Criteria
✅ AC-1: Rejects workflows missing mandatory labels
✅ AC-2: Validates label format (key=value)
✅ AC-3: Provides clear error messages  
✅ AC-4: Validation happens before persistence

## Technical Requirements
✅ Validator in pkg/datastorage/ - validator.go created
✅ Integration complete - handlers + main app
✅ Test coverage 85% (exceeds 70% requirement)
✅ Follows existing validation patterns

## Constraints Verified
✅ No breaking changes - all tests pass
✅ Performance <10ms - measured 3ms average
✅ All workflow types - tested 5 types

## Implementation Summary
Files Changed: 6 files
- pkg/datastorage/validator.go (new)
- pkg/datastorage/server/workflow_handlers.go (modified)
- cmd/datastorage/main.go (modified)
- test/unit/datastorage/validator_test.go (new)
- test/integration/datastorage/workflow_creation_test.go (modified)
- pkg/datastorage/config/validation_rules.yaml (new)

Tests Added: 8 tests (all behavior-focused)
Build Status: ✅ Success
Test Status: ✅ All Pass (85% coverage)

## Quality Assessment
TDD Compliance: ✅ Full RED-GREEN-REFACTOR
Code Quality: 92% confidence
Integration: ✅ Verified in main application
Standards: ✅ All Kubernaut guidelines followed

## Ready for Production
All acceptance criteria met. No blocking issues.

Next Steps:
- Commit changes
- Update BR-WORKFLOW-197 status to "Completed"
- Consider BR-WORKFLOW-198 (next in backlog)
```

---

## 🎯 Multi-BR Coordination

### Sequential Implementation

```
You: "Supervisor, implement these BRs in order:
     1. BR-WORKFLOW-197 (validation) - foundation
     2. BR-WORKFLOW-198 (search optimization) - depends on 197
     
     Ensure 197 fully approved before starting 198."
```

**Supervisor manages**:
- BR-197: Full cycle (RED → GREEN → REFACTOR → Approved)
- **THEN** BR-198: Starts after 197 complete

### Parallel Implementation (Different Areas)

```
You: "Supervisor, coordinate:
     - Worker A: BR-WORKFLOW-197 (validation)
     - Worker B: BR-API-045 (rate limiting)
     
     These are independent. Monitor both in parallel."
```

---

## 📝 BR Status Tracking

Supervisor can update BR status:

```
BR-WORKFLOW-197:
  Status: ✅ Approved → 🔄 In Progress → ✅ Completed
  Started: 2026-01-29 10:00
  Checkpoints: 
    - RED phase: 10:15 ✅
    - GREEN phase: 10:45 ✅
    - REFACTOR phase: 11:30 ✅
  Completed: 2026-01-29 11:45
  Confidence: 92%
```

---

## 🔗 Integration Checklist

For supervisor to work with your BRs:

- ✅ BR documents in accessible location
- ✅ Clear acceptance criteria in BRs
- ✅ Technical requirements specified
- ✅ Constraints documented
- ✅ BR status tracking (optional but helpful)

Ready to try? Point supervisor at your BR document!
