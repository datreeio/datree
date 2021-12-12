# Policy: security_best_practices
This policy allows DevOps engineers or developers to enforce industry-standard recommended security
best practices such as setting user permissions and more. This policy requires the `kind`
to be set to Pod, Deployment, StatefulSet or DaemonSet.

__This policy helps to enforce the following security best practices:__
* Ensure user permission is set
* Ensure groups permission is set
* Ensure privilege escalation is set
* Ensure that filesystem accessibility is set
* Ensure deprecated resource is not used
* Ensure privileged flag is set
* Ensure that imagePullPolicy is set to Always

## Ensure user permission is set
This flag controls which user ID the containers are run with. If this value is 0, then
the pod is run with non-root user.

### When this rule is failing?
If the `runAsUser` key is missing from the securityContext section:  
```
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
```

__OR__ an invalid value of `runAsUser` is used:
```
kind: Pod
spec:
  securityContext:
    runAsUser: 11000
```

## Ensure groups permission is set
The `runAsGroup` field controls how processes interact with files. If it was ommitted,
the gid would remain as 0, and the process would be able to interact with files owned
by root group.

### When this rule is failing?
If the `runAsGroup` key is missing from the securityContext section:  
```
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
    allowPrivilegedEscalation: false
```

__OR__ an invalid value of `runAsGroup` is used:
```
kind: Pod
spec:
  securityContext:
    runAsGroup: -1
```

## Ensure privilege escalation is set
This flag controls whether process can gain more privileges than it's parent process.
It controls the `no_new_privs` flag which is set on the container process.

### When this rule is failing?
If the flag is not a boolean: 
```
kind: Pod
spec:
  securityContext:
    allowPrivilegeEscalation: false
```

## Ensure that filesystem accessibility is set
Requires containers to run with a read-only root filesystem if set to true i.e. no
writeable layer.

### When this rule is failing?
If the flag is not a boolean:
```
kind: Pod
spec:
  securityContext:
    readOnlyRootFilesystem: false
```

## Ensure deprecated resource is not used
PodSecurityPolicy is a type of resource that is deprecated in Kubernetes v1.21 and is
set to be removed in Kubernetes v1.25. It is against security best practices and other
resources can be used in favor of it such as Pod or Deployment.

### When this rule is failing?
If the `kind` field is set to PodSecurityPolicy:
```
kind: PodSecurityPolicy
```

## Ensure privileged flag is set
The `privileged` flag allows a container to run in privileged mode where processes are
essentially equivalent to root on the host. It also allows the process to access host's
resources and devices.

### When this rule is failing?
If the flag is not a boolean:
```
kind: Pod
spec:
  securityContext:
    privileged: not-false
```

__OR__ it is missing from the securityContext section:
```
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 1000
    allowPrivilegedEscalation: true
```

## Ensure that imagePullPolicy is set to Always
The `imagePullPolicy` field for a container affect when the kubelet attempts to download
the specified image. If the value is set to Always, then everytime the kubelet is launched
the kubelet queries the container image and resolves the name to an image digest.

### When this rule is failing?
If the field is not Always:
```
kind: Pod
spec:
  containers:
    - name: image
      imagePullPolicy: Never
```

## Policy author
Antariksh Verma \\ [yummyweb](https://github.com/yummyweb)

## Policy Sources
- [Pod Security Policy Documentation](https://kubernetes.io/docs/concepts/policy/pod-security-policy/)
- [Kubernetes Pod Configuration](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
- [Aqua Security Best Practices](https://www.aquasec.com/cloud-native-academy/kubernetes-in-production/kubernetes-security-best-practices-10-steps-to-securing-k8s/)
