# Policy: No. of container instances
  

__This policy helps to enforce the following labels best practices:__
* [Ensure suitable no of container instances]

* ## Ensure container instance is specified

### When this rule is failing?
If the `instances` key is missing from the labels section:  
```
kind: Deployment
metadata:
  labels:
    owner: yoda-at-datree.io
```

__OR__ a instances > 4 is used:
```
kind: Deployment
metadata:
  labels:
    instances: 5
```
