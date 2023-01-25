<p align="center">
 <img src="https://github.com/datreeio/datree/blob/main/images/datree_GitHub_hero.png" alt="datree=github" border="0" />
</p>
 
<p align="center">
 <img src="https://img.shields.io/travis/com/datreeio/datree/staging?label=build-staging" target="_blank"></a>
 <img src="https://img.shields.io/travis/com/datreeio/datree/main?label=build-main" target="_blank"></a>
 <img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fdatreeio%2Fdatree&count_bg=%2379C83D&title_bg=%23555555&icon=github.svg&icon_color=%23E7E7E7&title=views+%28today+%2F+total%29&edge_flat=false" target="_blank"></a>
 <img src="https://img.shields.io/github/downloads/datreeio/datree/total.svg" target="_blank"></a>
 <img src="https://goreportcard.com/badge/github.com/datreeio/datree" target="_blank"></a>
</p>
  
<p align="center">
  <a href="https://bit.ly/3BHwCEG" target="_blank">
   <img src="https://img.shields.io/badge/Slack-4A154B?logo=slack&color=black&logoColor=white&style=for-the-badge alt="Join our Slack!" width="80" height="30">
  </a> 
</p>

<p align="center">
  <a href="https://hub.datree.io/#utm_source=github&utm_medium=organic_oss"><strong>Explore the docs »</strong></a>
  <br />
</p>

# Datree

[Datree](https://www.datree.io/) (pronounced `/da-tree/`) secures your Kubernetes by blocking the deployment of misconfigured resources.

## ✌️ Quick-start in two steps

Install Datree to get insights on the status of your cluster and enforce your desired policies on new resources.

> **NOTE:**  
> By default, Datree does not block misconfigured resources, it only monitors and alerts about them.  
> To enable **enforcement mode**, see the [documentation](https://hub.datree.io/setup/behavior#options).

### 1. Add the Datree Helm repository
Run the following command in your terminal:
```terminal
helm repo add datree-webhook https://datreeio.github.io/admission-webhook-datree
helm repo update
```

### 2. Install Datree on your cluster
Replace `<DATREE_TOKEN>` with the token from your [dashboard](https://app.datree.io/), and run the following command in your terminal:

```terminal
helm install -n datree datree-webhook datree-webhook/datree-admission-webhook --debug \
--create-namespace \
--set datree.token=<DATREE_TOKEN> \
--set datree.clusterName=$(kubectl config current-context)
```

This will create a new namespace (datree), where Datree’s services and application resources will reside. `datree.token` is used to connect your dashboard to your cluster. Note that the installation can take up to 5 minutes.

## ⚙️ How it works

Datree scans Kubernetes resources against a centrally managed policy, and blocks those that violate your desired policies.

Datree comes with multiple pre-built policies covering various use-cases, such as workload security, high availability, ArgoCD best practices, NSA hardening guide, and [many more](https://hub.datree.io/built-in-rules). 

In addition to our built-in rules, you can write [any custom rule you wish](https://hub.datree.io/custom-rules-overview) and then run it against your Kubernetes configurations to check for rule violations. The custom rule engine is based on JSON Schema.

## 🤩 Additional features

Datree offers a suite of features to make adoption seamless:
* **Monitoring** - Datree is first installed in monitoring mode that reports on policy violations, rather than block their deployments.
* [**CLI**](https://hub.datree.io/cli/getting-started) - Help your developers find misconfigurations in their configs before deploying them, by integrating Datree into their CI.
* **Misconfiguration prioritization** - Datree makes it easy to improve the quality of your cluster by prioritizing the misconfigurations to be fixed.
* **Cluster score** - Rank the stability of your cluster based on the number of detected misconfigurations.

## 📊 Management dashboard (web application)

Datree can be customized via code (policy as code) or via a management dashboard. The dashboard offers the following capabilities in an intuitive visual interface: 
* Customize policies
* Edit rules failure message
* Issue tokens
* View policy check history
* Configure Kubernetes schema version

<img src="/images/dashboard-policies.png" alt="Datree-saas" width="70%">

## Contributing

[Contributions](https://github.com/datreeio/datree/issues?q=is%3Aissue+is%3Aopen+label%3A%22up+for+grabs%22) are welcome!

[![Contributors](https://contrib.rocks/image?repo=datreeio/datree)](https://github.com/datreeio/datree/graphs/contributors)

Thank you to all the people who already contributed to Datree ❤️
