# Policy: security_checks
This policy checks for container level security parameters like privilege escalation and root filesystem permissions. 
Developers and administrators usually skip over these parameters, since ignorance of these parameters can result in severe exploits therefore it is essential to check for them before deploying manifest files in production.
Due to their severity these parameters are also a part of NSA Kubernetes Hardening Guide. More about NSA Kubernetes Hardening Guide can be read [here.](https://media.defense.gov/2021/Aug/03/2002820425/-1/-1/1/CTR_KUBERNETESHARDENINGGUIDANCE.PDF)

__This policy enforces the following best practices-__

```
spec:
  containers:
    securityContext:
      allowPrivilegeEscalation: false
```

__and__

```
spec:
  containers:
    securityContext:
      readOnlyRootFilesystem: true
```

### Resources scanned by this policy
- Deployment
- StatefulSet
- Pods

### Description of these parameters/controls
___Privilege Escalation___ - Attackers may gain access to a container and uplift its privilege to enable excessive capabilities.

_Solution/Remediation_ - If your application does not need it, make sure the allowPrivilegeEscalation field of the securityContext is set to false..


___Immutable Container filesystem___ - Mutable container filesystem can be adjusted to inject malicious code or data into containers. Use immutable (read only) filesystem to limit potential attacks.

_Solution/Remediation_ - Set the readOnlyFilesystem field of securityContext to true. In case your application requires writable filesystem then it is recommended to mount secondary filesystems
