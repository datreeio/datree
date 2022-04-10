# Recommended policies for production workloads
It is recommended for the kubernetes workloads deployed in production to support the industry best practices and have hardened rules enforced on them.

__The following policies are named as per their respective use cases:__
* Security

## Security

#### WORKLOAD_INVALID_CONTAINER_CAPABILITIES_VALUE
This rule ensures that no container is allowed to use the insecure/dangerous linux capabilities which causes the applications in the cluster and the cluster itself to be vulnerable to security attacks.

##### When this rule is failing?
When any kubernetes workload uses the follwing linux capabilities.
- ALL
- SYS_ADMIN
- NET_ADMIN
- CHOWN
- DAC_OVERRIDE
- FSETID
- FOWNER
- MKNOD
- NET_RAW
- SETGID
- SETUID
- SETFCAP
- SETPCAP
- SYS_CHROOT
- KILL
- AUDIT_WRITE  

```
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo-4
spec:
  containers:
  - name: sec-ctx-4
    image: gcr.io/google-samples/node-hello:1.0
    securityContext:
      capabilities:
        add: [ "NET_ADMIN","SYS_ADMIN"]
```
## Policy author
Dhanush M / [dhanushm7](https://github.com/dhanushm7)
