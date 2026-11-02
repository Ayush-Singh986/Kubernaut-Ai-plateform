# DD-REGISTRY-001: Container Registry Purpose Classification

**Date**: December 17, 2025
**Status**: ✅ **APPROVED**
**Version**: 1.0
**Authority**: ⭐ **AUTHORITATIVE** - Single source of truth for registry usage
**Related**: ADR-028 (Container Registry Policy)

---

## 🎯 **Context**

Kubernaut uses two Quay.io organizations for container image storage. Confusion exists about which registry to use for different purposes (development, staging, production). This document provides authoritative classification.

---

## ✅ **Decision**

**APPROVED**: Two-registry strategy with clear purpose classification:

| Registry | Purpose | Usage | Access |
|---|---|---|---|
| **`quay.io/Ayush-Singh986/`** | Development & Testing | Local development, integration tests, CI builds | Public/Private (development team) |
| **`quay.io/kubernaut/`** | Staging & Production | Staging environments, production deployments, releases | Under Kubernaut control |

---

## 📋 **Registry Specifications**

### **Registry 1: `quay.io/Ayush-Singh986/` (Development)**

**Purpose**: Development and testing

**Usage Scenarios**:
- Local development image builds
- CI/CD pipeline test images
- Integration test infrastructure
- Feature branch images
- Personal experimentation

**Image Naming Convention**:
```
quay.io/Ayush-Singh986/<service-name>:<tag>

Examples:
quay.io/Ayush-Singh986/gateway:dev-feature-123
quay.io/Ayush-Singh986/datastorage:test-abc123
quay.io/Ayush-Singh986/must-gather:dev-latest
```

**Access Control**:
- Public repositories: Read access for all
- Private repositories: Write access for development team

**Retention Policy**:
- Development images: 30 days
- Test images: 14 days
- Latest/stable tags: Indefinite

---

### **Registry 2: `quay.io/kubernaut/` (Staging & Production)**

**Purpose**: Staging and production deployments

**Usage Scenarios**:
- Staging environment deployments
- Production deployments
- Official releases (v1.0.0, v1.1.0, etc.)
- Customer-facing images
- Documented public images (e.g., must-gather)

**Image Naming Convention**:
```
quay.io/kubernaut/<service-name>:<tag>

Examples:
quay.io/kubernaut/gateway:v1.0.0
quay.io/kubernaut/datastorage:v1.0.0
quay.io/kubernaut/must-gather:latest
quay.io/kubernaut/must-gather:v1.0.0
```

**Access Control**:
- Public repositories: Read access for all
- Write access: Controlled by Kubernaut organization (CI/CD automation + release managers)

**Retention Policy**:
- Production images: Indefinite (immutable)
- Staging images: 90 days
- Release tags: Indefinite (never deleted)

---

## 🔄 **Image Promotion Workflow**

### **Development → Staging → Production**

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Development Build                                    │
│ Registry: quay.io/Ayush-Singh986/                                │
│ Tag: dev-feature-123, dev-abc123, test-xyz                  │
│ Purpose: Local testing, CI validation                       │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Feature complete + tests pass
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 2: Staging Build                                        │
│ Registry: quay.io/kubernaut/                                │
│ Tag: staging-v1.0.0-rc1, staging-latest                     │
│ Purpose: Pre-production validation                          │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Staging validation complete
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 3: Production Release                                   │
│ Registry: quay.io/kubernaut/                                │
│ Tag: v1.0.0, v1.0, v1, latest                               │
│ Purpose: Production deployment                              │
└─────────────────────────────────────────────────────────────┘
```

### **Promotion Triggers**

| Stage | Trigger | Promoted To | Authority |
|---|---|---|---|
| Development | Feature complete + tests pass | Staging | Release manager approval |
| Staging | Staging validation complete | Production | Release manager + QA approval |

---

## 🚀 **CI/CD Integration**

### **Development Builds** (CI Pipeline)

```yaml
# .github/workflows/dev-build.yml
name: Development Build
on:
  push:
    branches: [develop, feature/*]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build and push to development registry
        run: |
          docker build -t quay.io/Ayush-Singh986/${SERVICE}:dev-${GITHUB_SHA:0:7} .
          docker push quay.io/Ayush-Singh986/${SERVICE}:dev-${GITHUB_SHA:0:7}
```

### **Production Releases** (Release Pipeline)

```yaml
# .github/workflows/release.yml
name: Production Release
on:
  push:
    tags: ['v*.*.*']

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Build and push to production registry
        run: |
          docker build -t quay.io/kubernaut/${SERVICE}:${GITHUB_REF_NAME} .
          docker push quay.io/kubernaut/${SERVICE}:${GITHUB_REF_NAME}
          docker tag quay.io/kubernaut/${SERVICE}:${GITHUB_REF_NAME} quay.io/kubernaut/${SERVICE}:latest
          docker push quay.io/kubernaut/${SERVICE}:latest
```

---

## 📊 **Registry Selection Decision Tree**

```
Is this image for...
├─ Local development? → quay.io/Ayush-Singh986/
├─ Integration testing? → quay.io/Ayush-Singh986/
├─ CI/CD pipeline? → quay.io/Ayush-Singh986/
├─ Feature branch? → quay.io/Ayush-Singh986/
│
├─ Staging environment? → quay.io/kubernaut/
├─ Production deployment? → quay.io/kubernaut/
├─ Official release? → quay.io/kubernaut/
└─ Customer-facing documentation? → quay.io/kubernaut/
```

---

## 🔒 **Security & Access Control**

### **`quay.io/Ayush-Singh986/` (Development)**

**Access Levels**:
- **Read**: Public (all developers, CI/CD)
- **Write**: Development team members
- **Admin**: Repository owner (@Ayush-Singh986)

**Robot Accounts**:
- `ci-builder`: Write access for CI/CD pipelines
- `test-runner`: Read access for integration tests

### **`quay.io/kubernaut/` (Staging & Production)**

**Access Levels**:
- **Read**: Public (all users, customers)
- **Write**: CI/CD automation (production pipeline only)
- **Admin**: Kubernaut organization admins

**Robot Accounts**:
- `release-automation`: Write access for production releases only
- `staging-automation`: Write access for staging builds only

---

## 📝 **Image Tagging Strategy**

### **Development Registry** (`quay.io/Ayush-Singh986/`)

| Tag Pattern | Example | Purpose |
|---|---|---|
| `dev-<git-sha>` | `dev-abc123` | Development builds (commit-specific) |
| `test-<name>` | `test-integration` | Test-specific images |
| `feature-<name>` | `feature-auth` | Feature branch builds |
| `dev-latest` | `dev-latest` | Latest development build |

### **Production Registry** (`quay.io/kubernaut/`)

| Tag Pattern | Example | Purpose |
|---|---|---|
| `v<major>.<minor>.<patch>` | `v1.0.0` | Semantic versioning (immutable) |
| `v<major>.<minor>` | `v1.0` | Minor version tracking (updates) |
| `v<major>` | `v1` | Major version tracking (updates) |
| `latest` | `latest` | Latest stable release |
| `staging-<name>` | `staging-v1.0.0-rc1` | Staging candidates |

---

## 🔗 **Integration with ADR-028**

This DD extends ADR-028 (Container Registry Policy) with registry purpose classification:

| ADR-028 | DD-REGISTRY-001 |
|---|---|
| Approved registries | Purpose classification |
| Base image sources | Development vs production |
| Security scanning | Registry-specific policies |
| Air-gapped mirroring | Promotion workflow |

**Authority Hierarchy**:
1. **ADR-028**: Defines WHICH registries are approved
2. **DD-REGISTRY-001** (this doc): Defines WHEN to use each registry

---

## ✅ **Consequences**

### **Positive**

1. ✅ **Clear Purpose**: No ambiguity about which registry to use
2. ✅ **Controlled Production**: Production images protected by access control
3. ✅ **Development Freedom**: Developers can iterate freely in development registry
4. ✅ **Audit Trail**: Clear promotion path (dev → staging → prod)
5. ✅ **Cost Optimization**: Development images auto-expire (30 days), production images retained

### **Negative**

1. ⚠️ **Image Promotion Overhead**: Requires CI/CD pipeline automation
2. ⚠️ **Registry Duplication**: Same image may exist in both registries
3. ⚠️ **Access Management**: Two registries to manage (users, robot accounts)

### **Neutral**

1. 🔄 **Migration**: Existing images need clear purpose classification
2. 🔄 **Documentation**: All docs must specify correct registry

---

## 📊 **Approval**

**Status**: ✅ **APPROVED** (2025-12-17)
**Confidence**: **100%**
**Priority**: **P1 - Required for v1.0 Release**

**Rationale**:
1. Clarifies existing confusion about registry usage
2. Protects production images with controlled access
3. Enables development team iteration without production risk
4. Aligns with industry best practices (dev/staging/prod separation)

---

## 🔗 **Related Documents**

- **ADR-028**: Container Registry Policy (defines approved registries)
- **ADR-027**: Multi-Architecture Build Strategy (image build patterns)
- **BR-PLATFORM-001**: Must-Gather Diagnostic Collection (references production registry)

---

## 📝 **Update History**

| Version | Date | Change | Author |
|---|---|---|---|
| 1.0 | 2025-12-17 | Initial version | AI Assistant per user clarification |

---

**Next Steps**:
1. Update all documentation to reference `quay.io/kubernaut/` for production images
2. Update CI/CD pipelines to implement promotion workflow
3. Configure robot accounts with appropriate access levels
4. Create image promotion automation scripts

---

**Authority**: This document is the authoritative source for registry purpose classification. All services, documentation, and scripts MUST reference this DD when determining which registry to use.

