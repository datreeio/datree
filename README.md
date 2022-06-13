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

[Datree](https://datree.io/#utm_source=github&utm_medium=organic_oss) is a CLI tool that supports Kubernetes admins in their roles by preventing developers from making errors in Kubernetes configurations that can cause clusters to fail in production. Our CLI tool is open source, enabling it to be supported by the Kubernetes community.

It’s far more effective than manual processes, such as sending an email to a slew of developers, begging them to set various limits, which likely falls on deaf ears because developers are already overwhelmed.

## ⚙️ How it works

The CLI integration provides a policy enforcement solution for Kubernetes to run automatic checks on every code change for rule violations and misconfigurations. When rule violations are found, Datree produces an alert that guides the developer to fix the issue inside the CI process - or even earlier as a pre-commit hook - while explaining the reason behind the rule.

## ⏩ Quick-start in two steps

### 1. Install the latest release on your CLI

_Linux & MacOS:_ `curl https://get.datree.io | /bin/bash`  
_Windows:_ `iwr -useb https://get.datree.io/windows_install.ps1 | iex`

_Other installation options (Homebrew, Docker, etc.) can be found [here](https://hub.datree.io/#1-install-the-datree-cli)_

### 2. Pass Datree a Kubernetes manifest file to scan

`datree test [k8s-manifest-file]`

...and voilà, you just ran your first policy check! 🥳

## [Command Line Interface](https://hub.datree.io/cli-output) (CLI)

<img src="https://clipublic.s3.amazonaws.com/live.gif" alt="Datree-cli" width="60%" height="50%">

## [Web Application Interface](https://hub.datree.io/centralized-policy) (Dashboard)

<img src="https://user-images.githubusercontent.com/19731161/130956287-ca44e831-46ba-48fa-96eb-be8e23d43bdf.png" alt="Datree-saas" width="60%" height="50%">

<img src="https://user-images.githubusercontent.com/19731161/130957021-4b825b82-01e1-47ba-bf6f-68003f08a532.png" alt="Datree-saas" width="60%" height="50%">

## 🏛️ Architecture

![Architecture](https://github.com/datreeio/datree/blob/main/images/datree_architecture_light.png#gh-light-mode-only)  
![Architecture](https://github.com/datreeio/datree/blob/main/images/datree_architecture_dark.png#gh-dark-mode-only)  

## 🔌 Helm plugin

[Datree's Helm plugin](https://github.com/datreeio/helm-datree) can be accessed through the helm CLI, to provide a seamless experience to Helm users:

`helm plugin install https://github.com/datreeio/helm-datree`

## 🗂 Kustomize support

Datree comes with out-of-the-box [support for Kustomize](https://hub.datree.io/kustomize-support):

`datree kustomize test [kustomization.yaml dir path/]`

## 🤖 Built-in schema validation & policy check

Every policy check will (also) validate your [Kubernetes schema](https://hub.datree.io/schema-validation). In addition, there are 30 battle-tested rules for you to select to create your policy.

The policy rules cover a variety of [Kubernetes resources and use cases](https://hub.datree.io/built-in-rules):

- Workload
- CronJob
- Containers
- Networking
- Security
- Deprecation
- Others

## 🔧 Custom rules

In addition to our built-in rules, you can write [any custom rule](https://hub.datree.io/custom-rules-overview) you wish, and then run them against your Kubernetes configurations to check for rule violations. The custom rule engine is based on JSON Schema.

## 🔗 CI/CD integrations

Like any linter or static code analysis tool, Datree's command-line tool can be **integrated with all CI/CD platforms** to automatically scan every code change and provide feedback as part of the workflow. In the [docs](https://hub.datree.io/cicd-examples), you can find examples of some of the common CI/CD platforms.

If you run into any difficulties with CI/CD integrations, please join our [community Slack channel](https://bit.ly/3BHwCEG) or open an issue, and we'd be happy to guide you through it.

## Contributing

[Contributions](https://github.com/datreeio/datree/issues?q=is%3Aissue+is%3Aopen+label%3A%22up+for+grabs%22) are welcome!

[![Contributors](https://contrib.rocks/image?repo=datreeio/datree)](https://github.com/datreeio/datree/graphs/contributors)

Thank you to all the people who already contributed to Datree ❤️
