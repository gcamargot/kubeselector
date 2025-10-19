kubeselector

kubeselector is an interactive CLI tool built in Go that simplifies Kubernetes resource management.
It provides a fast and user-friendly way to search, filter, and act on cluster resources — directly from your terminal.

Features

Search across resources: quickly find Pods, Secrets, ConfigMaps, and more by name or label.

Namespace-aware: automatically lists resources in the current namespace, or switch with a flag.

Interactive selection: choose one or multiple items from a searchable list (powered by fuzzy finder).

One-key actions: delete, describe, or fetch logs directly.

kubectl-compatible: uses your current kubeconfig, context, and permissions.

Extensible: modular design ready for additional resource types or custom actions.

Lightweight binary: written in Go — no dependencies, no Python, no Node.

Example Usage
# Select and delete a pod interactively
kubeselector pods --delete

# Search for secrets containing “db”
kubeselector secrets --filter db

# View logs of a selected pod
kubeselector pods --logs

Architecture

kubeselector is built on top of:

client-go
 – for Kubernetes API access.

cobra
 – for CLI structure and subcommands.

survey
 or bubbletea
 – for interactive UI/selection.

Roadmap

 Support for more resource types (Deployments, ConfigMaps, etc.)

 Multi-select delete mode.

 Integration with kubectx for context switching.

 “Safe mode” confirmation before destructive actions.

 Export results in YAML/JSON.

Installation
go install github.com/<yourusername>/kubeselector@latest


Or clone and build manually:

git clone https://github.com/<yourusername>/kubeselector.git
cd kubeselector
make build

Gaston Camargo
