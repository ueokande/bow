package main

type Config struct {
	ContextName string
	Namespace   string
	KubeConfig  string
	Query       Query
}

type Query []string
