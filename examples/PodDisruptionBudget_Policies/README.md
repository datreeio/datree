# Policy:Pod_Disruption_Budget
A `PodDisruptionBudget` can be created for each Application.A `PodDisruptionBudget` limits the number of pods of a replicated application that are down simultaneously from `voluntary disruptions` like
* Deleting the deployment that manages the Pod
* Updating a deployment's pod template causing a Restart
* Deleting a Pod(Ex: By Accident)

__This policy helps to Check the following Pod_Disruption_Budget:__
* __Only one of the following Policies Should Be used__
* [Ensure Pod Disruption Budget with minAvailable Selector is Set](#Ensure-Pod-Disruption-Budget-with-minAvailable-Selector-is-Set)
* [Ensure Pod Disruption Budget with maxUnvailable Selector is Set](#Ensure-Pod-Disruption-Budget-with-maxUnvailable-Selector-is-Set)

## Ensure Pod Disruption Budget with minAvailable Selector is Set
Having a `PodDisruptionBudget` hepls K8's Admis to Set Minimum Number of Pods that should be Available When `voluntary disruptions by K8's Admins` Occur like 
* Draining a node for repair or upgrade
* Draining a node from a cluster to scale the cluster down
* Removing a pod from a node to permit something else to fit on the node

### When this rule is failing?
 If `minAvailable` Selector is not Set to an Integer or Percentage
 ```
 apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdbdemo
spec:
  minAvailable: NA
 ```
 __OR__ a different `PDBSelector` is used:
 ```
 apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdbdemo
spec:
  maxUnavailable: 1
```
## Ensure Pod Disruption Budget with maxUnavailable Selector is Set
Having a `PodDisruptionBudget` hepls K8's Admis to Set Maximum Number of Pods that should be Unavailable When `voluntary disruptions by K8's Admins` Occur like 
* Draining a node for repair or upgrade
* Draining a node from a cluster to scale the cluster down
* Removing a pod from a node to permit something else to fit on the node

### When this rule is failing?
 If `maxUnvailable` Selector is not Set to an Integer or Percentage
 ```
 apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdbdemo
spec:
  maxUnvailable: NA
 ```
 __OR__ a different `PDBSelector` is used:
 ```
 apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdbdemo
spec:
  minAvailable : 2
  ```
## Policy author
Avinash \\ [Avinash](https://github.com/Avinashnayak27)




