# Policy: labels_best_practices (example policy)
Kubernetes labels enable engineers to perform in-cluster object searches, apply bulk configuration changes, and more. Labels can help simplify and solve many day-to-day challenges encountered in Kubernetes environments if they are set correctly.  

# Version label
This label can be used to provide a version label so that for example if more than one people are trying to use the same container of kubernetes and both of them try to use same version then while commiting to the original project then there can be a misconfiguration.Like starting with 
* `1.0`
And if more than one people are using the same configuration they can make that as 1.0 - SNAPSHOT versions


# contains Version label
This label is useful for example if there is any version with same name already exists in the kubernetes contaner. And if it exists we can tell them there already exists please change the versioning number from that
* `true`
* `false`