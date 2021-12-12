# Recommended policies for production workloads
It is recommended for the kubernetes workloads deployed in production to support the industry best practices and have hardened rules enforced on them.

__The following policies are named as per their respective use cases:__
* Governance
* Cost Reduction
* Security

## Governance

This policy consists of three custom rules
1. `PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION`
2. `WORKLOAD_MISSING_NAMESPACE_FOR_NAMESPACED_RESOURCES`
3. `INGRESS_INCORRECT_DOMAIN_NAME`

#### PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION
This rule ensures that the our highly available worloads used in production are not denied any voluntary disruption. \
If you set `maxUnavailable` to 0% or 0, or you set `minAvailable` to 100% or the number of replicas, you are requiring zero voluntary evictions. When you set zero voluntary evictions for a workload object such as ReplicaSet, then you cannot successfully drain a Node running one of those Pods. If you try to drain a Node where an unevictable Pod is running, the drain never completes. \
ref: https://kubernetes.io/docs/tasks/run-application/configure-pdb/

##### When this rule is failing?
When any PodDisruptionBudget resource has `maxUnavailable` set to **0** or **0%** or `minAvailable` is set to **100%**
```
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: zk-pdb
  namespace: zookeeper
spec:
  maxUnavailable: 0%
  selector:
    matchLabels:
      app: zookeeper
---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: ms-pdb
  namespace: mesos
spec:
  maxUnavailable: 0
  selector:
    matchLabels:
      app: mesos
---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: ka-pdb
  namespace: kafka
spec:
  minAvailable: 100%
  selector:
    matchLabels:
      app: kafka
```

#### WORKLOAD_MISSING_NAMESPACE_FOR_NAMESPACED_RESOURCES
This rule ensures that all the namespace scoped resources are namescoped as it prevents kubernetes users from accidently deploying workloads into default namespace when they don't specify namespace.

##### When this rule is failing?
When any namespace scoped resource does not have any namespace specified.

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

#### INGRESS_INCORRECT_DOMAIN_NAME
This rule ensures that all the host domain names specified in the Ingress resource are valid. By default the kubernetes api does not disallow the creation of invalid hostnames.


##### When this rule is failing?
When any `Ingress` resource does not have domain name in the valid format.

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-example-ingress
  namespace: test
spec:
  rules:
    - host: datree.io-com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: service1
                port:
                  number: 80
```

## Cost Reduction

#### SERVICE_DENY_TYPE_LOADBALANCER
This rule ensures that the kubernetes applications are not exposed using service type `LoadBalancer` as an additional load balancer is created each time when  any new application is exposed through this service type. Instead it is recommended to the clusterIP or Ingress to expose the same set of services without undergoing a cost overhead.

##### When this rule is failing?
When any `Service` resource has type **LoadBalancer** specified.

```
apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: test
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
  type: LoadBalancer

```
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
