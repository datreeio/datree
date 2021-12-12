# Policy: stability_best_practices
Kubernetes podAntiAffinity is a feature that helps you avoid placing pods on the same node to prevent downtime.

__This policy helps to enforce the following labels best practices:__
* [Prevent containers from running on the same node if multiple replicas are specified](#prevent-containers-from-running-on-the-same-node-if-multiple-replicas-are-specified)

## Prevent containers from running on the same node if multiple replicas are specified
Inter-pod anti-affinity allow you to constrain which nodes your pod is eligible to be scheduled based on labels on pods that are already running on the node rather than based on labels on nodes. Refer to [documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity) for more details.

### When this rule is failing?
If `podAntiAffinity` is missing when multiple replicas are specified:
```
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```
## Policy author
Brijesh Shah \\ [brijeshshah13](https://github.com/brijeshshah13)