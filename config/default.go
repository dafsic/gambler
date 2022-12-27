package config

// DefaultGamblerCfg returns the default config
func DefaultGamblerCfg() *GamblerCfg {
	return &GamblerCfg{
		LogLevel: "info",
		TrxNode:  "127.0.0.1:8090",
		Addr:     "",
		Pk:       "",
		Pool:     "",
		Refund:   "",
	}
}
