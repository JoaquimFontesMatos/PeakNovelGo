package utils

import "os/exec"

type ScriptExecutor interface {
	ExecuteScript(script string, args ...string) ([]byte, error)
}

type ScriptError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type RealScriptExecutor struct{}

func (r *RealScriptExecutor) ExecuteScript(script string, args ...string) ([]byte, error) {
	cmd := exec.Command(script, args...)
	return cmd.CombinedOutput()
}
