
# Custom Policies To Validate Kubernetes Mainfests Using Datree CLI

Kubernetes labels enable engineers to perform in-cluster object searches, apply bulk configuration changes, and more. Labels can help simplify and solve many day-to-day challenges encountered in Kubernetes environments if they are set correctly.

__This custom policies helps to enforce the following labels best practices for use cases:__
* [Ensure strategy has pre-defined labels ](#ensure-strategy-has-pre-defined-labels)
* [Ensure RollingUpdate strategy has maxSurge and maxUnavailable labels](#ensure-RollingUpdate-strategy-has-maxSurge-and-maxUnavailable-labels)
* [Ensure pre-defined DnsPolicy values are used for pods](#ensure-pre-defined-DnsPolicy-values-are-used-for-pods)
* [Ensure custom nodeselector has pre-defined disktype value](#ensure-custom-nodeselector-has-pre-defined-disktype-value)


## Ensure strategy has pre-defined labels 
Kubernetes deployment strategies are use to replace existing pods with new pods in two ways.
* `Recreate` kill all existing pods before creating new ones.
* `RollingUpdate` replace the old ReplicaSets by new one using rolling update i.e gradually scale down the old ReplicaSets and scale up the new one.

This rule will also ensure that only pre-approved `strategy.type` label values are used
### When this rule is failing?
If the `strategy` key is missing from the spec section:  
```
kind: Deployment
metadata:
spec:
  replicas: 2
  template:
```

__OR__ a different `strategy.type` value is used:
```
kind: Deployment
spec:
  template:
  strategy:
    type: Ab
```


## Ensure RollingUpdate strategy has maxSurge and maxUnavailable labels
RollingUpdate strategy replace the old ReplicaSets by new one using rolling update i.e gradually scale down the old ReplicaSets and scale up the new one.
Rolling update config params. Present only if `DeploymentStrategyType = RollingUpdate`. This rule will ensure that `rollingupdate` labels are numeric
* `maxSurge`
* `maxUnavailable`  
params control the desired behavior of rolling update.
### When this rule is failing?
while `DeploymentStrategyType` = `RollingUpdate` but `maxSurge` and `maxUnavailable`  are not defined in spec.strategy.rollingUpdate section
```
kind: Deployment
metadata:
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
```
__OR__ assigned non-numeric labels

```
kind: Deployment
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: ab
      maxUnavailable: aa
```


## Ensure pre-defined DnsPolicy values are used for pods
DNS policies can be set on a per-pod basis. Currently Kubernetes supports the following pod-specific DNS policies. 
* `Default`- The Pod inherits the name resolution configuration from the node that the pods run on
* `ClusterFirst` - Any DNS query that does not match the configured cluster domain suffix, such as "www.kubernetes.io", is forwarded to the upstream nameserver inherited      from the node
* `ClusterFirstWithHostNet` - For Pods running with hostNetwork, you should explicitly set its DNS policy "ClusterFirstWithHostNet"
* `None` -  It allows a Pod to ignore DNS settings from the Kubernetes environment

These policies are specified in the `dnsPolicy` field of a PodSpec.

### When this rule is failing?
If the `dndPolicy` key is missing from the template.spec section:
```
kind: Deployment
spec:
  replicas: 2
  template:
     metadata:
     spec:
       container:
       
```
__OR__ a different `dnsPolicy` value is used:
```
kind: Deployment
spec:
  replicas: 2
  template:
     metadata:
     spec:
       container:
       dnsPolicy: ab
```

## Ensure custom nodeselector has pre-defined disktype value
`nodeSelector` is the simplest recommended form of node selection constraint. nodeSelector is a field of PodSpec. It specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). The most common usage is one key-value pair.

`nodeSelector:
   disktype: ssd`   
### When this rule is failing?
If the `nodeSelector` key is missing from the template.spec section:
```
kind: Deployment
spec:
  replicas: 2
  template:
     metadata:
     spec:
       container:
```
__OR__ `nodeSelector.disktype` value is not equal to `ssd`:
```
kind: Deployment
spec:
  replicas: 2
  template:
     metadata:
     spec:
       container:
       nodeSelector:
          disktype: hdd
```




## Policy author
Kapakayala Naga sai krishna vinay swami (aka siddhu)\\ [siddhusniper](https://github.com/siddhusniper)

