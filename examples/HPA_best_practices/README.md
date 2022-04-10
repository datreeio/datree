# Policy: HPA_best_practices 
When auto-scaling resource utilization is triggered with HPA (HorizontalPodAutoscaler), a `targetCPUUtilizationPercentage` should be defined in Manifest File to get a Threshold CPU Limit and it will Scale a Pod and Create a New Replica
## Ensure targetCPUUtilizationPercentage is Set
 Having an `targetCPUUtilizationPercentage` is Useful to create Replicas Based on CPU_Utilization
### When this rule is failing? 
If `targetCPUUtilizationPercentage` is not found in Manifest file of `Kind: HorizontalPodAutoscaler`
```
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  minReplicas: 1
  maxReplicas: 10
```
## Policy author
Avinash \\ [Avinash](https://github.com/Avinashnayak27)