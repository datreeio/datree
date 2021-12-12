# Policy: preventing_use_of_removed_deprecated_apis

Deprecated apis are those that are not under active development and are scheduled to
be removed/stopped being supported in the future versions. The rules in this policy help detect the use of those apis that have been dropped in the version of kubernetes that is mentioned in the corresponding rule name and prevent them from landing in the kubernetes configuration file.

# How this helps

Many times developers don't have time to keep track of all the breaking changes happening within a software ecosystem. In this context for example a developer
might not even know that the `apiVersion` he/she is using has been deprecated/removed
from the newer kubernetes version. These rules would help identify what's breaking his
kubernetes configuration and also teach him to stop using those apis in future versions
of kubernetes.

# List of rules

- [Prevent deprecated APIs in Kubernetes v1.26 (v1.26 not yet released)]
- [Prevent deprecated APIs in Kubernetes v1.25 (v1.25 not yet released)]
- [Prevent deprecated APIs in Kubernetes v1.22]

# How these rules work

All of the three rules are exactly the same except for the kubernetes versions they are defined for and the list of values they check against.The rules check for two specific properties `apiVersion` and `kind` on each of the kubernetes config files (yaml files). All the deprecated values for `apiVersion` are matched with the corresponding values for `kind` as dictated by the kubernetes schema validation. Here are those values :-

#### In v1.26

`apiVersion` :-

-autoscaling/v2beta2

`kind` :-

-HorizontalPodAutoscaler

#### In v1.25

`apiVersion` :-

-batch/v1beta1
-discovery.k8s.io/v1beta1
-events.k8s.io/v1beta1
-autoscaling/v2beta1
-policy/v1beta1
-node.k8s.io/v1beta1

`kind` :-

-CronJob
-EndpointSlice
-Event
-HorizontalPodAutoscaler
-PodDisruptionBudget
-PodSecurityPolicy
-RuntimeClass

#### In v1.22

`apiVersion` :-

-admissionregistration.k8s.io/v1beta1
-apiextensions.k8s.io/v1beta1
-apiregistration.k8s.io/v1beta1
-authentication.k8s.io/v1beta1
-authorization.k8s.io/v1beta1
-certificates.k8s.io/v1beta1
-coordination.k8s.io/v1beta1
-extensions/v1beta1
-networking.k8s.io/v1beta1
-rbac.authorization.k8s.io/v1beta1
-scheduling.k8s.io/v1beta1
-storage.k8s.io/v1beta1

`kind` :-

-MutatingWebhookConfiguration
-CustomResourceDefinition
-APIService
-TokenReview
-LocalSubjectAccessReview
-SelfSubjectAccessReview
-SubjectAccessReview
-CertificateSigningRequest
-Lease
-Ingress
-IngressClass
-ClusterRole
-ClusterRoleBinding
-Role
-RoleBinding
-PriorityClass
-CSIDriver
-CSINode
-StorageClass
-VolumeAttachment

Whenever a match is found the `datree` cli throws an error and a message describing the failure.

## Example

### api `batch/v1beta1` is not supported in v1.25 for kind `CronJob`

#### YAML

```
apiVersion: batch/v1beta1
kind: CronJob
```

### Output on datree cli

![image](https://github.com/xoldyckk/image/blob/main/customError.png?raw=true)

### api `batch/v1` recommended by kubernetes to be used with kind `CronJob` v1.25 & onwards

#### YAML

```
apiVersion: batch/v1beta1
kind: CronJob
```

## Policy author

Adarsh Tiwari \\ [xoldyckk](https://github.com/xoldyckk)

# Other helpful documentation related to this :-

- [Kubernetes Documentation](https://kubernetes.io/docs/reference/using-api/deprecation-guide/)
