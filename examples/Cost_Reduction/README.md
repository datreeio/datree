# Recommended policies for production workloads
It is recommended for the kubernetes workloads deployed in production to support the industry best practices and have hardened rules enforced on them.

__The following policies are named as per their respective use cases:__
* Cost Reduction

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

## Policy author
Dhanush M / [dhanushm7](https://github.com/dhanushm7)
