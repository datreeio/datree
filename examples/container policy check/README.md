# Policy: Containers_best_practices

Pods are the smallest units of deployable ,which you can create and manage in Kubernetes.

As Pod contain different Phase State i,e.. Pending,Running,Succeded,Failed and Unknown.

# What are Container states?

As Kubernetes tracks the state of each container inside a Pod.
A scheduler assigns a Pod that starts creating containers for 
that Pod using a container runtime. There are three possible container states: Waiting, Running, and Terminated.

It helps to easily identify the Containers state label and verify the correct state and have a valid state,
A failed or terminated state gives an error.

Ensure that each container has its correct state which are used.
Ensure that each containers has a configured with running status that indicates the container is executing without issues.
Ensure containers has valid its valid state.

# When this rule is failing?

*If the environment key is missing from the labels section or if a different environment value is used.

*This policy fails when it could not find appropriate container state label during container running.

# When this policy is Pass?

* This policy will pass when state includes information like the container’s sate label for each process.

# Ensure that each container has its correct state which are used (mentioned below):


* When it is in `Waiting` State':
If a container is not in either the Running or Terminated state, 
it is Waiting. A container in the Waiting state is still running the operations it requires in order to complete.

* When it is in `Running` State: 
The Running state indicates that a container is executing.

* When it is in `Terminated` State':
A container in the Terminated state began execution and then either ran to completion or failed for some reason.



# Ensure that each containers has a configured with running status that indicates the container is executing without issues.

defaultMessageOnFailure: Containers stop executing as running state failed.

# Ensure containers has valid its valid state.

defaultMessageOnFailure: Accept only approved container states (`waiting`, `running` and `terminated`)
