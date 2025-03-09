package utils

import "os/exec"

// ScriptExecutor defines an interface for executing scripts with provided arguments and returning the output or an error.
type ScriptExecutor interface {
	ExecuteScript(script string, args ...string) ([]byte, error)
}

// ScriptError represents an error structure typically used for script execution errors.
// It contains an error message and an HTTP status code.
type ScriptError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// RealScriptExecutor is a concrete implementation for executing external scripts using the exec.Command API.
type RealScriptExecutor struct{}

// ExecuteScript runs the specified script with optional arguments and returns combined stdout and stderr output.
//
// Parameters:
//	- script string (start command for the script to run)
//	- args ...string (optional number of args that could be added to the start command)
//
// Returns:
//	- []byte (response of the script)
//	- error (error that occurred in the execution of the script)
func (r *RealScriptExecutor) ExecuteScript(script string, args ...string) ([]byte, error) {
	cmd := exec.Command(script, args...)
	return cmd.CombinedOutput()
}
