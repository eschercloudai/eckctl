# eckctl

## Introduction
`eckctl` is a simple utility for working with the EscherCloud Kubernetes (ECK) API, and which provides a marginal UX improvement over a bunch of curl commands wrapped in a shell script.

## Installation

Grab a binary from the [Releases](https://github.com/eschercloudai/eckctl/releases) page.

It probably helps to familiarise yourself with the [ECK architecture and principles](https://docs.eschercloud.ai/Kubernetes/Reference/overview).

## Configuration

`eckctl` can be configured via `~/.eckctl.yaml`:

```yaml
username: "n.jones@eschercloud.ai"
password: "hunter2"
url: "https://unikorn.nl1.eschercloud.dev"
project: "abc123"
```

These options can also be set (or overriden) via the command line.

## Usage

### Get projects

```shell
% eckctl get projects
Name: njones-demo	ID: abc123
Name: hotdog	Description: Management cluster ID: abc123
```

### List images

Get a list of compatible images:

```shell
% eckctl get images
Name: eck-230408-ac77c987	UUID: ab59076c-142a-4da6-9149-e8426fb93aff	Created: 2023-04-08 09:34:16 +0000 UTC	Kubernetes version: v1.25.6	NVIDIA driver version: 525.85.05
Name: eck-230414-cc034e2d	UUID: e1ad0805-ecd9-4162-b243-e2bfd85e8aed	Created: 2023-04-14 12:47:16 +0000 UTC	Kubernetes version: v1.26.3	NVIDIA driver version: no_gpu
```

### List networks

Get external networks:

```shell
% eckctl get networks
Name: Internet	ID: c9d130bc-301d-45c0-9328-a6964af65579
```

### List ECK control planes:

```shell
% eckctl get controlplanes
Name: dibnah	Status: Provisioned	Version: 1.0.1
```

### List clusters

```shell
% eckctl get clusters --controlplane default
Name: demo	Version: v1.26.1	Status: Provisioned
```

### Get versions

```shell
% eckctl get versions
Cluster Bundles:
Name: kubernetes-cluster-1.0.0	Version: 1.0.0
Name: kubernetes-cluster-1.1.0	Version: 1.1.0
Name: kubernetes-cluster-1.2.0	Version: 1.2.0
Control Plane Bundles:
Name: control-plane-1.0.0	Version: 1.0.0
Name: control-plane-1.0.1	Version: 1.0.1
```

### Create a control plane

```shell
% eckctl create controlplane --name default --version 1.0.1
```

### Create a cluster

Creating a cluster requires a number of parameters to be defined, including workload pools plus associated flavours and whether to enable autoscaling.  For now this definition needs to be crafted as a blob of JSON, see the example in the [examples](https://github.com/eschercloudai/eckctl/tree/main/examples) folder.  A cluster can be created once the associated ECK control plane is in status `Provisioned`:

```shell
% eckctl get controlplane --name default
Name: default	Status: Provisioned	Version: 1.0.1
% eckctl create cluster --name demo \
  --controlplane default --json ./examples/cluster.json
% eckctl get cluster --name demo --controlplane default
Cluster: demo, version: v1.26.1, status: Provisioning.
└── Pools:
    └── Name: worker	Flavor: g.2.standard	Image: eck-230408-3464bc03
```

### Retrieve kubeconfig

Once the cluster's status is `Provisioned` you can retrieve its kubeconfig:

```shell
% eckctl get clusters --controlplane default --name demo
Cluster: demo, version: v1.26.1, status: Provisioned.
└── Pools:
    └── Name: worker	Flavor: g.2.standard	Image: eck-230408-3464bc03
% eckctl get kubeconfig --controlplane default --cluster demo > ~/kubeconfig-demo
% export KUBECONFIG=~/kubeconfig-demo
% kubectl version -o json | jq .serverVersion
{
  "major": "1",
  "minor": "26",
  "gitVersion": "v1.26.1",
  "gitCommit": "8f94681cd294aa8cfd3407b8191f6c70214973a4",
  "gitTreeState": "clean",
  "buildDate": "2023-01-18T15:51:25Z",
  "goVersion": "go1.19.5",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```
### Testing Using Docker Image
Use the below snippet to create the image and then use the commands above to test things work as expected.

```shell
docker build -t eckctl:v0.0.0 -f docker/Dockerfile .
docker run --name eckctl -it --rm -v /home/$USER/.eckctl.yaml:/home/eckctl/.eckctl.yaml eckctl:v0.0.0 "eckctl get images"
```
