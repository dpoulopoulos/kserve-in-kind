# KServe in kind

Welcome to the KServe in kind repository! This project is designed as a one-stop resource for those looking to set up
a new kind (Kubernetes in Docker) cluster and install KServe on it. If you're new to Kubernetes or KServe, or if you're
simply looking to spin up a new environment, you've come to the right place.

![kserve-in-kind](images/kind-kserve.png)

## What is kind?

[kind](https://kind.sigs.k8s.io/) stands for "Kubernetes in Docker" and is a tool for running local Kubernetes clusters
using Docker container "nodes". It is primarily designed as a Kubernetes testbed.

## What is KServe?

[KServe](https://kserve.github.io/) provides a Kubernetes Custom Resource Definition for serving machine learning models
on arbitrary frameworks. It aims to solve production model serving use cases by providing performant, high abstraction
interfaces for common ML frameworks like TensorFlow, PyTorch, and XGBoost.

## Why KServe in kind?

Setting up a Kubernetes cluster and then installing and configuring KServe can be a complex and time-consuming task,
particularly for those who are new to the ecosystem. This repository aims to simplify this process by providing all the
necessary scripts, configuration files, and documentation to get you up and running quickly. This enables you to focus
more on your ML models and less on the infrastructure.

This repository uses kind version v0.20.0. Additionally, it installs and configures:

- Kubernetes `v1.27.3`
- Knative `v1.11.0`
- Istio `v1.18.0`
- Cert Manager `v1.13.0`
- KServe `0.11.0`
- MetalLB `v0.13.7`

# What You'll Need

To complete the installation of KServe on a kind cluster you'll need kind installed, as well as `kubectl`. To avoid
any issues install a `kubectl` version that is compatible with the Kubernetes version you'll install with kind.

# Procedure

Follow the steps below to install KServe on a kind cluster:

1. Create a new kind cluster:

   ```shell
   kind create cluster --config kind/config.yaml
   ```

1. Install Istio:

   ```shell
   kubectl apply -k manifests/istio/overlays/kind
   ```

1. Install knative-serving:
   ```shell
   kubectl apply -k manifests/knative-serving/overlays/kind
   ```

    > This command attempts to deploy Knative Serving; however, some components may fail to install initially due to missing CRDs (Custom Resource Definitions). Running the command a second time should resolve the issue. This is a common procedure that you can follow when you install other components of the cluster, like Istio and KServe.

1. Install cert-manager:

    ```shell
    kubectl apply -k manifests/cert-manager/base/
    ```

1. Install KServe:

   ```shell
   kubectl apply -k manifests/kserve/base/
   ```

1. Install and configure MetalLB:

    a. Define the IP address range you want to use:

    ```shell
    export IP_ADDRESS_RANGE=$(
        docker network inspect -f '{{.IPAM.Config}}' kind \
        | tr -d '[]{}' \
        | awk -F' ' '{print $1}' \
        | sed -E 's/([0-9]+\.[0-9]+)\.[0-9]+(\.[0-9]+\/[0-9]+)/\1.100\2/' \
        | sed 's/\/[0-9]\+/\/24/'
    )
    ```
    
    b. Create the `metallb-config.yaml` file:

    ```shell
    j2 manifests/metallb/base/upstream/metallb-config.yaml.j2 > manifests/metallb/base/upstream/metallb-config.yaml
    ```

    c. Install MetalLB:

    ```shell
    kubectl apply -k manifests/metallb/base
    ```

    > Parts of this installation will fail because certain pods are not yet running. To wait until all pods are ready, execute the following command:
    > ```shell
    > kubectl wait --namespace metallb-system  --for=condition=ready pod  --selector=app=metallb  --timeout=90s
    > ```
    > When every pod is ready run the MetalLB installation command again.

# Verify

Run the following guide to verify that your MetalLB and KServe installations work as expected.

## MetalLB

To verify that MetalLB works as expected, follow the steps below:

1. Apply the MetalLB example:

   ```shell
   kubectl apply -f examples/metallb/foo-bar.yaml
   ```

1. Wait until the pods are ready:

   ```shell
   kubectl wait --namespace default \
                --for=condition=ready pod \
                --selector=app=http-echo \
                --timeout=90s
   ```

1. Get the Load Balancer IP:

   ```shell
   LB_IP=$(kubectl get svc/foo-service -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
   ```

1. Invoke the service:

   ```shell
   for _ in {1..10}; do
     curl ${LB_IP}:5678
   done
   ```

   The output of this command should print `foo` and `bar` on separate lines.

## KServe

To verify that KServe works as expected, follow the steps below:

1. `cd` into the example directory:

   ```shell
   cd examples/kserve/sklearn
   ```

1. Apply the KServe example:

   ```shell
   kubectl apply -f iris.yaml
   ```

1. Wait until the predictor pod is ready:

   ```shell
   kubectl wait --namespace default \
                --for=condition=ready pod \
                --selector=component=predictor \
                --timeout=90s
   ```

1. Get the ingress host:

   ```shell
   export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   ```

1. Get the ingress port:

   ```shell
   export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
   ```

1. Run the Python prediction script:

   ```
   python3 predict.py
   ```

   The output of the script should be:

   ```shell
   Status code: 200
   Predictions: [1, 1]
   ```
