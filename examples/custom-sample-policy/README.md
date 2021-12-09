# Policy: custom_sample_policy (example policy)

__This policy helps to enforce the following labels best practices:__
* [Ensure correct environment labels are used](#ensure-correct-environment-labels-are-used)
* [Ensure each container has a configured memory request](#ensure-each-container-has-a-configured-memory-request)
* [Ensure each container has a configured cpu request](#ensure-each-container-has-a-configured-cpu-request)
* [Ensure each container image has a pinned tag version](#ensure-each-container-image-has-a-pinned-tag-version)
* [Ensure each container has a configured cpu limit within range](#ensure-each-container-has-a-configured-cpu-limit-within-range)
* [Ensure Deployment has replicas set within range](#ensure-deployment-has-replicas-set-within-range)

## Ensure correct environment labels are used
Having an env label is useful for performing bulk operations in specific environments or for filtering workloads according to their stage. This rule will also ensure that only pre-approved `environment` label values are used:
* `prod`
* `staging`
* `test`

### When this rule is failing?
If the `environment` key is missing from the labels section:  
```
kind: Deployment
metadata:
  labels:
    owner: jatin
```

__OR__ a different `environment` value is used:
```
kind: Deployment
metadata:
  labels:
    environment: qa
```

## Ensure each container has a configured memory request
Each container should have a configured memory request.

### When this rule is failing?
If the `memory` key is missing from the request section:  
```
containers:
      - name: nginx
        image: nginx:1.14.2
        resources:
          requests:
            cpu: "250m"
```

## Ensure each container has a configured cpu request
Each container should have a configured memory request.

### When this rule is failing?
If the `cpu` key is missing from the request section:  
```
containers:
      - name: nginx
        image: nginx:1.14.2
        resources:
          requests:
            memory: "64Mi"
```

## Ensure each container image has a pinned tag version
Its better to specify the version of image that you want to use,
If no version is provided then by default you will get the latest image.
Which can create some issues.

### When this rule is failing?
If the format of `image` key does not contain specific versioning 
```
containers:
    - name: nginx
    image: nginx@latest
```

## Ensure each container has a configured cpu limit within range
Cpu limit should be within the specified range.

### When this rule is failing?
If the `cpu` key is having some value which is out of specific range(250m-500m):  
```
resources:
    limits:
    memory: "128Mi"
    cpu: "450m"
```

## Ensure Deployment has replicas set within range
If the kind is Deployment then the replicas should within the specified range.

### When this rule is failing?
If the `replicas` key is having some value which is out of specific range(2-10):  
```
spec:
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend 
```

## Policy author
Jatin Motwani \\ [eyarz](https://github.com/jatinmotwani)
