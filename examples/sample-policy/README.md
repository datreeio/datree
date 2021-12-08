# Policy: labels_best_practices (example policy)
Kubernetes labels allow Engineers to perform in-cluster object searches, apply bulk configuration changes, and more. Labels can help simplify and solve many day-to-day challenges encountered in Kubernetes environments if they are set correctly.  

__This policy helps to enforce the following labels best practices:__
* [Ensure pre-defined environment labels are used](#ensure-pre-defined-environment-labels-are-used)
* [Ensure the owner label is used](#ensure-the-owner-label-is-used)
* Ensure all labels have a valid label value (RFC 1123)

## Ensure pre-defined environment labels are used
Having an env label is useful for performing bulk operations in specific environments or for filtering workloads according to their stage. This rules will also ensure that only pre-aprooved `environment` label values are used:
* `prod`
* `staging`
* `test`

### Failing this rule
If the `environment` key is missing from the labels section:  
```
kind: Deployment
metadata:
  labels:
    owner: yoda-at-datree.io
```

__OR__ a different `environment` value is used:
```
kind: Deployment
metadata:
  labels:
    environment: QA
```

## Ensure the owner label is used
An `owner` label is great not only for financial ownership but also useful for operational ownership. Consider adding an owner label to your workload with the name, email alias or Slack channel of the team responsible for the service. This will make it easier to alert the relevant team or team member when necessary.

### Failing this rule
If the `owner` key is missing from the labels section:  
```
kind: Deployment
metadata:
  labels:
    name: app-back
```

## Ensure all labels have a valid label value (RFC 1123)
Labels are nothing more than custom key-value pairs that are attached to objects and are used to describe and manage different Kubernetes resources. If the labels do not follow Kubernetes label syntax requirements (RFC-1123), they will not be applied properly.  

A lowercase RFC-1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name', or '123-abc', regex used for validation is 'a-z0-9?')  

### Failing this rule
If one of the labels values don't follow the RFC-1123 standart:  
```
kind: Deployment
metadata:
  labels:
    name: test@datree.io
```
