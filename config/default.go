package config

// DefaultAssistantNode returns the default config
func DefaultAssistantNode() *AssistantNode {
	return &AssistantNode{
		LogLevel: "info",
		API: API{
			Address: "127.0.0.1:6660",
			Token:   "",
			Timeout: 30,
		},
		MinerAPI: API{
			Address: "192.168.28.136:2345",
			Token:   "",
			Timeout: 30,
		},
		DaemonAPI: API{
			Address: "192.168.28.174:1234",
			Token:   "",
			Timeout: 30,
		},
	}
}
