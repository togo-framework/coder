// Package coder is a togo exec provider backed by Coder (https://coder.com):
// it runs each autopilot job inside an isolated, disposable Coder workspace on
// your own infrastructure instead of the local machine.
//
// It registers into autopilot's `exec` capability slot. Opt in with
// `togo provider:use exec coder` (or TOGO_EXEC_PROVIDER=coder); it is off by
// default so `local` stays the default. Fully dynamic over CLI / .env / settings:
//
//	CODER_TEMPLATE   workspace template to provision (required)
//	CODER_REPO_DIR   repo path inside the workspace   (default ~/workspace)
//	CODER_BIN        coder CLI binary                 (default "coder")
//	CODER_KEEP       "1" keeps the workspace after use (default: delete)
//	CODER_URL/CODER_TOKEN   used by `togo provider:test`
package coder

import (
	"github.com/togo-framework/providers"
	"github.com/togo-framework/togo"
)

// Name is the plugin's stable identifier and its `exec` backend name.
const Name = "coder"

func init() {
	// PriorityService (before autopilot registers its local default at a later
	// priority): install ourselves into the exec slot only when selected.
	togo.RegisterProviderFunc(Name, togo.PriorityService, func(k *togo.Kernel) error {
		providers.Use(k, providers.CapExecute, Name, newExec(k), false)
		if k.Log != nil {
			k.Log.Info("plugin active", "plugin", Name)
		}
		return nil
	})
}
