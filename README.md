[![datree-logo](https://raw.githubusercontent.com/datreeio/datree/main/images/datree_LOGO-180px.png)](#) 

![Travis (.com) branch](https://img.shields.io/travis/com/datreeio/datree/staging?label=build-staging)
![Travis (.com) branch](https://img.shields.io/travis/com/datreeio/datree/main?label=build-main)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fdatreeio%2Fdatree&count_bg=%2379C83D&title_bg=%23555555&icon=github.svg&icon_color=%23E7E7E7&title=views+%28today+%2F+total%29&edge_flat=false)](https://hits.seeyoufarm.com)
![open issues](https://img.shields.io/github/issues-raw/datreeio/datree)

## What is Datree?
[Datree](https://datree.io/?utm_source=github&utm_medium=organic_oss) is a CLI tool that supports Kubernetes admins in their roles, by preventing developers from making errors in Kubernetes configurations that can cause clusters to fail in production. Our CLI tool is open source, enabling it to be supported by the Kubernetes community.  

It’s far more effective than manual processes, such as sending an email to a slew of developers, begging them to set various limits, which likely falls on deaf ears because developers are already overwhelmed.

## How it Works
The CLI integration provides a policy enforcement solution for Kubernetes to run automatic checks on every code change for rule violations and misconfigurations. When rule violations are found, Datree produces an alert which guides the developer to fix the issue inside the CI process — or even earlier as a pre-commit hook — while explaining the reason behind the rule.

## Quick start in two steps
#### 1. Install the latest release on your CLI  
**Linux & MacOS:** `curl https://get.datree.io | /bin/bash`  
**Windows:** `iwr -useb https://get.datree.io/windows_install.ps1 | iex`  

#### 2. Pass datree a Kuberntes manifest file to scan
`datree test <k8s-manifest-file>`  

...and voilà, you just ran your first invocation! 🥳    

In your command line, you’ll see something like this:
<img src="https://clipublic.s3.amazonaws.com/live.gif" alt="Datree" width="600" height="400">  

#### Playground
[![katacoda-logo](https://raw.githubusercontent.com/datreeio/datree/main/images/katacoda-logo.png)](#)  
You can also checkout our [interactive demo scenario](https://www.katacoda.com/datree/scenarios/datree-demo) on Katacoda without having to install anything on your machine.  

#### Ready to review our "Getting Started" guide?
All the information needed to get started, as well as a bunch of other cool features (including how to set up your policy), can be found in [**our docs**](https://hub.datree.io/getting-started/?utm_source=github&utm_medium=organic_oss).

## Helm plugin
[Datree's Helm plugin](https://hub.datree.io/helm-plugin/?utm_source=github&utm_medium=organic_oss) can be accessed through the helm CLI, to provide a seamless experience to Helm users:  

`helm plugin install https://github.com/datreeio/helm-datree`  

## Built-in rules
Right now, there are 30 battle-tested rules for you to choose from.    
The rules cover different Kubernetes resources and use cases:
* [Workload](https://hub.datree.io/workload/?utm_source=github&utm_medium=organic_oss)
* [CronJob](https://hub.datree.io/cronjob/?utm_source=github&utm_medium=organic_oss)
* [Containers](https://hub.datree.io/containers/?utm_source=github&utm_medium=organic_oss)
* [Networking](https://hub.datree.io/networking/?utm_source=github&utm_medium=organic_oss)
* [Deprecation](https://hub.datree.io/deprecation/?utm_source=github&utm_medium=organic_oss)
* [Others](https://hub.datree.io/other/?utm_source=github&utm_medium=organic_oss)

## WARNING: Alpha release

Datree is at an early stage of development. We do our best to maintain backwards compatibility but there may be breaking changes in
the future to the command usage, flags, and configuration file formats. The CLI will output a warning message when a new version with breaking changes is detected.
We encourage you to use Datree to test your Kubernetes manifests files and Helm charts, see what
breaks, and [contribute](./CONTRIBUTING.md).

## Support

[Datree](https://datree.io/?utm_source=github&utm_medium=organic_oss) builds and maintains this project to make Kubernetes policies simple and accessible. Start with our [documentations](https://hub.datree.io/?utm_source=github&utm_medium=organic_oss) for quick tutorials and examples.

