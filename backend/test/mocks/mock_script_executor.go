package mocks

import (
	"backend/internal/utils"
	"encoding/json"
)

type MockScriptExecutorNetworkDown struct{}

func (m *MockScriptExecutorNetworkDown) ExecuteScript(script string, args ...string) ([]byte, error) {
	var scriptError utils.ScriptError

	scriptError.Status = 503
	scriptError.Error = "Network connection down"
	message, err := json.Marshal(scriptError)
	if err != nil {
		return nil, err
	}

	return []byte(message), nil
}

type MockScriptExecutorSourceWebsiteDown struct{}

func (m *MockScriptExecutorSourceWebsiteDown) ExecuteScript(script string, args ...string) ([]byte, error) {
	var scriptError utils.ScriptError

	scriptError.Status = 503
	scriptError.Error ="Source website down"
	message, err := json.Marshal(scriptError)
	if err != nil {
		return nil, err
	}

	return []byte(message), nil
}
