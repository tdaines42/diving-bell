# diving-bell

## Config

The config file `cluster-config.yaml` is currently expected to be in the folder you run `diving-bell` from.

Example config for a cluster:

    clusterName: "test-cluster"
    controlPlaneTarget: "10.17.1.0"
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
