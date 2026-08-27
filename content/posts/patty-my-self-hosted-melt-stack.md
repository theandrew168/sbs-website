---
date: 2026-08-22
title: "Patty: My Self-Hosted MELT Stack"
slug: "patty-my-self-hosted-melt-stack"
tags: ["VictoriaMetrics", "VictoriaLogs", "VictoriaTraces"]
draft: true
---

Talk about the Victoria stack: metrics, logs, and traces.
Ansible roles for each server.
Ansible roles for each client (agent).
Metrics are pushed (could pull, but server and DB metrics are local-only).
Logs are pushed via systemd-journal-remote (custom namespace for filtering).
Traces are pushed (still need to figure this out).

The MELT server runs all three programs, listening on both localhost (for Grafana) the internal VPC (10.136.0.4, for agent pushes).
A mounted volume holds all of the data (separate subdir per service).
Grafana accesses the data via the localhost listeners and various connectors (VictoriaMetrics, VictoriaLogs, and Jaeger).
Grafana's UI listens only on localhost, exposed to the internet (with TLS) via Caddy.
Caddy listens on 80/443 with auto-redirects for non-TLS conns.
UFW is configured to allow all traffic on 80/443 and internal traffic only (10.136.0.0/16) on the various service ports (8428, 9428, 10428).

## Open Questions

Victoria stack (metrics, logs, traces) vs Prom + Jaeger?

Separate logs w/ correlation IDs (trace_id and span_id) or embedded logs as Span Events?
In OpenTelemetry (OTel), correlating logs and traces using IDs (TraceId and SpanId) means keeping logs as a distinct telemetry signal that references a trace, while embedded logs as span events means injecting log-like data directly into the trace payload itself.
Recommendation: Use ID-based correlation for standard application logging.
The OpenTelemetry community has officially moved toward deprecating the Span Events API for general log-like events in favor of the formal Logs API.