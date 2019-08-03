package bow

// Config is a bow config
type Config struct {
	Namespace  string
	KubeConfig string
	Query      string
	Container  string
	Command    []string
	NoHosts    bool
}
