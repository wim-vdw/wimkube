# wimkube

An interactive Kubernetes CLI tool built with Go that simplifies common Kubernetes operations through an intuitive
command-line interface.

## Features

- **Context Management**: Switch between and manage Kubernetes contexts
- **Namespace Management**: View and switch between namespaces
- **Pod Operations**: List pods, view containers, and execute interactive shells
- **Interactive Menus**: User-friendly interactive prompts for all operations
- **Direct Commands**: Support for both interactive and direct command execution

## Installation

### Prerequisites

- Go 1.26 or later
- Access to a Kubernetes cluster
- kubectl configured with valid kubeconfig

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

## Examples

### Switch to a different context

```bash
# Interactive
wimkube context

# Direct command
wimkube context set production-cluster
```

### Change namespace

```bash
# Interactive
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

### Use custom kubeconfig

```bash
wimkube --kubeconfig /path/to/kubeconfig context list
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [huh](https://github.com/charmbracelet/huh) - Interactive forms and prompts
- [client-go](https://github.com/kubernetes/client-go) - Kubernetes API client
- [kubectl](https://github.com/kubernetes/kubectl) - Kubernetes command-line tool library

## Development

### Project Structure

```
wimkube/
├── cmd/
│   ├── root.go       # Root command and configuration
│   ├── context.go    # Context management commands
│   ├── namespace.go  # Namespace management commands
│   └── pod.go        # Pod management commands
├── internal/
│   ├── client.go     # Kubernetes client wrapper
│   └── kubeconfig.go # Kubeconfig operations
├── main.go           # Entry point
├── go.mod
└── README.md
```

### Building

```bash
go build -o wimkube
```

## License

See [LICENSE](LICENSE) file for details.

## Author

[wim-vdw](https://github.com/wim-vdw)
