package bow

// Config is a bow config
type Config struct {
	ContextName string
	Namespace   string
	KubeConfig  string
	Query       string
	Command     []string
	NoHosts     bool
}
