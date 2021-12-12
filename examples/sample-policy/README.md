# Policy: Checking Container Size and RunAsUSer

__This policy helps to enforce the following labels best practices:__
* [Ensure appropriate size of the container]
* [Ensure RunAsUser ]
* 
## Ensure Appropriate container size
Having correct container size is very important for performing bulk operations or for filtering workloads according to their stage. Type used:
* `g3.k3s.small`
* `g3.k3s.medium`
* `g3.k3s.large`

### When this rule is failing?
If the size key is 'g3.k3s.large`:  

```
kind: Deployment
metadata:
  labels:
    size: g3.k3s.large
```

## Ensure RunAsUser


### When this rule is failing?
If the `runAsUser` key is missing from the labels section:  
```
kind: Deployment
spec:
  containers:
    securityContext:
      runAsGroup: 3000
    
```

