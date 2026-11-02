# Test Naming Conventions

This document establishes clear naming conventions for different types of tests to ensure proper separation and organization.

## 📁 **Directory Structure**

### Business Requirement Tests
```
test/business-requirements/
├── ai/                    # AI and ML business requirements
│   └── *_business_test.go
├── workflow/              # Workflow execution business requirements
│   └── *_business_test.go
├── infrastructure/        # Platform and infrastructure business requirements
│   └── *_business_test.go
├── integration/          # Cross-component business requirements
│   └── *_business_test.go
└── shared/               # Business test utilities
    └── business_test_suite.go
```

### Unit Tests (Implementation Tests)
```
pkg/{component}/
├── component.go           # Production code
├── component_test.go      # Unit tests
├── interfaces.go          # Production interfaces
└── types.go              # Production types
```

## 🏷️ **File Naming Conventions**

| Test Type | File Pattern | Package Pattern | Example |
|-----------|-------------|-----------------|---------|
| **Business Requirements** | `*_business_test.go` | `{component}_business_test` | `effectiveness_assessment_business_test.go` |
| **Unit Tests** | `*_test.go` | `{component}_test` | `effectiveness_assessment_test.go` |
| **Integration Tests** | `*_integration_test.go` | `{component}_integration_test` | `ai_pipeline_integration_test.go` |

## 🧪 **Test Suite Naming**

### Business Requirement Tests
```go
// Pattern: {Component} - Business Requirements Validation
var _ = Describe("AI Effectiveness Assessment - Business Requirements Validation", func() {
    // Business requirement tests focus on outcomes
})
```

### Unit Tests
```go
// Pattern: {Component}.{Method} or {Component} Implementation
var _ = Describe("EffectivenessAssessor.ProcessPendingAssessments", func() {
    // Unit tests focus on implementation correctness
})
```

## 📋 **Business Requirement ID Format**

### ID Structure: `BR-{COMPONENT}-{NUMBER}`

| Component Category | ID Pattern | Example | Description |
|-------------------|------------|---------|-------------|
| **AI & ML** | `BR-AI-###` | `BR-AI-001` | AI systems, learning, recommendations |
| **Workflows** | `BR-WF-###` | `BR-WF-001` | Workflow execution, orchestration |
| **Infrastructure** | `BR-INF-###` | `BR-INF-001` | Platform, monitoring, performance |
| **Integration** | `BR-INT-###` | `BR-INT-001` | Cross-component, end-to-end flows |
| **Vector Database** | `BR-VDB-###` | `BR-VDB-001` | Vector operations, embeddings |
| **Orchestration** | `BR-ORK-###` | `BR-ORK-001` | Multi-service orchestration |

### Business Requirement Test Structure
```go
Describe("BR-AI-001: System Must Learn From Action Failures", func() {
    It("should reduce confidence for actions that fail repeatedly", func() {
        // BUSINESS REQUIREMENT: Actions with <50% success rate get reduced confidence

        // Given: Business scenario setup
        // When: Business action executed
        // Then: Business outcome validated
        // And: Business value measured
    })
})
```

## 🎯 **Context and Description Patterns**

### Business Requirement Tests
- **Focus**: Business outcomes and value
- **Language**: Stakeholder-friendly descriptions
- **Validation**: Business metrics and acceptance criteria

```go
Context("when system processes high-failure actions", func() {
    It("should reduce confidence below 50% threshold", func() {
        // Test validates business requirement
    })
})
```

### Unit Tests
- **Focus**: Implementation correctness
- **Language**: Technical, developer-focused
- **Validation**: Code behavior and edge cases

```go
Context("when ProcessPendingAssessments receives invalid input", func() {
    It("should return validation error", func() {
        // Test validates implementation behavior
    })
})
```

## 📊 **Measurement and Metrics**

### Business Requirement Tests
- **Metrics**: Business KPIs (accuracy %, cost reduction, time savings)
- **Thresholds**: Business-defined success criteria
- **Reporting**: Stakeholder-friendly outcomes

```go
// Business value measurement
finalAccuracy := calculateSystemAccuracy()
improvementPct := ((finalAccuracy - baselineAccuracy) / baselineAccuracy) * 100

Expect(improvementPct).To(BeNumerically(">=", 25),
    "System should achieve 25% accuracy improvement")
```

### Unit Tests
- **Metrics**: Technical metrics (execution time, memory usage, coverage)
- **Thresholds**: Technical requirements
- **Reporting**: Developer-focused diagnostics

```go
// Technical validation
result, err := processor.ProcessPendingAssessments(ctx)
Expect(err).ToNot(HaveOccurred())
Expect(result).ToNot(BeNil())
```

## 🔧 **Package Declaration Patterns**

### Business Requirement Tests
```go
package ai_business_test  // Component + business_test suffix

import (
    // Production code under test
    "github.com/Ayush-Singh986/kubernaut/pkg/ai/insights"

    // Shared business test utilities
    "github.com/Ayush-Singh986/kubernaut/test/business-requirements/shared"

    // Centralized mocks
    "github.com/Ayush-Singh986/kubernaut/pkg/testutil/mocks"
)
```

### Unit Tests
```go
package insights_test  // Component_test suffix

import (
    // Production code under test
    "github.com/Ayush-Singh986/kubernaut/pkg/ai/insights"

    // Centralized mocks for minimal dependencies
    "github.com/Ayush-Singh986/kubernaut/pkg/testutil/mocks"
)
```

## 🚦 **Migration Guidelines**

### Converting Unit Tests to Business Requirements Tests
1. **Identify business value**: Does this test validate a business requirement?
2. **Move to business-requirements directory**: Use proper naming convention
3. **Update package declaration**: Use `{component}_business_test` pattern
4. **Rewrite test descriptions**: Focus on business outcomes
5. **Add business metrics**: Measure business value, not just technical success
6. **Use business test suite**: Import and use `shared.BusinessTestSuite`

### Example Migration
```go
// BEFORE (Unit Test)
It("should process assessments without error", func() {
    err := assessor.ProcessPendingAssessments(ctx)
    Expect(err).ToNot(HaveOccurred())
})

// AFTER (Business Requirement Test)
It("should improve recommendation accuracy by 25% through learning", func() {
    // BUSINESS REQUIREMENT: System must learn and improve over time

    baselineAccuracy := 0.6
    targetImprovement := 0.25

    // Execute business scenario
    err := assessor.ProcessPendingAssessments(ctx)
    Expect(err).ToNot(HaveOccurred())

    // Measure business outcome
    currentAccuracy := calculateAccuracy()
    improvementPct := (currentAccuracy - baselineAccuracy) / baselineAccuracy

    Expect(improvementPct).To(BeNumerically(">=", targetImprovement),
        "System should achieve 25% accuracy improvement through learning")

    // Log business value
    suite.LogBusinessOutcome("BR-AI-001", businessMetric, improvementPct >= targetImprovement)
})
```

## ✅ **Validation Checklist**

### Business Requirement Test Checklist
- [ ] File ends with `_business_test.go`
- [ ] Package uses `{component}_business_test` pattern
- [ ] Test ID follows `BR-{COMPONENT}-{NUMBER}` format
- [ ] Test description focuses on business outcome
- [ ] Test validates business metrics/KPIs
- [ ] Test uses business test suite utilities
- [ ] Test logs business outcomes for stakeholders

### Unit Test Checklist
- [ ] File ends with `_test.go` (not `_business_test.go`)
- [ ] Package uses `{component}_test` pattern
- [ ] Test description focuses on implementation behavior
- [ ] Test validates technical correctness
- [ ] Test uses minimal, focused mocks
- [ ] Test covers edge cases and error conditions
- [ ] Test provides fast feedback for developers
