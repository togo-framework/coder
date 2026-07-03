<!-- togo-brand -->
<p align="center">
  <img src=".github/assets/togo-mark.svg" width="96" alt="togo" />
</p>
<h1 align="center">coder</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a> — the full-stack Go + React framework</sub></p>

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

MIT © fadymondy
