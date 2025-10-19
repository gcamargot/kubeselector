# kubeselector

An interactive command-line tool for Kubernetes that lets you easily browse and select resources like Pods, Secrets, and ConfigMaps then perform actions such as delete, describe, or view logs.

## Features

- **Interactive CLI tool** built in Go that simplifies Kubernetes resource management
- **Quickly find** Pods, Secrets, ConfigMaps by name or label
- **Automatically lists** resources in the current namespace, or switch with a flag
- **One-key actions**: delete, describe, or fetch logs directly
- Uses your current **kubeconfig, context, and permissions**

## Installation

```bash
go install github.com/gcamargot/kubeselector@latest
```

Or build from source:

```bash
git clone https://github.com/gcamargot/kubeselector.git
cd kubeselector
go build -o kubeselector .
```

## Usage

### Basic Commands

List and interact with pods:
```bash
kubeselector pods
```

List and interact with secrets:
```bash
kubeselector secrets
```

List and interact with configmaps:
```bash
kubeselector configmaps
```

### Command Options

#### Delete Resources
Delete a selected pod, secret, or configmap:
```bash
kubeselector pods --delete
kubeselector secrets --delete
kubeselector configmaps --delete
```

#### Describe Resources
Get detailed information about a resource:
```bash
kubeselector pods --describe
kubeselector secrets --describe
kubeselector configmaps --describe
```

#### View Pod Logs
Fetch logs from a selected pod:
```bash
kubeselector pods --logs
```

#### Filter Resources
Filter resources by name:
```bash
kubeselector secrets --filter db
kubeselector pods --filter nginx
```

#### Specify Namespace
Use a specific namespace instead of the current context:
```bash
kubeselector pods --namespace kube-system
kubeselector secrets -n production
```

### Example Workflows

1. **Delete a specific pod in production namespace:**
   ```bash
   kubeselector pods --namespace production --delete
   ```

2. **View logs from a pod with 'api' in its name:**
   ```bash
   kubeselector pods --filter api --logs
   ```

3. **Describe a secret containing 'db':**
   ```bash
   kubeselector secrets --filter db --describe
   ```

## Interactive Selection

When you run a command, kubeselector will present an interactive list where you can:
- Use **arrow keys** to navigate
- Press **/** to filter/search
- Press **Enter** to select
- Press **Ctrl+C** to cancel

## Requirements

- Go 1.21 or higher
- Access to a Kubernetes cluster
- Valid kubeconfig file

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
