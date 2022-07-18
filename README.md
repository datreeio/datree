<p align="center">
 <img src="https://github.com/datreeio/datree/blob/main/images/datree_GitHub_hero.png" alt="datree=github" border="0" />
</p>

<h1 align="center">
 Prevent Kubernetes Misconfigurations
</h1>
 
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

## 🤔 What is Datree?

[Datree](https://datree.io/) automatically validates Kubernetes objects for rule violations, ensuring no misconfigurations reach production. It’s an E2E policy enforcement solution that can be used on the command line, admission webhook, or even as a kubectl plugin.

It’s far more effective than manual processes, such as sending an email to a slew of developers, begging them to set various limits, which likely falls on deaf ears because developers are already overwhelmed.

## ⚙️ How it works

Datree scans Kubernetes configuration changes and validates them against a centrally managed policy for rule violations and misconfigurations.

![Architecture](https://github.com/datreeio/datree/blob/main/images/datree_architecture_light.png#gh-light-mode-only)  
![Architecture](https://github.com/datreeio/datree/blob/main/images/datree_architecture_dark.png#gh-dark-mode-only) 

The CLI interface can be run locally, as a pre-commit hook, or in your CI, to shift left misconfiguration detection. With the admission webhook interface, you can enforce the same policy on the cluster.

### Each Datree scan runs three validation on your Kubernetes objects:
* YAML validation
* Schema validation (Including CRD support)
* Policy check

Datree comes with dozens of battle-tested rules for you to select to create your policy. The policy rules cover a variety of Kubernetes resources such as workload security, networking availability, Argo best practices, NSA hardening guide, and [many more](https://hub.datree.io/built-in-rules). 

In addition to our built-in rules, you can write [any custom rule you wish](https://hub.datree.io/custom-rules-overview) and then run it against your Kubernetes configurations to check for rule violations. The custom rule engine is based on JSON Schema.


## ✌️ Quick-start in two steps

### 1. Install the latest release on your CLI

_Linux & MacOS:_ `curl https://get.datree.io | /bin/bash`  
_Windows:_ `iwr -useb https://get.datree.io/windows_install.ps1 | iex`

_Other installation options (Homebrew, Docker, etc.) can be found [here](https://hub.datree.io/#1-install-the-datree-cli)_

### 2. Pass Datree a Kubernetes manifest file to scan

`datree test [k8s-manifest-file]`

...and voilà, you just ran your first policy check! 🥳

<img src="https://clipublic.s3.amazonaws.com/live.gif" alt="Datree-cli" width="70%">

## 🔌 Helm plugin

[Datree's Helm plugin](https://github.com/datreeio/helm-datree) can be accessed through the helm CLI to provide a seamless experience to Helm users:

`helm plugin install https://github.com/datreeio/helm-datree`

## 🗂 Kustomize support

Datree comes with out-of-the-box [support for Kustomize](https://hub.datree.io/kustomize-support):

`datree kustomize test [kustomization.yaml dir path/]`

## Management dashboard (web application)

Datree can be customized via code (policy as code) or via a management dashboard. The dashboard offers the following capabilities in an intuitive visual interface: 
* Customize policies
* Edit rules failure message
* Issue tokens
* View policy check history
* Configure Kubernetes schema version

<img src="https://user-images.githubusercontent.com/19731161/130956287-ca44e831-46ba-48fa-96eb-be8e23d43bdf.png" alt="Datree-saas" width="55%">

## ✔️ Next step: Integrate into your workflow

From develop to runtime, you can use Datree in every step of your Kuberenetes pipeline to help you prevent misconfigurations:  
* [Develop (code)](https://hub.datree.io/#2-test-a-kubernetes-demo-manifest) - run the CLI locally (or as a pre-commit hook) to get instant validation  
* [Distribute (CI)](https://hub.datree.io/cicd-examples) - integrate with your CI platform to shift-left policy checks  
* [Deploy (CD)](https://github.com/datreeio/admission-webhook-datree) - gate your cluster with the admission webhook  
* [Runtime (production)](https://github.com/datreeio/kubectl-datree) - query deployed resources with the kubectl plugin to your know your status  

## Contributing

[Contributions](https://github.com/datreeio/datree/issues?q=is%3Aissue+is%3Aopen+label%3A%22up+for+grabs%22) are welcome!

[![Contributors](https://contrib.rocks/image?repo=datreeio/datree)](https://github.com/datreeio/datree/graphs/contributors)

Thank you to all the people who already contributed to Datree ❤️
