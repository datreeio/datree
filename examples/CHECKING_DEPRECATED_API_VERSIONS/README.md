# Policy: CHECKING_DEPRECATED_API_VERSIONS

A deprecated API is one that you are no longer recommended to use, due to changes in the API. While deprecated classes, methods, and fields are still implemented, they may be removed in future implementations, so you should not use them in new code, and if possible rewrite old code not to use them. This policy checks for the deprecated API Versions used in the Kubernetes Configuration Files.

**Objectives of this policy:**

- It prevents all the deprecated APIs in kubernetes cluster version v1.22 and above to reach the production.
- It gives an alert for the APIs which are expected to be deprecated in the kubernetes cluster version v1.26 and above.

**What makes this policy different from the existing comparable datree policies(deprecation policies)**

- Existing datree policies checks for the deprecated APIs in and above kubernetes cluster version v1.17.
- The newly defined policy checks for the deprecated APIs in and above kubernetes cluster version v1.22.
- It also warns for the APIs which are expected to be deprecated by kubernetes cluster version v1.26.

## Policy Explanation

### Custom rule 1: Prevent deprecated APIs in kubernetes v1.22

Targeted resources by this rule (types of `kind`):

- MutatingWebhookConfiguration
- ValidatingWebhookConfiguration
- CustomResourceDefinition
- APIService
- TokenReview
- LocalSubjectAccessReview
- SelfSubjectAccessReview
- SubjectAccessReview
- CertificateSigningRequest
- Lease
- Ingress
- IngressClass
- ClusterRole
- ClusterRoleBinding
- Role
- RoleBinding
- PriorityClass
- CSIDriver
- CSINode
- StorageClass
- VolumeAttachment

### When this rule is failing?

If one of the following API versions is used

- admissionregistration.k8s.io/v1beta1
- apiextensions.k8s.io/v1beta1
- apiregistration.k8s.io/v1beta1
- authentication.k8s.io/v1beta1
- authorization.k8s.io/v1beta1
- certificates.k8s.io/v1beta1
- coordination.k8s.io/v1beta1
- extensions/v1beta1
- networking.k8s.io/v1beta1
- rbac.authorization.k8s.io/v1beta1
- scheduling.k8s.io/v1beta1
- storage.k8s.io/v1beta1

**YAML**

```
   apiVersion: batch/v1beta1
   kind: CronJob
   metadata:
     name: hello
```

**Output on CLI**

![image](https://user-images.githubusercontent.com/68479079/145676250-88f0624b-f69a-488c-b1b7-569eee81c1c8.png)

### How to fix this failure

Use v1 instead of the deprecated version

**YAML**

```
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
```

##

### Custom rule 2: [Warning -> Not yet deprecated] Prevent deprecated APIs in kubernetes v1.25

Targeted resources by this rule (types of **`kind`**):

- CronJob
- EndpointSlice
- Event
- HorizontalPodAutoscaler
- PodDisruptionBudget
- PodSecurityPolicy
- RuntimeClass

### When this rule is failing?

If one of the following API versions is used

- batch/v1beta1
- discovery.k8s.io/v1beta1
- events.k8s.io/v1beta1
- autoscaling/v2beta1
- policy/v1beta1
- node.k8s.io/v1beta1

**YAML**

```
   apiVersion: admissionregistration.k8s.io/v1beta1
   kind: ValidatingWebhookConfiguration
   metadata:
     name: "pod-policy.example.com"
```

**Output on CLI**

![image](https://user-images.githubusercontent.com/68479079/145676540-a596b386-6cd6-463f-aee9-399541a7296f.png)

### How to fix this failure

Use v1 instead of the deprecated version

**YAML**

```
   apiVersion: admissionregistration.k8s.io/v1
   kind: ValidatingWebhookConfiguration
   metadata:
     name: "pod-policy.example.com"
```

##

### Custom rule 3: [Warning -> Not yet deprecated] Prevent deprecated APIs in kubernetes v1.26

Targeted resources by this rule (types of **`kind`**):

- HorizontalPodAutoscaler

### When this rule is failing?

If one of the following API versions is used

- autoscaling/v2beta2

**YAML**

```
   apiVersion: autoscaling/v2beta2
   kind: HorizontalPodAutoscaler
   metadata:
      name: nginx
```

**Output on CLI**

![image](https://user-images.githubusercontent.com/82814375/145682522-06c32061-ca09-42f2-bba3-e52e3cefa4b5.jpeg)

### How to fix this failure

Use v2 instead of the deprecated version

**YAML**

```
   apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
      name: nginx
```

## Policy authors

- Rupesh Agarwal \\ [Rupesh-1302](https://github.com/Rupesh-1302)
- Mahesh Gajakosh \\ [CodeAbsolute](https://github.com/CodeAbsolute)
- Nidhi Daulat \\ [nidhidaulat16](https://github.com/nidhidaulat16)
- Ranjeet Suthar \\ [RanjeetNSuthar](https://github.com/RanjeetNSuthar)

## Policy Sources

- [Kubernetes Documentation](https://kubernetes.io/docs/reference/using-api/deprecation-guide/#what-to-do)
