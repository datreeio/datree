<p align="center">
 <img src="https://github.com/datreeio/datree/blob/main/images/datree_GitHub_hero.png" alt="datree=github" border="0" />
</p>
 
<p align="center">
 <img src="https://img.shields.io/github/v/release/datreeio/datree" />
 <img src="https://github.com/datreeio/datree/actions/workflows/release.yml/badge.svg" />
 <img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fdatreeio%2Fdatree&count_bg=%2379C83D&title_bg=%23555555&icon=github.svg&icon_color=%23E7E7E7&title=views+%28today+%2F+total%29&edge_flat=false" target="_blank"></a>
 <img src="https://img.shields.io/github/downloads/datreeio/datree/total.svg" target="_blank"></a>
 <img src="https://goreportcard.com/badge/github.com/datreeio/datree" target="_blank"></a>
</p>

<p align="center">
  <a href="https://hub.datree.io/#utm_source=github&utm_medium=organic_oss"><strong>Explore the docs »</strong></a>
  <br />
</p>

# Datree [DEPRECATED]

[Datree](https://www.datree.io/) (pronounced `/da-tree/`) was built to secure Kubernetes workloads by blocking the deployment of misconfigured resources. **Since July 2023, the commercial company that supports and actively maintains this project has been closed.**

## Migrating to the (fully) open-source version of Datree 

For existing users, it is still possible to run Datree as a standalone: https://hub.datree.io/cli/offline-mode

## What will not be available anymore

All the archived open source repositories under datreeio org will no longer be maintained and accept any new code changes, including any security patches.
In addition, the following key capabilities will not longer be available anymore:  
* Centralized policy registry
* Automatic Kubernetes schema validation
* Access to the dashboard and all of its components (e.g. activity-log page, token management, etc.)

## ⚙️ How it works

Datree scans Kubernetes resources against a centrally managed policy, and blocks those that violate your desired policies.

Datree comes with over 100 rules covering various use-cases, such as workload security, high availability, ArgoCD best practices, NSA hardening guide, and [many more](https://hub.datree.io/built-in-rules). 

In addition to our built-in rules, you can write [any custom rule you wish](https://hub.datree.io/custom-rules-overview) and then run it against your Kubernetes configurations to check for rule violations. Custom rules can be written in [JSON schema](https://hub.datree.io/custom-rules/custom-rules-overview) or in [Rego](https://hub.datree.io/custom-rules/rego-support).

## Contributing

We want to thank our contributors for helping us build Datree ❤️
  
[![Contributors](https://contrib.rocks/image?repo=datreeio/datree)](https://github.com/datreeio/datree/graphs/contributors)
