<p align="center">
<img src="https://challengepost-s3-challengepost.netdna-ssl.com/photos/production/software_thumbnail_photos/001/770/244/datas/medium.png" width="300px"/>
<p>

# Datree Automated : Datree Policy Generator

## Important Links

- [Devpost Submission](https://devpost.com/software/shhhh) - Please take a look here for a better explanation of the project
- [GitHub Repository](https://github.com/sauravhiremath/policy-fy) - Clone to try it out for yourself!
- [Demo Video](https://youtu.be/o3AQUmHE-Ms) - Walkthrough of our project

## What it does

Building new policies for each use case seems redundant and this tool automates the process.
The CLI tool takes in a Kubernetes config file, parses it using our custom Datree policy generative algorithm to produce a policy.yml file that can be published to test configurations using Datree.

**Features**

- Parsing YAML config properties
- Support for Resource Limits. Ex: maximum: 25
- Supports enums, string and limit values

## How we built it

The tool is written in TypeScript. Node and Commander are used to run it in CLI.

## Challenges we ran into

We ran into a bunch of challenges while building this, mostly while trying to figure out where we can find the Kubernetes schemas to use ([the swagger file](https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/swagger.json) on the K8s repository is very difficult to parse).

We tried using [Kubeconform](https://github.com/yannh/kubeconform) but ultimately decided to use [Kubernetes JSON Schema](https://github.com/instrumenta/kubernetes-json-schema) for our use case.

Because of the decision we made to pivot on the last day of the hackathon, our biggest constraint became time. However, we believe our hack is original and innovative enough for us to take this risk and showcase all the possibilities of our idea!

## Accomplishments that we're proud of

- Getting it to work
- Making it stable enough to be used for generating policies

## What we learned

- Increased our understanding and interest in Kubernetes
- The problem that Datree is solving and a bit about how it works under the hood

## What's next for Datree Automated

We want to integrate our tool into the existing Datree CLI package and provide an automated solution for every user who wants to build a custom policy.
Our tool can also be built into a web app that can generate a Datree policy from a user's Kubernetes configuration file.

Thanks for checking out our hack! 🚀

_We would love to talk all about it or if you run into any issues trying to run it_

## Author

Kartik Choudhary \\ [kartikcho](https://github.com/kartikcho)
Saurav M H \\ [sauravhiremath](https://github.com/sauravhiremath)
