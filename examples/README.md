# Cloud Native Hackathon Guide - December 10-12, 2021
### Datree’s mission: To improve the workflow of Kubernetes administrators and users

**If at any point you find yourself lost, join our [Slack community](https://bit.ly/3BHwCEG) for assistance.**

### Compete to win
* 1st prize is a laptop;  
* 2nd prize is a mobile phone;  
* 3rd prize is an iPad or headphones  

## Required knowledge to participate in this Hackathon
* Basic knowledge with operating Git and GitHub
* Basic knowledge with reading and writing JSON and YAML
* Basic knowledge about using and writing [JSON Schema](https://json-schema.org/)
* Familiarity with Datree's [Policy as code](https://hub.datree.io/policy-as-code) and [custom rules](https://hub.datree.io/custom-rules) concepts

## The Rules
1. Choose a use case that would help other Kubernetes admins or users in their work. Such as:
   1. Stability
   2. Cost reduction
   3. Security
   4. Industry best practices
   5. Check lists
   6. Any other creative idea you have will be welcomed
2. Solve the use case by creating a policy with relevant custom rules. See the [step-by-step](#this-is-a-step-by-step-guide-to-submit-one-piece-of-work-to-datrees-hackathon) below
3. Include well-written documentation in a README file. The README should include the following information:
   1. An overview about the use case that the policy is solving
   2. An explanation about why it’s a worthy use case and how the use case helps Kubernetes admins in their work
   3. The rules to support the use case
   4. Details for each rule, which should include what can cause the rule to fail and how to fix it. [See an example here](https://hub.datree.io/ensure-labels-value-valid).
4. Provide a basic Kubernetes manifests to test (failing and passing) each custom rule. [See an example here](https://github.com/datreeio/datree/tree/main/examples/sample-policy)
5. Please note: your policy must represent a workable solution.

## How to win
1. Submit your work by opening a pull request on [Datree’s repo](https://github.com/datreeio/datree). You may submit more than one policy.
2. Your score for winning will be based on whether your use case works according to what you present, including detailed documentation, and how useful the use case is to Kubernetes admins. Quality matters over quantity.
3. (Optional) Are you proud of your work? If so, don’t forget to Tweet or post on LinkedIn that you’ve participated in this #CloudNativeHackathon and tag us so we can like and share!

## Support
1. If technical issues arise, support will be provided via our [Slack channel](https://bit.ly/3BHwCEG). 
2. If you find a bug, open an issue in [our project](https://github.com/datreeio/datree/issues) -- we’ll prioritize it and answer.

## YouTube Video
[![You Tube Video](https://img.youtube.com/vi/Cgmvs3UFPIQ/0.jpg)](https://www.youtube.com/watch?v=Cgmvs3UFPIQ)
   
## This is a step-by-step guide to submit one piece of work to Datree’s hackathon

1. [Sign up](https://app.datree.io/#hackathon) for Datree and follow the instructions to install Datree’s CLI on your machine
2. Fork [Datree's project](https://github.com/datreeio/datree):
  <img width="1679" alt="Screen Shot 2021-11-07 at 11 19 57" src="https://user-images.githubusercontent.com/1208902/142754175-099d9d47-fa83-415c-bec6-c0373a65e1cc.png">

3. Clone the forked project to your local machine:  
   * Click on the green code button<sup>[1]</sup>  
   * Choose the cloning method (HTTPS, SSH, CLI)<sup>[2]</sup> 
   * Copy the link<sup>[3]</sup>  
   * Open a terminal or command line  
   * Direct it to the directory where we want to store the local repo with the cd command  
   * Run: `git clone <copied link>`  
   * Run: `cd datree`  

<img width="391" style="float: right;" alt="Screen Shot 2021-11-07 at 11 45 07" src="https://user-images.githubusercontent.com/1208902/142754350-8aca5344-08a1-4a07-aa02-a1f9e15f08e3.png"> 

4. Create your policy:
  
   * Open the project in your favourite IDE (e.g. [VSCode](https://code.visualstudio.com/download), [WebStorm](https://www.jetbrains.com/webstorm/download))
   * Go to the examples directory.
   * Create a new directory with the name of your policy. The code structure should resemble the [sample-policy](https://github.com/datreeio/datree/tree/main/examples/sample-policy) directory
   * Go to your newly created directory
   * Create a policy with custom rules using [this](https://hub.datree.io/custom-rules) guide
   * Create a README.md file that describes the policy and rules
   * Create a pass.yaml and fail.yaml that contain a k8s manifest example that pass and fail your policy respectively. You can create more than one of each.

5. Save the changes - once you made the changes save them with git locally by committing them:
   * Run `git add .` to add all changed files.  
     You can select specific files by running `git add [file1] [file2]`
   * Run `git commit -m "[meaningful commit message]"`
6. Run `git push` to push the changes to your remote (forked) repository
7. When you’re satisfied with your work and you’re ready to submit it for review, create a pull request:
   * Go to your forked repository in GitHub.
   * You should see the **Compare & pull request** notification. Click on it
   ![Screen Shot 2021-11-08 at 15 00 40](https://user-images.githubusercontent.com/1208902/142755492-57262458-87d1-4f3c-a9e7-00bf00d1313c.png)
    *If you can’t find it. Go to the pull requests tab and click the New Pull Request button
   * Write a title and a description using [this](https://github.com/datreeio/datree/blob/ada7baa263f7dee8b43c99bc50868bf6b90e0857/CONTRIBUTING.md#-commit-message-format) guide

   * Click the "compare & pull request" button
   * Please note, the pull request is visible in the [origin repository](https://github.com/datreeio/datree)
8. The team will review and approve the request or will ask for changes and clarifications.
   * Please visit the pull request page often and see if any changes were requested.  
   

   

