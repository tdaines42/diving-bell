# diving-bell

## Usage

    Manage a k8s cluster using kubectl, terraform, and skuba

    Usage:
      diving-bell [command]

    Available Commands:
      bootstrap   Bootstrap the cluster using skuba
      config      Print the config to the console
      deploy      Deploy a cluster with one command
      help        Help about any command
      provision   Provision the cluster using terraform
      status      Get the current status of the cluster

    Flags:
          --config string   config file (default is $HOME/.diving-bell.yaml)
      -h, --help            help for diving-bell

    Use "diving-bell [command] --help" for more information about a command.`

## Config

Example config for a cluster at `$HOME/.diving-bell.yaml`:

    clusterName: "test-cluster"
    controlPlaneTarget: "10.17.1.0"
    terraformWorkspacePath: "~/github/skuba/ci/infra/libvirt"
    managers:
            - user: "sles"
              target: "10.17.2.0"
              hostName: "testing-master-0"
            - user: "sles"
              target: "10.17.2.1"
              hostName: "testing-master-1"
            - user: "sles"
              target: "10.17.2.2"
              hostName: "testing-master-2"

    workers:
            - user: "sles"
              target: "10.17.3.0"
              hostName: "testing-worker-0"

            - user: "sles"
              target: "10.17.3.1"
              hostName: "testing-worker-1"
