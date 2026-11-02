# Kubernaut

**AIOps Platform for Intelligent Kubernetes Remediation**

[![Go Reference](https://pkg.go.dev/badge/github.com/jordigilh/kubernaut.svg)](https://pkg.go.dev/github.com/jordigilh/kubernaut)
[![Go Version](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org/dl/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.32+-blue.svg)](https://kubernetes.io/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI](https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/actions/workflows/ci-pipeline.yml/badge.svg)](https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/actions/workflows/ci-pipeline.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Ayush-Singh986/Kubernaut-Ai-plateform/badge)](https://scorecard.dev/viewer/?uri=github.com/Ayush-Singh986/Kubernaut-Ai-plateform)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13485/badge)](https://www.bestpractices.dev/projects/13485)
[![Latest Release](https://img.shields.io/github/v/release/Ayush-Singh986/Kubernaut-Ai-plateform)](https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/releases/latest)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

> This repository is a fork of [Jordi Gil's Kubernaut project](https://github.com/jordigilh/kubernaut). The upstream Go module path and copyright notices are retained.

Kubernaut closes the loop from Kubernetes alert to automated remediation. It operates in two modes: **autonomously** — detecting signals, investigating root causes, and executing fixes end-to-end without human involvement — and **interactively** — letting operators join an in-progress investigation via MCP or A2A, guide the agent, and approve remediations in real time. The LLM-powered agent uses native Go client-go bindings against the Kubernetes API, Prometheus, and log endpoints to investigate, select a remediation workflow, and execute the fix — or escalate to a human with a full RCA when it can't. See the **[full documentation](https://Ayush-Singh986.github.io/kubernaut-docs/)** for architecture, installation, and usage guides.

<div align="center">
  <video src="https://github.com/user-attachments/assets/b95290db-412b-4d6d-81b8-f766ef4657e2" controls width="100%"></video>
</div>

<details>
<summary>See autonomous mode demo (non-interactive)</summary>
<p align="center">
  <img src="https://raw.githubusercontent.com/Ayush-Singh986/kubernaut-demo-scenarios/main/scenarios/crashloop/crashloop-lite.gif" alt="CrashLoopBackOff demo — autonomous alert-to-fix" width="800"/>
</p>
<p align="center">
  This is 1 of 37 runnable scenarios in <a href="https://github.com/Ayush-Singh986/kubernaut-demo-scenarios"><strong>kubernaut-demo-scenarios</strong></a> — clone it and run any of them against your own cluster.
</p>
</details>

---

## Why

Kubernetes operators spend hours manually triaging alerts, diagnosing root causes from scattered logs and metrics, and executing remediation steps from runbooks that drift out of date. The response depends on tribal knowledge, human availability, and often happens at 3am.

Rule-based remediation tools help with known, deterministic problems — "if X, do Y." But when the same symptom has multiple root causes, or the right fix depends on context the rule can't see, they fall short.

Kubernaut bridges that gap. It uses an LLM agent that investigates the actual root cause through native Go bindings against the Kubernetes API and observability stack, selects the right remediation from a workflow catalog, executes it, and verifies the fix worked — escalating to humans only when it should. Rule-based tools are thermostats. Kubernaut is a diagnostician that also adjusts the thermostat.

**[Why Kubernaut? — full comparison with rule-based tools](https://Ayush-Singh986.github.io/kubernaut-docs/latest/getting-started/why-kubernaut/)**

---

## What It Can Triage & Fix

Validated against [37 runnable scenarios](https://github.com/Ayush-Singh986/kubernaut-demo-scenarios) covering real-world Kubernetes failure modes:

| Category | Examples |
|---|---|
| **Workload failures** | Crash loops, stuck rollouts, memory leaks/escalation, SLO error-budget burn |
| **Resource & capacity** | HPA maxed out, PDB deadlocks, PVC capacity exhaustion (forecasted), DB connection pool saturation |
| **Infrastructure** | Node failures, PVC/StatefulSet binding failures, scheduling taints, OCP operator/build/RBAC/SCC violations |
| **Networking & delivery** | NetworkPolicy misconfigurations, Istio authorization failures, GitOps drift, cert-manager failures |
| **Complex root-cause chains** | Cross-namespace dependency tracing, cascading multi-service failures, prioritizing causation over severity, filtering unrelated noise |
| **Safety & reasoning robustness** | Prompt-injection detection (shadow agent), resistance to misleading alert descriptions, duplicate-alert suppression |

See the [full scenario catalog](https://github.com/Ayush-Singh986/kubernaut-demo-scenarios/blob/main/docs/scenarios.md) for exact alerts, fault injection, and remediation workflows per scenario.

---

## What It Does

- **Detects** — Ingests Prometheus AlertManager alerts and Kubernetes Events, validates resource scope, and deduplicates by fingerprint
- **Triages** — Resolves signal severity through a multi-tier pipeline (firing alerts, Prometheus rule evaluation, LLM-based triage) and derives grounded signal names from infrastructure context
- **Investigates** — Performs live root cause analysis using Kubernetes inspection tools, configurable observability toolsets (Prometheus, etc.), and remediation history. Runs **autonomously** end-to-end, or **interactively** with an operator guiding the investigation in real time via MCP tools or A2A sessions
- **Integrates** — Exposes MCP and A2A (Agent-to-Agent) protocol endpoints through the API Frontend, enabling external agents, UIs, and automation to interact with Kubernaut via OIDC-authenticated sessions. Operators can take over autonomous sessions mid-flight, review findings, and approve next steps
- **Remediates** — Selects and executes a workflow from a searchable catalog via Tekton Pipelines, Kubernetes Jobs, or Ansible (AWX/AAP), with optional human approval gates
- **Closes the loop** — Notifies the team (Slack, webhook), evaluates whether the fix worked via health checks, alert resolution, and spec hash drift detection, and feeds effectiveness scores back into future investigations

<details>
<summary>Architecture</summary>

![Kubernaut Remediation Pipeline](https://raw.githubusercontent.com/Ayush-Singh986/kubernaut-docs/main/docs/assets/images/pipeline-phases.svg)

See the interactive version with per-phase detail: [How It Works](https://Ayush-Singh986.github.io/kubernaut-docs/latest/#how-it-works).

</details>

<details>
<summary>Services</summary>

| Service | Path | Description |
|---|---|---|
| **Gateway** | `cmd/gateway` | Signal ingestion — AlertManager webhooks and Kubernetes Events |
| **Signal Processing** | `cmd/signalprocessing` | Signal enrichment, deduplication, and routing |
| **Remediation Orchestrator** | `cmd/remediationorchestrator` | CRD lifecycle orchestration across the pipeline |
| **AI Analysis** | `cmd/aianalysis` | Investigation controller — dispatches to Kubernaut Agent |
| **Kubernaut Agent** | `cmd/kubernautagent` | LLM-powered RCA, workflow selection, and MCP tool execution |
| **API Frontend** | `cmd/apifrontend` | External protocol layer — MCP, A2A, OIDC auth, severity triage |
| **Workflow Execution** | `cmd/workflowexecution` | Tekton Pipeline / Job / Ansible execution engine |
| **Data Storage** | `cmd/datastorage` | Workflow catalog, audit trail, and persistence (PostgreSQL) |
| **Notification** | `cmd/notification` | Slack, webhook, and console notification delivery |
| **Effectiveness Monitor** | `cmd/effectivenessmonitor` | Post-remediation health checks and effectiveness scoring |
| **Auth Webhook** | `cmd/authwebhook` | Kubernetes authentication webhook for service identity |
| **Fleet Metadata Cache** | `cmd/fleetmetadatacache` | Multi-cluster metadata cache for cross-cluster workflow targeting |

</details>

---

## Try It Out

Don't have a cluster yet? Follow the
[Fleet Demo Quick Start](docs/operations/deployment/FLEET_DEMO_QUICKSTART.md) for the
operator-backed qualification path or the user-managed Fleet setup. The repository no
longer ships a one-command application demo installer.

The former one-command local and Fleet demo helpers were removed with the application
Helm deployment path. For a disposable two-cluster qualification environment, run the
repository's real Fleet journey with `PRESERVE_E2E_CLUSTER=true make test-e2e-fleet`.
For a user-managed deployment, follow the [Kubernaut Operator installation guide](https://github.com/Ayush-Singh986/kubernaut-operator/blob/main/docs/installation/00-quickstart.md),
then use the [Autonomous](docs/operations/deployment/AUTONOMOUS_MODE_SETUP_GUIDE.md),
[Interactive](docs/operations/deployment/INTERACTIVE_MODE_SETUP_GUIDE.md), or
[Fleet setup](docs/operations/deployment/FLEET_SETUP_GUIDE.md) guide.

---

## Installation

Kubernaut is deployed through the [Kubernaut Operator installation guide](https://github.com/Ayush-Singh986/kubernaut-operator/blob/main/docs/installation/00-quickstart.md).
The application Helm chart is deprecated and is no longer shipped from this repository.

Once installed, pick a setup guide for how you want to trigger remediation — the three are
independent and can be combined on the same install:

| Mode | Trigger | Guide |
|---|---|---|
| Autonomous | Alert-driven (Prometheus/AlertManager → Gateway) | [Autonomous Mode Setup Guide](docs/operations/deployment/AUTONOMOUS_MODE_SETUP_GUIDE.md) |
| Interactive | Chat-driven (Console/APIFrontend, human-in-the-loop) | [Interactive Mode Setup Guide](docs/operations/deployment/INTERACTIVE_MODE_SETUP_GUIDE.md) |
| Fleet | Multi-cluster (hub + spoke, layers onto either mode above) | [Fleet Setup Guide](docs/operations/deployment/FLEET_SETUP_GUIDE.md) |

---

## End-to-End Flow and Deployment

### How a remediation runs

In autonomous mode, Prometheus/AlertManager sends an alert to the **Gateway**. The
Gateway validates and records the signal, then creates a `RemediationRequest`.
**Signal Processing** enriches and deduplicates it, and the **Remediation Orchestrator**
coordinates the lifecycle through investigation and execution. **AI Analysis** dispatches
the **Kubernaut Agent**, which investigates the affected resources using Kubernetes and
configured observability data, identifies a likely root cause, and selects an eligible
workflow. **Workflow Execution** runs that workflow using the configured backend, such
as Tekton, a Kubernetes Job, or Ansible. Kubernaut then sends notifications, checks
whether the remediation worked, and records the outcome for audit and future analysis.

Interactive mode starts an investigation through the API Frontend instead of an alert;
an operator can guide the investigation and approve actions. Fleet mode extends either
entry point to target workloads on managed clusters. These modes can be combined;
their prerequisites and configuration are described in the mode-specific guides above.

### Deployment steps

Kubernaut's supported deployment is managed by the **Kubernaut Operator**. This
repository no longer ships an application Helm chart or a one-command installer. Use
the [Operator installation guide](https://github.com/Ayush-Singh986/kubernaut-operator/blob/main/docs/installation/00-quickstart.md)
for the exact commands and configuration.

1. **Prepare a Kubernetes cluster** that meets the project's supported version
   requirements, and ensure you have permission to install cluster resources. Arrange
   persistent storage for the services that need it and an LLM provider credential.
   Autonomous mode also needs Prometheus/AlertManager; workflow backends and their
   credentials depend on the workflows you plan to run.
2. **Install the Kubernaut Operator** by following its quick-start guide. Install and
   configure the dependencies it documents, including PostgreSQL, Valkey, and any
   required credentials and policy ConfigMaps.
3. **Configure the Kubernaut custom resource** (`kubernaut.ai/v1alpha2`) with the
   namespace, service settings, Secret references, policies, and enabled features for
   your environment, then apply it as described in the Operator guide. The Operator
   reconciles this resource and deploys the Kubernaut services.
4. **Set up an entry point**: configure AlertManager and Gateway using the
   [Autonomous Mode guide](docs/operations/deployment/AUTONOMOUS_MODE_SETUP_GUIDE.md),
   or configure the API Frontend and authentication using the
   [Interactive Mode guide](docs/operations/deployment/INTERACTIVE_MODE_SETUP_GUIDE.md).
   For multi-cluster targeting, follow the
   [Fleet Setup guide](docs/operations/deployment/FLEET_SETUP_GUIDE.md) as well.
5. **Verify and exercise the installation** by checking that the Operator and Kubernaut
   workloads are ready in their namespace, then send a test alert or start an
   interactive investigation. Confirm the investigation, workflow execution, and
   resulting notification or effectiveness status in the configured environment.

For a disposable two-cluster qualification environment, see the
[Fleet Demo Quick Start](docs/operations/deployment/FLEET_DEMO_QUICKSTART.md). It
documents the requirements and command for the repository's operator-backed test
journey; it is not an application installer.

---

## Signed & Verified

Every container image published to `quay.io/kubernaut-ai` is keylessly signed with
[Cosign](https://github.com/sigstore/cosign) (Sigstore/Fulcio/Rekor) and carries a
[SLSA build provenance](https://slsa.dev/) attestation plus a CycloneDX SBOM attestation,
generated entirely inside GitHub Actions via OIDC — no long-lived signing keys involved.
See [`.github/workflows/release.yml`](.github/workflows/release.yml) for the full pipeline.

Verify any image before you pull it:

```bash
# Resolve the latest published release. GitHub's `releases/latest` endpoint always
# excludes drafts and pre-releases, so this never resolves to an -rcN build.
VERSION=$(curl -fsSL https://api.github.com/repos/Ayush-Singh986/Kubernaut-Ai-plateform/releases/latest | grep -o '"tag_name": *"[^"]*"' | cut -d'"' -f4)

# Verify the signature (keyless, tied to this repo's release workflow)
cosign verify quay.io/kubernaut-ai/gateway:${VERSION} \
  --certificate-identity-regexp "https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/.github/workflows/release.yml@.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"

# Verify the SLSA build provenance attestation
cosign verify-attestation quay.io/kubernaut-ai/gateway:${VERSION} \
  --type slsaprovenance \
  --certificate-identity-regexp "https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/.github/workflows/release.yml@.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
```

Security posture is tracked continuously via [OpenSSF Scorecard](https://scorecard.dev/viewer/?uri=github.com/Ayush-Singh986/Kubernaut-Ai-plateform)
(badge above) and [CodeQL](.github/workflows/codeql.yml) static analysis. Found a vulnerability? See [SECURITY.md](SECURITY.md) for our disclosure process.

---

## Documentation

| Resource | Link |
|---|---|
| **User & Operator Guide** | [Ayush-Singh986.github.io/kubernaut-docs](https://Ayush-Singh986.github.io/kubernaut-docs/) |
| **Architecture Overview** | [Architecture](https://Ayush-Singh986.github.io/kubernaut-docs/latest/getting-started/architecture-overview/) |
| **Developer Guide** | [docs/DEVELOPER_GUIDE.md](docs/DEVELOPER_GUIDE.md) |
| **Must-Gather Diagnostics** | [cmd/must-gather/README.md](cmd/must-gather/README.md) |
| **DevConf.US presentation — September 24, 2026** | [When Your Kubernetes Cluster Fixes Itself at 3am](https://Ayush-Singh986.github.io/kubernaut-devconf-reference/) |

---

## Related Repositories

| Repository | Description |
|---|---|
| [kubernaut-operator](https://github.com/Ayush-Singh986/kubernaut-operator) | Operator and Helm bootstrap chart — deploys and manages Kubernaut services on generic Kubernetes, Kind, and OpenShift |
| [kubernaut-console](https://github.com/Ayush-Singh986/kubernaut-console) | Kubernaut Console — web UI for interactive investigation and remediation |
| [kubernaut-docs](https://github.com/Ayush-Singh986/kubernaut-docs) | Documentation website (MkDocs Material) |
| [kubernaut-demo-scenarios](https://github.com/Ayush-Singh986/kubernaut-demo-scenarios) | Demo scenarios, scripts, and recordings |

---

## Roadmap

### v1.6 — Fleet Operations (release candidate)

- **Fleet operations** — Multi-cluster remediation orchestration through a pluggable scope-checking adapter: Red Hat ACM/OCM for ACM shops, or a control-plane-free mode backed by the built-in Fleet Metadata Cache (FMC) for GitOps and standalone clusters — with Rancher and Clusterpedia adapters architected for other vendor fleet platforms ([#54](https://github.com/jordigilh/kubernaut/issues/54))
- **ServiceNow incident triage** (v1.6.1) — Consume ServiceNow incidents as signals through the API Frontend, enabling Kubernaut to investigate and remediate ITSM tickets alongside Kubernetes alerts ([#1338](https://github.com/jordigilh/kubernaut/issues/1338))

**[Full roadmap](docs/roadmap/ROADMAP.md)** — for released features, see the [CHANGELOG](CHANGELOG.md).

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. In short: create a feature branch, implement with tests, update docs, and open a PR.

---

## License

Apache License 2.0 — see [LICENSE](LICENSE).

---

**Website**: [kubernaut.ai](https://kubernaut.ai) · **Issues**: [GitHub Issues](https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/issues) · **Discussions**: [GitHub Discussions](https://github.com/Ayush-Singh986/Kubernaut-Ai-plateform/discussions)

**Kubernaut** — From alert to remediation, intelligently.
