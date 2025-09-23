# mcc

[![Go Reference](https://pkg.go.dev/badge/github.com/onggizam/mcc.svg)](https://pkg.go.dev/github.com/onggizam/mcc)
[![Release](https://img.shields.io/github/v/release/onggizam/mcc)](https://github.com/onggizam/mcc/releases)
[![Homebrew](https://img.shields.io/badge/homebrew-available-blue)](https://github.com/onggizam/homebrew-mcc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> Manage & switch multiple Kubernetes kubeconfigs with a single command

`mcc` is a lightweight CLI tool to manage and switch between multiple **Kubernetes cluster configs**.  
It stores kubeconfig files in `~/.kube/multi` and switches them into `~/.kube/config` when needed.

## Features

- 📂 **Organize kubeconfigs** — keep multiple cluster configs under `~/.kube/multi`
- 🔄 **One-command switch** — instantly swap the active cluster
- 🗑 **Clean up easily** — remove configs you no longer need
- 📝 **Readable list** — view stored clusters, created date, server info, and active flag

## Installation

### Homebrew (recommended)

```bash
brew tap onggizam/mcc
brew install mcc
```

### Build from source

```bash
git clone https://github.com/onggizam/mcc.git
cd mcc
bash scripts/build.sh
```

## Usage

First, export your cluster kubeconfig:

```bash
cat ~/.kube/config

# or generate from cluster
kubectl get cm kubeadm-config -n kube-system -o yaml > cluster1
```

## Commands & Options

| Command         | Options             | Description                                                                  |
| --------------- | ------------------- | ---------------------------------------------------------------------------- |
| `mcc add`       | `-f, --file <path>` | Source kubeconfig file (default: current `~/.kube/config`)                   |
|                 | `-n, --name <name>` | Name to store the kubeconfig as (**required**)                               |
|                 | `--force`           | Overwrite if the name already exists                                         |
| `mcc ch <name>` | `--backup`          | Switch to the given cluster config; optionally backup the current config     |
| `mcc list`      | _(none)_            | Show stored kubeconfigs in a table (NO., NAME, CREATED, SERVER, ACTIVE)      |
| `mcc delete`    | `<name>`            | Remove the stored kubeconfig with the given name                             |
| `mcc now`       | _(none)_            | Show current active cluster info + cluster-wide Pod summary (all namespaces) |
| `mcc version`   | _(none)_            | Print the current version of mcc                                             |

### Examples

```bash
mcc add -f ./myconfig -n cluster1 --force   # Add (overwrite if exists)
mcc ch cluster1 --backup                    # Switch and backup current config
mcc list                                    # Show all stored configs
mcc delete cluster1                         # Delete stored config
mcc now                                     # Show current cluster info and Pods summary
mcc version                                 # Show current version
```

Check the version:

```bash
mcc version
```

## 🤝 Contributing

Contributions are welcome!
