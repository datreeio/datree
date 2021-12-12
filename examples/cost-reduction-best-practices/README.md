# Policy: cost_reduction_best_practices
Kubernetes resource requests and limits enable engineers to ensure that their workloads are not over or under-utilized. These policies are implemented with multiple property paths to provide adaptation to multiple Kubernetes objects such as Pod, Deployment.

__This policy helps to enforce the following best practices:__
* [Ensure each container has a configured CPU request within range](#ensure-each-container-has-a-configured-cpu-request-within-range)
* [Ensure each container has a configured CPU limit within range](#ensure-each-container-has-a-configured-cpu-limit-within-range)
* [Ensure each container has a configured memory request within range](#ensure-each-container-has-a-configured-memory-request-within-range)
* [Ensure each container has a configured memory limit within range](#ensure-each-container-has-a-configured-memory-limit-within-range)

## Ensure each container has a configured CPU request within range
### When this rule is failing?
If `requests.cpu` is missing:
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

__OR__ the value is outside of the configured range (100m-250m):
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      resources:
        requests:
          cpu: "50m"
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            requests:
              cpu: "50m"
```

## Ensure each container has a configured CPU limit within range
### When this rule is failing?
If `limits.cpu` is missing:
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

__OR__ the value is outside of the configured range (500m-1000m):
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      resources:
        limits:
          cpu: "1500m"
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: "1500m"
```

## Ensure each container has a configured memory request within range
### When this rule is failing?
If `requests.memory` is missing:
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

__OR__ the value is outside of the configured range (512Mi-1024Mi):
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      resources:
        requests:
          memory: "256Mi"

```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "256Mi"
```

## Ensure each container has a configured memory limit within range
### When this rule is failing?
If `limits.memory` is missing:
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

__OR__ the value is outside of the configured range (2048Mi-4096Mi):
```
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      resources:
        limits:
          memory: "5120Mi"
```
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            limits:
              memory: "5120Mi"
```

## Policy author
Brijesh Shah \\ [brijeshshah13](https://github.com/brijeshshah13)