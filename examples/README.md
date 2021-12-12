# Recommended policies for production workloads
It is recommended for the kubernetes workloads deployed in production to support the industry best practices and have hardened rules enforced on them.

__The following policies are named as per their respectives use cases:__
* Governance
* Cost Reduction
* Security

## Governance

This policy consists of three custom rules
1. `PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION`
2. `WORKLOAD_MISSING_NAMESPACE_FOR_NAMESPACED_RESOURCES`
3. `INGRESS_INCORRECT_DOMAIN_NAME`

### PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION
This rule ensures that the our highly available worloads used in production are not denied any voluntary disruption. \
If you set `maxUnavailable` to 0% or 0, or you set `minAvailable` to 100% or the number of replicas, you are requiring zero voluntary evictions. When you set zero voluntary evictions for a workload object such as ReplicaSet, then you cannot successfully drain a Node running one of those Pods. If you try to drain a Node where an unevictable Pod is running, the drain never completes. \
ref: https://kubernetes.io/docs/tasks/run-application/configure-pdb/

### WORKLOAD_MISSING_NAMESPACE_FOR_NAMESPACED_RESOURCES
This rule ensures that all the namespaced resources are namescoped as it prevents kubernetes users from accidently deploying workloads into default namespace when they don't specify namespace.

### INGRESS_INCORRECT_DOMAIN_NAME
This rule ensures that all the host domain names specified in the Ingress resource are valid. By default the kubernetes api does not disallow the creation of invalid hostnames.

## Cost Reduction

### SERVICE_DENY_TYPE_LOADBALANCER
This rule ensures that the kubernetes applications are not exposed using service type `LoadBalancer` as an additional load balancer is created each time when  any new application is exposed through this service type. Instead it is recommended to the clusterIP or Ingress to expose the same set of services without undergoing a cost overhead.

## Security

### WORKLOAD_INVALID_CONTAINER_CAPABILITIES_VALUE
This rule ensures that no container is allowed to have the insecure/dangerous linux capabilities which causes the applications in the cluster to be vulnerable to security attacks.

## Policy author
Dhanush M \\ [Dhanush M](https://github.com/dhanushm7)
