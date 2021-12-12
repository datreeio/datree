# Policy: security_best_practices
This policy allows DevOps engineers or developers to enforce industry-standard recommended security
best practices such as setting user permissions and more. This policy requires the `kind`
to be set to Pod, Deployment, StatefulSet or DaemonSet.

__This policy helps to enforce the following security best practices:__
* Ensuring the user permission is set
* Ensuring the groups permission is set
* Ensuring a privilege escalation flag is set

## Ensure user permission is set
This flag controls which user ID the containers are run with. 

### When this rule is failing?
If the `runAsUser` key is missing from the securityContext section:  
```
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
```

__OR__ a invalid value of `runAsUser` is used:
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
    runAsGroup: 3000
```

## Ensure privelege escalation is set
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

## Policy author
Antariksh Verma \\ [yummyweb](https://github.com/yummyweb)

## Policy Sources
- [Kubernetes Documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
- [AquaSec Security Best Practices](https://www.aquasec.com/cloud-native-academy/kubernetes-in-production/kubernetes-security-best-practices-10-steps-to-securing-k8s/)
