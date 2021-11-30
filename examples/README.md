# Hackathon Guide - December 10-12, 2021

**If at any point you find yourself lost, join our [Slack community](https://bit.ly/3BHwCEG) for assistance.**


## This is a step-by-step guide to submit one piece of work to Datree’s hackathon:


1. [Sign up](https://app.datree.io/#hackathon) for Datree and follow the instructions to install Datree’s CLI on your machine
2. Fork [Datree's project](https://github.com/datreeio/datree):
  <img width="1679" alt="Screen Shot 2021-11-07 at 11 19 57" src="https://user-images.githubusercontent.com/1208902/142754175-099d9d47-fa83-415c-bec6-c0373a65e1cc.png">
  
<img width="391" style="float: right;" alt="Screen Shot 2021-11-07 at 11 45 07" src="https://user-images.githubusercontent.com/1208902/142754350-8aca5344-08a1-4a07-aa02-a1f9e15f08e3.png"> 

3. Clone the forked project to your local machine:  
   * Click on the green code button<sup>[1]</sup>  
   * Choose the cloning method (HTTPS, SSH, CLI)<sup>[2]</sup> 
   * Copy the link<sup>[3]</sup>  
   * Open a terminal or command line  
   * Direct it to the directory where we want to store the local repo with the cd command  
   * Run: `git clone <copied link>`  
   * Run: `cd datree`  


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
   

   

