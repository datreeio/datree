# Policy: security_best_practices
Ingress host name should be a valid organization name. Invalid host names can cause problems with access to the API. These policies are implemented with multiple property paths to provide adaptation to multiple Kubernetes objects such as Pod, Deployment.

__This policy helps to enforce the following best practices:__
* [Prevent containers from running without a read-only root filesystem](#prevent-containers-from-running-without-a-read-only-root-filesystem)
* [Ensure containers do not allow privilege escalation](#ensure-containers-do-not-allow-privilege-escalation)
* [Ensure containers do not run processes with root privileges](#ensure-containers-do-not-run-processes-with-root-privileges)

## Prevent containers from running without a read-only root filesystem
### When this rule is failing?
If `securityContext.readOnlyRootFilesystem` is missing:
```
apiVersion: v1
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

__OR__ the value of `securityContext.readOnlyRootFilesystem` is `false`:

```
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      securityContext:
        readOnlyRootFilesystem: false
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
          securityContext:
            readOnlyRootFilesystem: false
```




## Ensure containers do not allow privilege escalation

### When this rule is failing?
If `securityContext.allowPrivilegeEscalation` is missing:
```
apiVersion: v1
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

__OR__ the value of `securityContext.allowPrivilegeEscalation` is `true`:

```
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      securityContext:
        allowPrivilegeEscalation: true
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
          securityContext:
            allowPrivilegeEscalation: true
```

## Ensure containers do not run processes with root privileges

### When this rule is failing?
If `securityContext.runAsUser` or `securityContext.runAsNonRoot` is missing:
```
apiVersion: v1
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

__OR__ the value of `securityContext.runAsUser` is greater than `0` and `securityContext.runAsNonRoot` is `true`:

```
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      securityContext:
        runAsUser: 1000
        runAsNonRoot: true
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
          securityContext:
            runAsUser: 1000
            runAsNonRoot: true
```

## Ensure containers do not expose sensitive host system directories

### When this rule is failing?
If `mountPath` is missing:
```
apiVersion: v1
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

__OR__ the value of `mountPath` is one of the directories listed by the organization:

```
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
        - containerPort: 80
      volumeMounts:
        - name: nginx-certs
            mountPath: /etc/nginx/certs
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
          volumeMounts:
            - name: nginx-certs
              mountPath: /etc/nginx/certs
```
## Policy author
Brijesh Shah \\ [brijeshshah13](https://github.com/brijeshshah13)