<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/coder</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/coder"><img src="https://pkg.go.dev/badge/github.com/togo-framework/coder.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/coder
```

<!-- /togo-header -->

A togo **exec** provider backed by [Coder](https://coder.com): runs each autopilot
job inside an **isolated, disposable Coder workspace** on your own infrastructure
instead of the local machine — so agents get governed, reproducible environments.

Registers into [autopilot](https://github.com/togo-framework/autopilot)'s `exec`
capability slot. Off by default — opt in per project:

```bash
togo install togo-framework/coder
togo provider:use exec coder
togo config:set CODER_TEMPLATE omni-dev
```

## How it works
`Acquire` → `coder create autopilot-<id> --template <CODER_TEMPLATE>`; the impl
(claude/omnigent) then runs **inside** the workspace via `coder ssh … -- bash -lc
'cd <repo> && …'`; `release` → `coder delete`. Nothing touches your dev tree.

## Config (dynamic — CLI / .env / settings)

| Key | Default | Purpose |
|---|---|---|
| `CODER_TEMPLATE` | — (required) | workspace template to provision |
| `CODER_REPO_DIR` | `~/workspace` | repo path inside the workspace |
| `CODER_BIN` | `coder` | coder CLI binary |
| `CODER_KEEP` | — | `1` keeps the workspace after use |
| `CODER_URL` / `CODER_TOKEN` | — | used by `togo provider:test` |
<!-- togo-sponsors -->
---

<div align="center">
  <h3>💎 Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><img src=".github/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
    &nbsp;&nbsp;&nbsp;&nbsp;
    <a href="https://one-studio.co"><img src=".github/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
