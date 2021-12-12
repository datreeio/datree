# Datree Custom Policy Example
Datree policies help kubernetes administrators to enforce security and best practices policies on resource manifest files. It comes with a good set of default policies which can be found [here](https://hub.datree.io/built-in-rules).

In this example folder, we are showcasing how to create your own custom policies that satisfy your use case. Datree custom rule engine is based on powerful [JSON Schema](https://json-schema.org/ "https://json-schema.org/"), so it supports both YAML and JSON declarative syntax.

This sample policy file included in this folder defines three custom policies:
* [governance_policy](#policy-governance_policy)
* [security_policy](#policy-security_policy)
* [stability_policy](#policy-stability_policy)

## Policy: governance_policy

This policy demonstrates how you can create rules that follow your organizations best practices.

__Rules included:__
* [Ensure each container has a configured CPU and memory less than max value](#ensure-each-container-has-a-configured-cpu-and-memory-less-than-max-value)
* [Ensure that we do not have more than 3 containers per pod](#ensure-that-we-do-not-have-more-than-3-containers-per-pod)

### Ensure each container has a configured CPU and memory less than max value
An organization might have max limit on CPU and memory value possible for any container. This rule can be used to enforce that limit. In this example policy, the max value of CPU is set to `500m` and for Memory is set to `512Mi`.

**When this rule is failing?**

If the value of cpu/memory in requests OR limits section of resource is higher than max value defined in the rule.
```
resources:
  requests:
    memory: "64Mi"
    cpu: "250m"
  limits:
    memory: "128Mi"
    cpu: "1300m"
 ```

### Ensure that we do not have more than 3 containers per pod
Pods can have multiple tightly coupled containers. An organization might want set a limit on max number of containers allowed in one pod for better policy management. In this rule, we have set the max number of containers in one pod to 3.

**When this rule is failing?**

When more than three containers are defined in the pod yaml file, this rule would fail.

## Policy: security_policy

This policy demonstrates how you can create rules that follow security best practices.

__Rules included:__
* [Ensure that we use internal image registry for our docker images](#ensure-that-we-use-internal-image-registry-for-our-docker-images)
* [Ensure that we run the containers as non-root users and runAsGroup & fsGroup values are set](#ensure-that-we-run-the-containers-as-non-root-users-and-runasgroup--fsgroup-values-are-set)
* [Ensure that we run the containers in unprivileged mode](#ensure-that-we-run-the-containers-in-unprivileged-mode)
* [Ensure that we do not AllowPrivilegeEscalation for production containers](#ensure-that-we-do-not-allowprivilegeescalation-for-production-containers)

### Ensure that we use internal image registry for our docker images
Organization should restrict the image registry allowed for pulling docker images. Unknown registries might be holding images with severe security vulnerabilities. This rule can reduce the threat exposure caused by pulling risky images.
Organizations can also choose to modify this rule to enforce that developers pull images from a particular publisher only. This publisher can be company's user account storing custom images.

**When this rule is failing?**

When image name does not start with `eu.gcr.io`, `asia.gcr.io` or `us.gcr.io`

```
containers:
  - name: nginx
    image: nginx:latest
 ```

### Ensure that we run the containers as non-root users and runAsGroup & fsGroup values are set
It is a safe practice to not allow containers to run as root user. This rule ensures that the value of `runAsUser` is greater than 0. The rule also enforces the developers to define `runAsGroup` and `fsGroup` values in the `securityContext` section.

**When this rule is failing?**

When the manifest file specifies `runAsUser` less than 1 OR `runAsGroup`/`fsGroup` are missing

```
securityContext:
    runAsUser: 0
```

### Ensure that we run the containers in unprivileged mode
When the container is running in privileged mode, then it can access host's resources and kernel capabilities. An attacker can use this privilege to exploit our infrastructure. The rule enforces that no container runs with `privileged: true` setting.

**When this rule is failing?**

When the manifest file specifies `privileged: true` or any non-boolean value.
```
containers:
  - name: nginx
    image: nginx:latest
    securityContext:
      privileged: true
```

### Ensure that we do not AllowPrivilegeEscalation for production containers
This rule is included to demonstrate the if/then condition of json schema. This rule enforces that `environment` label is set for the pod, and if the environment is production, then `allowPrivilegeEscalation` should be explicitly set to `false`. This behavior is required to effectively enforce `MustRunAsNonRoot`.

**When this rule is failing?**

When `environment` key is not included in list of labels
```
metadata:
  name: backend
  labels:
    name: backend
```
OR When `allowPrivilegeEscalation` is true or not defined for a production environment
```
name: app
  image: asia.gcr.io/myproject/app:v3
  securityContext:
    allowPrivilegeEscalation: true
```

## Policy: stability_policy

This policy demonstrates how you can create rules that ensure the stability of our kubernetes clusters.

__Rules included:__
* [Ensure that we do not configure process namespace sharing for a pod](#ensure-that-we-do-not-configure-process-namespace-sharing-for-a-pod)
* [Ensure that we do not use latest tag for our docker images](#ensure-that-we-do-not-use-latest-tag-for-our-docker-images)

### Ensure that we do not configure process namespace sharing for a pod
This directive allows processes and filesystems in a container to be visible to all other containers in that pod. It can sometimes cause unexpected behavior. This rule ensures that `shareProcessNamespace` is not set to `true`

**When this rule is failing?**

When `shareProcessNamespace` is set to true or a non-boolean value.
```
spec:
  shareProcessNamespace: true
```

### Ensure that we do not use latest tag for our docker images
It is a best practice to not include `latest` tag in image name. This rule checks that the image name does not end with `latest` tag.

**When this rule is failing?**

When latest tag is used in image name
```
containers:
  - name: nginx
    image: nginx:latest
```
## Policy author
Vijay Nandwani \\ [vijaynandwani](https://github.com/vijaynandwani/)
