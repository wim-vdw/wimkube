# wimkube

An interactive Kubernetes CLI tool built with Go that simplifies common Kubernetes operations through an intuitive
command-line interface.

## Features

- **Context Management**: Switch between and manage Kubernetes contexts
- **Namespace Management**: View and switch between namespaces
- **Pod Operations**: List pods, view containers, execute interactive shells, and retrieve container logs
- **Interactive Menus**: User-friendly interactive prompts for all operations
- **Direct Commands**: Support for both interactive and direct command execution

## Installation

### Prerequisites

- Go 1.26 or later
- Access to a Kubernetes cluster
- kubectl configured with valid kubeconfig

### Homebrew (macOS)

```bash
brew install wim-vdw/tap/wimkube
```

Or using a two-step tap:

```bash
brew tap wim-vdw/tap
brew install wimkube
```

Supports both Apple Silicon (arm64) and Intel (amd64) Macs.

### Build from Source

```bash
git clone https://github.com/wim-vdw/wimkube.git
cd wimkube
go build -o wimkube
```

### Install

```bash
go install github.com/wim-vdw/wimkube@latest
```

## Usage

### Global Flags

- `--kubeconfig`: Path to the kubeconfig file (default: `~/.kube/config`)
- `-t, --request-timeout`: Timeout in seconds for Kubernetes API requests (default: 30)
- `-h, --help`: Display help message
- `-v, --version`: Display version information

### Version Information

**Display detailed version information:**

```bash
wimkube version
```

This shows the version, Go version, Git commit, build time, and OS/Arch.

### Context Management

**Interactive menu:**

```bash
wimkube context
```

**List all contexts:**

```bash
wimkube context list
```

**Get current context:**

```bash
wimkube context get
```

**Set current context:**

```bash
wimkube context set <context-name>
```

### Namespace Management

**Interactive menu:**

```bash
wimkube namespace
```

**List all namespaces:**

```bash
wimkube namespace list
```

**Get current namespace:**

```bash
wimkube namespace get
```

**Set current namespace:**

```bash
wimkube namespace set <namespace-name>
```

### Pod Management

**Interactive menu:**

```bash
wimkube pod
```

**List all pods in current namespace:**

```bash
wimkube pod list
```

**List containers in a pod:**

```bash
wimkube pod list-containers <pod-name>
```

**Execute interactive shell in a container:**

```bash
wimkube pod exec <pod-name> <container-name>
```

**Get the logs of a container:**

```bash
wimkube pod logs <pod-name> <container-name>
```

## Examples

### Switch to a different context

```bash
# Interactive - follow the prompts
wimkube context

# Direct command
wimkube context set production-cluster
```

### Change namespace

```bash
# Interactive - follow the prompts
wimkube namespace

# Direct command
wimkube namespace set default
```

### Access a pod container

```bash
# Interactive - follow the prompts
wimkube pod

# Direct command
wimkube pod exec my-pod my-container
```

### Retrieve container logs

```bash
# Interactive - follow the prompts
wimkube pod

# Direct command
wimkube pod logs my-pod my-container
```

### Use custom kubeconfig

```bash
wimkube --kubeconfig /path/to/kubeconfig context list
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [huh](https://charm.land/huh/v2) - Interactive forms and prompts
- [client-go](https://github.com/kubernetes/client-go) - Kubernetes API client
- [kubectl](https://github.com/kubernetes/kubectl) - Kubernetes command-line tool library
- [term](https://pkg.go.dev/golang.org/x/term) - Terminal handling

## Development

### Project Structure

```
wimkube/
├── cmd/
│   ├── root.go       # Root command and configuration
│   ├── context.go    # Context management commands
│   ├── namespace.go  # Namespace management commands
│   ├── pod.go        # Pod management commands
│   └── version.go    # Version command
├── internal/
│   ├── client.go     # Kubernetes client wrapper
│   └── kubeconfig.go # Kubeconfig operations
├── main.go           # Entry point
├── go.mod
└── README.md
```

## License

See [LICENSE](LICENSE) file for details.

## Author

[wim-vdw](https://github.com/wim-vdw)
