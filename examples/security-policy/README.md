# Policy: security_best_practices
Kubernetes security allows users to take security measures using CLI or using configuration files. This policy describes the best practices for ensuring security of kubernetes cluster.

This policy helps to ensure following security best practices:
- [Use image digest instead of tags](#use-image-digest-instead-of-tags)
- [Ensure read only root file system is true](#ensure-read-only-root-file-system-is-true)
- [Ensure containers run with non root access](#ensure-containers-run-with-non-root-access)
- [Ensure no new privileges are set](#ensure-no-new-privileges-are-set)
- [Ensure security context is set](#ensure-security-context-is-set)
- [Missing seLinuxOptions in the securityContext](#missing-selinuxoptions-in-the-securitycontext)
- [Ensure deprecated Pod Security Policy not used](#ensure-deprecated-pod-security-policy-not-used)
- [Ensure imagePullPolicy set to Always](#ensure-imagepullpolicy-set-to-always)
- [Ensure default service account is not used](#ensure-default-service-account-is-not-used)
- [Check no default service account is used](#check-no-default-service-account-is-used)
- [Ensure all capabilities are droped](#ensure-all-capabilities-are-droped)
- [Ensure seccompProfile is set](#ensure-seccompprofile-is-set)

## Use image digest instead of tags

Always use image digests for image in kubernetes cluster. Tags are mutable so there is a chance of getting a wrong or malicious image. On the other hand, image digests are immutable. This rule will check if image digest used or not.

### When this rule is failing?

If the image used in configuration file do not have image digest.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.2.4 # tag is used
    ports:
    - containerPort: 80
```
## Ensure read only root file system is true

It is good practice to have the root file system read-only i.e. you can't modify the root folder. This rule will check if your configuration file has 'readOnlyRootFilesystem' = true.

### When this rule is failing?

If 'readOnlyRootFilesystem' is not set to true.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:latest
    securityContext:
        readOnlyRootFilesystem: false
    ...
```
## Ensure containers run with non root access

Letting the containers have root access is a bad security practice. This rule will check if the runAsUser is set 1000 or below.

### When this rule is failing?

If 'runAsUser' is either not set or have access > 1000
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  # securityContext:
  #   runAsUser: 1000
  ...
```
## Ensure no new privileges are set

When the allowPrivilegeEscalation is set to true it can give more privileges which is a security problem. This rule ensure allowPrivilegeEscalation is set to false.

### When this rule is failing?

If allowPrivilegeEscalation is set to true or not set.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
    securityContext:
      allowPrivilegeEscalation: true
```
## Ensure security context is set

Security context in container and pod configurations is almost essential to keep the cluster secure. So, this rule ensure that all forms of pod and container configurations have set their security context.

### When this rule is failing?

If security context is not set.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  labels:
    app: test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      # securityContext:
      #     # runAsUser: 1000
      #     runAsGroup: 3000
      #     fsGroup: 2000
      containers:
      - name: test
        image: busybox@sha256:23whsy235hshlkjglho355
        imagePullPolicy: Always
        # securityContext:
        #     readOnlyRootFilesystem: true
        #     runAsUser: 5000
        #     seLinuxOptions:
        #       level: RunDefault
        ports:
        - containerPort: 80
        command: [ "sh", "-c", "sleep 1h" ]

```

## Missing seLinuxOptions in the securityContext
seLinuxOptions gives more security to cluster. Note you have to enable seLinux for it. This rule is disabled by default. This checks if seLinuxOptions is set or not.

### When this rule is failing?

If seLinuxOptions is not set.
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  securityContext:
      seLinuxOptions:
          level: "s0:c123,c456"
    # ...
```

## Ensure deprecated Pod Security Policy not used

Pod security Policy is now deprecated. So don't use it. This rule checks whether Pod Security policy is not used.

### When this rule is failing?
If Pod security Policy is initiated in the config file
```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: example
spec:
# ...
```

## Ensure imagePullPolicy set to Always

It is good practice to set imagePullPolicy to Always. You can also use AlwaysPullImage admission controller. This rule ensure that imagePullPolicy is set to 'Always'.

### When this rule is failing?
If imagePullPolicy is set to other options instead of 'Always'

```yaml
piVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox
    command: [ "sh", "-c", "sleep 1h" ]
    imagePullPolicy: IfNotPresent
    volumeMounts:
    - name: sec-ctx-vol
```

## Ensure default service account is not used
This rule ensure if automountServiceAccountToken is set to false.

### When this rule is failing?
If automountServiceAccountToken is not set to false or if not set.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: 2015-06-16T00:12:59Z
  name: build-robot
  namespace: default
  resourceVersion: "272500"
  uid: 721ab723-13bc-11e5-aec2-42010af0021e
# automountServiceAccountToken: false
```

## Check no default service account is used

This is similar to the above rule but for pods. It checks whether `automountServiceAccountToken` is set to `false` for pod or not.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  serviceAccountName: build-robot
  automountServiceAccountToken: true
  volumes:
  - name: sec-ctx-vol
  [...]
```

## Ensure all capabilities are droped
It is recommended to drop all capabilities first and then add capabilities that requires for the configuration. This rule checks if all the capabilities are droped or not.

### When this rule is failing?
If all capabilities are not droped.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
    # securityContext:
    #   capabilities:
    #       drop:
    #           - all
```

## Ensure seccompProfile is set

seccomp stands for Secure Computing Mode. It is used to sandbox the privileges of a process, restricting the calls it is able to make from userspace. This rule checks if seccompProfile is set.

### When this rule is failing?
If seccompProfile is not set.
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  serviceAccountName: build-robot
  automountServiceAccountToken: false
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox
    command: [ "sh", "-c", "sleep 1h" ]
    imagePullPolicy: Always
    # securityContext:
    #     seccompProfile:
    #         type: RuntimeDefault
```


# Policy author

Abhradeep Chakraborty \ [Abhra303](https://github.com/Abhra303)