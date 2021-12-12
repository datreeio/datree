# Policy: horizontal_pod_autoscaling
Horizontal Pod Autoscaling allows automatic scaling of workload to match demand. This works by increasing or decreasing the number of Pods. This technique prevents unexpected cost explosions so you can focus on saving costs. With this configuration, Kubernetes engineers can work on actually reducing costs and scaling their pods rather than ensuring correct configuration. On top of that, this policy also verifies the usage of resources and resource metrics like cpu.

__This policy helps to enforce the following HPA workloads:__
* [Ensure target cpu utilization is set](#ensure-target-cpu-utilization-is-set)
* [Ensure max replicas is set and valid](#ensure-max-replicas-is-set-and-valid)
* [Ensure min replicas is set and valid](#ensure-min-replicas-is-set-and-valid)
* [Ensure scale target ref is configured properly](#ensure-scale-target-ref-is-configured-properly)

## Ensure target cpu utilization is set
The field `targetCPUUtilizationPercentage` defines the target for when the pods are to be scaled. CPU Utilization is the average CPU usage of all pods in a deployment divided by the requested CPU of the deployment. If the mean of CPU utilization is higher than the target, then the pod replicas will be readjusted.

### When this rule is failing?
If the `targetCPUUtilizationPercentage` key is missing from the labels section:  
```
kind: HorizontalPodAutoscaler
spec:
  maxReplicas: 10
  minReplicas: 1
```

__OR__ a value outside of the range, 1 - 100 was used:
```
kind: HorizontalPodAutoscaler
spec:
  targetCPUUtilizationPercentage: 200
```

## Ensure max replicas is set and valid

### When this rule is failing?
If the `owner` key is missing from the labels section:  
```
kind: HorizontalPodAutoscaler
metadata:
  labels:
    name: app-back
```

## Ensure min replicas is set and valid
The field `minReplicas` define the minimum number of replicas of a resource. As a best practice, it should be set to two, hence the minimum minReplicas one can set it 2.

### When this rule is failing?
If `minReplicas` is not present:
```
kind: HorizontalPodAutoscaler
spec:
  maxReplicas: 10
  targetCPUUtilizationPercentage: 50
```
## Policy author
Nishant Verma \\ [theowlet](https://github.com/theowlet)
