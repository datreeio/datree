# Policy: preventing_use_of_removed_deprecated_apis

Deprecated apis are those that are not under active development and are scheduled to be removed/stopped being supported in the newer/future versions. The rules in this policy help detect the use of those apiVersions that have been/will be dropped in the corresponding version of kubernetes that is mentioned in the rule name and prevent them from landing in the kubernetes configuration objects/files.

# How this helps

Many a times developers don't have time to keep track of all the breaking changes happening within a software ecosystem. In this context for example a developer
might not even know that the `apiVersion` he/she is using has been deprecated/removed from the newer kubernetes version. These rules would help identify what's breaking his kubernetes configuration and also keep him updated to these latest changes so that he/she stops using those apiVersions in their kubernetes configuration so that it works for newer kubernetes versions as well.

# List of rules

- [Prevent deprecated APIs in Kubernetes v1.26 (v1.26 not yet released)]
- [Prevent deprecated APIs in Kubernetes v1.25 (v1.25 not yet released)]
- [Prevent deprecated APIs in Kubernetes v1.22]

# How these rules work

All of the three rules are exactly the same except for the kubernetes versions they are defined for and the set of values they check against. The rules check for two specific properties `apiVersion` and `kind` on each of the kubernetes config files (yaml files). All the deprecated values for `apiVersion` are matched with the corresponding values for `kind` as dictated by the kubernetes schema validation. Here are those set of apiVersions that are/will be no longer served by kubernetes in the versions specified and onwards :-

#### In v1.26

| deprecated `apiVersion` |         `kind`          | recommended `apiVersion` |
| :---                    |          :---:          |                     ---: |
| autoscaling/v2beta2     | HorizontalPodAutoscaler |           autoscaling/v2 |

#### In v1.25

| deprecated `apiVersion`  |        `kind`           |       recommended `apiVersion` |
| :---                     |         :---:           |                           ---: |
| batch/v1beta1            |         CronJob         |                       batch/v1 | 
| discovery.k8s.io/v1beta1 |      EndpointSlice      |            discovery.k8s.io/v1 | 
| events.k8s.io/v1beta1    |         Event           |               events.k8s.io/v1 | 
| autoscaling/v2beta1      | HorizontalPodAutoscaler |                 autoscaling/v2 | 
| policy/v1beta1           |   PodDisruptionBudget   |                      policy/v1 | 
| policy/v1beta1           |   PodSecurityPolicy     | `No official replacements yet` | 
| node.k8s.io/v1beta1      |      RuntimeClass       |                 node.k8s.io/v1 | 

#### In v1.22


| deprecated `apiVersion`              |            `kind`              |        recommended `apiVersion` |
| :---                                 |            :---:               |                            ---: |
| admissionregistration.k8s.io/v1beta1 |  MutatingWebhookConfiguration  | admissionregistration.k8s.io/v1 | 
| admissionregistration.k8s.io/v1beta1 | ValidatingWebhookConfiguration | admissionregistration.k8s.io/v1 | 
| apiextensions.k8s.io/v1beta1         |    CustomResourceDefinition    |         apiextensions.k8s.io/v1 | 
| apiregistration.k8s.io/v1beta1       |          APIService            |       apiregistration.k8s.io/v1 | 
| authentication.k8s.io/v1beta1        |          TokenReview           |        authentication.k8s.io/v1 | 
| authorization.k8s.io/v1beta1         |    LocalSubjectAccessReview    |         authorization.k8s.io/v1 | 
| authorization.k8s.io/v1beta1         |    SelfSubjectAccessReview     |         authorization.k8s.io/v1 | 
| authorization.k8s.io/v1beta1         |      SubjectAccessReview       |         authorization.k8s.io/v1 | 
| certificates.k8s.io/v1beta1          |   CertificateSigningRequest    |          certificates.k8s.io/v1 | 
| coordination.k8s.io/v1beta1          |             Lease              |          coordination.k8s.io/v1 | 
| extensions/v1beta1                   |           Ingress              |            networking.k8s.io/v1 | 
| networking.k8s.io/v1beta1            |           Ingress              |            networking.k8s.io/v1 | 
| networking.k8s.io/v1beta1            |         IngressClass           |            networking.k8s.io/v1 | 
| rbac.authorization.k8s.io/v1beta1    |         ClusterRole            |    rbac.authorization.k8s.io/v1 | 
| rbac.authorization.k8s.io/v1beta1    |       ClusterRoleBinding       |    rbac.authorization.k8s.io/v1 | 
| rbac.authorization.k8s.io/v1beta1    |            Role                |    rbac.authorization.k8s.io/v1 | 
| rbac.authorization.k8s.io/v1beta1    |         RoleBinding            |    rbac.authorization.k8s.io/v1 | 
| scheduling.k8s.io/v1beta1            |        PriorityClass           |            scheduling.k8s.io/v1 | 
| storage.k8s.io/v1beta1               |          CSIDriver             |               storage.k8s.io/v1 | 
| storage.k8s.io/v1beta1               |           CSINode              |               storage.k8s.io/v1 | 
| storage.k8s.io/v1beta1               |         StorageClass           |               storage.k8s.io/v1 | 
| storage.k8s.io/v1beta1               |        VolumeAttachment        |               storage.k8s.io/v1 | 

#### Whenever a `deprecated apiVersion` ends up with the corresponding `kind` in a kubernetes configuration object `datree` cli throws an error and a message describing the failure.

## Example

### apiVersion `batch/v1beta1` will not be supported in v1.25 for kind `CronJob`

##### YAML

```
apiVersion: batch/v1beta1
kind: CronJob
```

##### Output on the datree cli

![image](https://github.com/xoldyckk/image/blob/main/customError.png?raw=true)

### apiVersion `batch/v1` is recommended by kubernetes to be used with kind `CronJob` in v1.25 & onwards

##### YAML

```
apiVersion: batch/v1
kind: CronJob
```

## Note

Kubernetes schema validation throws error of type `could not find schema for ....` whenever the `apiVersion` and `kind` properties
do not match their definition for the current kubernetes version. To prevent this behaviour I have provided two workarounds that I'm going
to list right here :-

### Workaround 1 :-

#### Run the checks according to the version specified in the comments for each config in pass/fail yaml files. For example :- 
```
# Test using v1.23 or later    <-------------   THIS COMMENT
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
```
#### You can check each config individually using the version specified in the comment (v1.23 for this case).
#### Example using the cli :- ```datree test `your config file here` --schema-version "1.23.0" ```

### Workaround 2 :-

#### Pass the `--ignore-missing-schemas` flag to the datree test.
#### Example using the cli :- ```datree test `your config file here` --ignore-missing-schemas```

## Policy author

Adarsh Tiwari \\ [xoldyckk](https://github.com/xoldyckk)

# Other helpful documentation related to this :-

- [Kubernetes Documentation](https://kubernetes.io/docs/reference/using-api/deprecation-guide/)
