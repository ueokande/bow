# bow

Bow allow you to exec a command on multiple container for Kubernetes.

![Screenshot](screenshot.gif)

## Install

Download a latest version of the binary from [releases][], or use `go get` as follows:

```console
$ go get -u github.com/ueokande/bow
```

## Usage

```console
$ bow [OPTIONS] POD_SELECTOR -- COMMAND ARGS...

Flags:
  -c, --container string    Container name when multiple containers in pod
      --kubeconfig string   Path to kubeconfig file to use
  -n, --namespace string    Kubernetes namespace to use.
      --no-hosts            Do not print hosts
```

## License

MIT

[releases]: https://github.com/ueokande/bow/releases
