# Datree Admission Webhook

<p align="center">
<img src="https://github.com/datreeio/admission-webhook-datree/blob/main/internal/images/diagram.png" width="80%" />
</p>
  
# Overview
Datree offers cluster integration that allows you to validate your resources against your configured policy upon pushing them into a cluster, by using an admission webhook.

The webhook will catch **create**, **apply** and **edit** operations and initiate a policy check against the configs associated with each operation. If any misconfigurations are found, the webhook will reject the operation, and display a detailed output with instructions on how to resolve each misconfiguration.

## Webhook validation triggers

K8s use different abstractions to simplify and automate complex processes. For example, when explicitly applying an object type “Deployment”, under the hood, K8s will “translate” this object into implicit objects of type “Pod.”

When installed on your cluster, other policy enforcement tools will validate both explicit and implicit objects. This approach may create a lot of noise and false positive failures since it will cause the webhook to validate objects that the users don’t manage and, in some cases, are not even accessible.

To avoid such issues, we decided to define the specific operations that the admission webhook should validate:

- Kubectl - validate objects that were created or updated using kubectl `create`, `edit`, and `apply` commands. Objects that were implicitly created (e.g., pods created via deployment) are ignored since the webhook validates the deployment that generated them and is accessible to the user.
- Gitops CD tools - validate objects that were explicitly created and distinguish them from other objects (custom resources) that were implicitly created during the installation and are required for the ongoing operation of these tools (e.g., ArgoCD, FluxCD, etc.)

## Prerequisites

The webhook officially supports **Kubernetes version _1.19_** and higher, and has been tested with EKS.

# Installation

## Deploy with Helm

```bash
  # Install and create namespace with Helm
  helm repo add datree-webhook https://datreeio.github.io/admission-webhook-datree/
  helm repo update

  # Already existing `datree` namespace
  kubectl create ns datree
  helm install -n datree datree-webhook datree-webhook/datree-admission-webhook --set datree.token=<DATREE_TOKEN>
```
 
For more information see [Datree webhook Helm chart](https://github.com/datreeio/admission-webhook-datree/tree/gh-pages).

## Deploy with installation script

During the installtion the script will require to enter the Datree token during installation.

```bash
# Install with prompting Datree token
bash <(curl https://get.datree.io/admission-webhook)

# Install without prompting Datree token
DATREE_TOKEN=[your-token] bash <(curl https://get.datree.io/admission-webhook)

```

### Prerequisites

The following applications must be installed on the machine:

- kubectl
- openssl - _required for creating a certificate authority (CA)._
- curl

## Usage

Once the webhook is installed, every hooked operation will trigger a Datree policy check. If no misconfigurations were found, the resource will be applied/updated.
For any misconfigurations that were found the following output will be displayed:

![deny-example](https://raw.githubusercontent.com/datreeio/admission-webhook-datree/main/internal/images/deny-example.png)

## Behavior

### Token

🤫 Since your token is sensitive and you would not want to keep it in your repository, we recommend to set/change it by running a separate `kubectl patch` command:

```yaml
kubectl patch deployment webhook-server -n datree -p '
spec:
  template:
    spec:
      containers:
        - name: server
          env:
            - name: DATREE_TOKEN
              value: "<your-token>"'
```

Simply replace `<your-token>` with your actual token, then copy the entire command and run it in your terminal.

### Other settings

1. Create a YAML file in your repository with this content:

```yaml
spec:
  template:
    spec:
      containers:
        - name: server
          env:
            - name: DATREE_POLICY
              value: ""
            - name: DATREE_VERBOSE
              value: ""
            - name: DATREE_OUTPUT
              value: ""
            - name: DATREE_NO_RECORD
              value: ""
```

2. Change the values of your settings as you desire.
3. Run the following command to apply your changes to the webhook resource:

```bash
kubectl patch deployment webhook-server -n datree --patch-file /path/to/patch/file.yaml
```

## Ignore a namespace

Add the label `"admission.datree/validate=skip"` to the configuration of the namespace you would like to ignore:

```bash
kubectl label namespaces default "admission.datree/validate=skip"
```

To delete the label and resume running the datree webhook on the namespace again:

```bash
kubectl label namespaces default "admission.datree/validate-"
```

## Uninstallation

To uninstall the webhook, copy the following command and run it in your terminal:

```bash
bash <(curl https://get.datree.io/admission-webhook-uninstall)
```

To uninstall the helm release, copy the following command and run it in your terminal:

```bash
helm uninstall datree-webhook -n datree
kubectl delete ns datree
```

## Local development

To run the webhook locally (in development), view our [developer guide](/guides/developer-guide.md).
