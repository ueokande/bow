package main

type Config struct {
	ContextName string
	Namespace   string
	KubeConfig  string
	Query       string
	Command     []string
}
