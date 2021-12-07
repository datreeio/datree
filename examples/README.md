# Cloud Native Hackathon Guide - December 10-12, 2021
### Datree’s mission: To improve the workflow of Kubernetes administrators and users

**If at any point you find yourself lost, join our [Slack community](https://bit.ly/3BHwCEG) for assistance.**

### :trophy: Compete to win
* 1st prize is a laptop  
* 2nd prize is a mobile phone  
* 3rd prize is an iPad or headphones  

## :mortar_board: Required knowledge to participate in this Hackathon
* Basic knowledge with operating Git and GitHub
* Basic knowledge with reading and writing JSON and YAML
* Basic knowledge about using and writing [JSON Schema](https://json-schema.org/)
* Familiarity with Datree's [Policy as code](https://hub.datree.io/policy-as-code) and [custom rules](https://hub.datree.io/custom-rules) concepts

## :clipboard: Instructions
1. Choose a use case that would help other Kubernetes admins or users in their work. Such as:  
   a. Stability  
   b. Cost reduction  
   c. Security  
   d. Industry best practices  
   e. Governance  
   f. Any other creative idea you have will be welcomed  
   
_:bulb:ㅤWe also welcome all innovative uses of Datree that help the community - please get your innovative idea pre-approved by our judges_  

2. Solve the use case by [creating a policy](#how-to-write-a-policy) with relevant custom rules
3. Include well-written documentation in a README file   
__The README should include the following information:__  
   ㅤa. An overview about the use case that the policy is solving  
   ㅤb. An explanation about why it’s a worthy use case and how the use case helps Kubernetes admins in their work  
   ㅤc. A list of the rules in the policy and how they support the use case  
   ㅤd. Any other type of documentation that can help evaluate and understand the rules in your policy, and how they are related to the use case that you chosed
4. Provide a basic Kubernetes manifests to test (failing & passing) each custom rule ([see an example here](https://github.com/datreeio/datree/tree/main/examples/sample-policy))  
5. Submit your work bt opening a Pull Request to the examples dirctory (`datree/examples`)

_:warning:ㅤPlease note: your policy must represent a workable solution_

## :tada: How to win
1. Submit your work by opening a pull request on [Datree’s repo](https://github.com/datreeio/datree) (you may submit more than one policy).
2. Your score for winning will be based on whether your use case works according to what you present, including detailed documentation, and how useful the use case is to Kubernetes admins. Quality matters over quantity.
3. [OPTIONAL] Are you proud of your work? If so, don’t forget to Tweet or post on LinkedIn that you’ve participated in this #CloudNativeHackathon and tag us so we can like and share!

### YouTube Video
[![You Tube Video](https://img.youtube.com/vi/Cgmvs3UFPIQ/0.jpg)](https://www.youtube.com/watch?v=Cgmvs3UFPIQ)

## :oncoming_police_car: How to write a policy
1. [Sign up](https://app.datree.io/#hackathon) for Datree and follow the instructions to install Datree’s CLI on your machine
2. Using JSON Schema, create [custom rules](https://hub.datree.io/custom-rules-overview) that are relevant to your policy  
:point_right:ㅤYou can use this [online YAML Schema Validator](https://yamlschemavalidator.datree.io/) to easily test your custom rule logic before adding it to your policy
3. Add the custom rules to your [policy file](https://hub.datree.io/policy-as-code#go-policiesyaml)
4. Publish the policy (`datree publish policy-name.yaml`) and verify it is working as expected

## :computer: How to submit a pull request
1. Fork [Datree's project](https://github.com/datreeio/datree):
  <img width="1679" alt="Screen Shot 2021-11-07 at 11 19 57" src="https://user-images.githubusercontent.com/1208902/142754175-099d9d47-fa83-415c-bec6-c0373a65e1cc.png">

2. Clone the forked project to your local machine:  
   * Click on the green code button<sup>[1]</sup>  
   * Choose the cloning method (HTTPS, SSH, CLI)<sup>[2]</sup> 
   * Copy the link<sup>[3]</sup>  
   * Open a terminal or command line  
   * Direct it to the directory where we want to store the local repo with the cd command  
   * Run: `git clone <copied link>`  
   * Run: `cd datree`  

<img width="391" style="float: right;" alt="Screen Shot 2021-11-07 at 11 45 07" src="https://user-images.githubusercontent.com/1208902/142754350-8aca5344-08a1-4a07-aa02-a1f9e15f08e3.png"> 

3. Add your policy:
   * Open the project in your favourite IDE (e.g. VSCode, WebStorm, etc.)
   * Go to the examples directory (`datree/examples`) and create a new directory with the name of your policy 
   * Add your policy and additinal required files (README, tests yamls, etc.)
   * The code structure should resemble the [sample-policy](https://github.com/datreeio/datree/tree/main/examples/sample-policy) directory
   
4. Save the changes - once you made the changes save them with git locally by committing them:
   * Run `git add .` to add all changed files.  
     You can select specific files by running `git add [file1] [file2]`
   * Run `git commit -m "[meaningful commit message]"`
5. Run `git push` to push the changes to your remote (forked) repository
6. When you’re satisfied with your work and you’re ready to submit it for review, create a pull request:
   * Go to your forked repository in GitHub
   * Click on the "**Compare & pull request**" button
   ![Screen Shot 2021-11-08 at 15 00 40](https://user-images.githubusercontent.com/1208902/142755492-57262458-87d1-4f3c-a9e7-00bf00d1313c.png)
    *Please note: the pull request is visible in the [origin repository](https://github.com/datreeio/datree)
   * Write a title and a description using [this](https://github.com/datreeio/datree/blob/ada7baa263f7dee8b43c99bc50868bf6b90e0857/CONTRIBUTING.md#-commit-message-format) guide
7. The team will review and approve the request or will ask for changes and clarifications.
   * Please visit the pull request page often and see if any changes were requested.  
   
## :ambulance: Support
* If technical issues arise, support will be provided via our [Slack channel](https://bit.ly/3BHwCEG). 
* If you find a bug, open an issue in [our project](https://github.com/datreeio/datree/issues) -- we’ll prioritize it and answer.

