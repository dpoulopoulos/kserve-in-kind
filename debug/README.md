# Debugging in Kubernetes

This directory provides detailed instructions on debugging a running process within a Kubernetes
Pod. The steps outlined will guide you through the process of constructing a debug pod, injecting it
as an ephemeral container into the target pod, and subsequently initiating a Go debugging session
using Delve. You can use this knowledge to debug your KServe Inference Services, or the any KServe
component you need, such as the KServe or Knative controllers.

# What You'll Need

To complete this guide, you'll need a working Docker installation, a kind cluster, and `kubectl`.
You can follow the [KServe in kind](../README.md) guide to set up the kind cluster.

# Procedure

Follow the steps below to start your debugging session:

1. Build your debug container image:

   ```shell
   docker build -t debian:debug .
   ```

1. Push the pod on kind:

   ```shell
   kind load docker-image debian:debug
   ```

1. Build the client and server images you'll be using a debug targets:

    ```shell
    docker build -t client:distroless -f client/Dockerfile client/
    docker build -t server:distroless -f server/Dockerfile server/
    ```

1. Push the client and server images on kind:

    ```shell
    kind load docker-image client:distroless
    kind load docker-image server:distroless
    ```

1. Inject the debug container into the target pod. To run a live debugging session in a pod, your
   ephemeral container should have the `SYS_PTRACE` capability and run as root. Unfortunately, you
   cannot set this through the `kubectl debug` tool yet. Instead, you need to inject this by sending
   a POST request directly to the API server:

    a. Start a proxy server using `kubectl` to access the API server:

        ```shell
        kubectl proxy &
        ```
    
    c. Create a POST request to inject the ephemeral container:

        ```shell
        curl localhost:8001/api/v1/namespaces/default/pods/client/ephemeralcontainers \
        -XPATCH -H 'Content-Type: application/strategic-merge-patch+json' \
        -d '
        {
            "spec":
            {
                "ephemeralContainers":
                [
                    {
                        "name": "debugger",
                        "command": ["/bin/bash"],
                        "image": "debian:debug",
                        "targetContainerName": "server",
                        "stdin": true,
                        "tty": true,
                        "securityContext": {
                            "capabilities": {
                                "add": ["SYS_PTRACE"]
                            },
                            "runAsNonRoot": false,
                            "runAsUser": 0,
                            "runAsGroup": 0
                        }
                    }
                ]
            }
        }'
        ```

    d. Repeat the process for the server pod.

1. Start the Delve debug server and attach on PID `1`:

    ```shell
    kubectl exec client --container=debugger -- \
    dlv --headless --accept-multiclient --api-version=2 \
        --listen  127.0.0.1:1111 \
        --continue  --log  --log-dest=dlv.log  \
        attach 1
    ```

    > Repeat the process for the server Pod but note to change the port that the server is listening
    > to. The debuging configuration uses `1112` by default.

1. Set up port forwarding to reach the two debugging servers running on the two Pods:

    ```shell
    kubectl port-forward po/client 1111:1111 &
    kubectl port-forward po/server 1112:1112 &
    ```

1. Navigate to VS Code's debugging UI and start a new debugging session. You can use the
   configuration provided under the `vscode` directory as a starting point. Set your breakpoints and
   start debugging:

    ```
    kubectl port-forward po/client 6000:6000 &
    curl localhost:6000
    ```

# Conclusion

You've successfully set up a debugging session for a running process in Kubernetes. You can use this
knowledge to debug your KServe Inference Services, or the any KServe component you need. Moreover,
the debug container you built includes a number of useful tools that you can use to debug your Pods
in a live environment, even if they do not have a shell.