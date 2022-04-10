# Policy: deployment_best_practices 
Deployment is one the Key part of Running an Application.This policiy ensures whether K8's Admins Applied Correct Configuration for RollingUpdateDeployment.

__This policy helps to enforce the following deployment best practices:__
* [Define Rolling Update Strategy](#Define-Rolling-Update-Strategy)
* [Ensure revisionHistoryLimit was Set](#Ensure-revisionHistoryLimit-was-Set)

## Define Rolling Update Strategy

Mentioning Rolling Update can be used to adjust the `parameters` in Manifest File
* `maxSurge` define how many pod we can add at a time during the rolling update
* `maxUnavailable` define how many pod can be unavailable during the rolling update

### When this rule is failing?
A different `Deployment Strategy` value is used:
```
spec:
  replicas: 3
  strategy:
        type: Recreate
 ```
 __OR__ If `maxSurge` or `maxUnavailable` are missing
 ```
 spec:
  replicas: 3
  strategy:
        type: RollingUpdate
 ```

 ## Ensure revisionHistoryLimit was Set

If Deployment of New Version is Messed up to Roll Back to the Previous Version We Use `Kubectl` commands and to get previous History of Deployment we can add `revisionHistoryLimit` in Manifest Files to get previous versions upto certain versions

### When this rule is failing?
`revisionHistoryLimit` wasn't set to an Integer
```
spec:
 revisionHistoryLimit: NA
```
## Policy author
Avinash \\ [Avinash](https://github.com/AvinashNayak27)


