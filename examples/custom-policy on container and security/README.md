# Custom Policy for Cloud Native Hackathon 2021

I have created the 4 Custom Policy for Kubernetes clusters containers
the list as follows:

1. [Ensure existance for namespace](#)
2. [Check existance for more container than expected](#)
3. [Ensure existance for exposed port of container](#)
4. [Ensure existance for correct OS or image of the container](#)

# Documentation

## 1. Ensure the existance of namespace(Imposatant security policy)

### 1. Use case scenario of namespace

You can think of a Namespace as a virtual cluster inside your Kubernetes cluster.
You can have multiple namespaces inside a single Kubernetes cluster,
and they are all logically isolated from each other.
They can help you and your teams with organization, security,
and even performance!

#### The Most important use case using namespaces is that whenever there is on the cluster it helps to restrict the attack at very small so tha it protect the other cluster and container to remain unaffected.

### 2. When the policy failing?

When there is absences of the of the namespace in the cluster

```bash
apiVersion: v1
kind: Pod
metadata:
  name: fail-policy
  labels:
    name: mypod
spec:
  containers:
    - name: mypod
      image: nginx
```

As in the above yaml file there not existance of namespace

### 3. When the policy Passing

When there is presences of the of namespace in the cluster's yaml file

```bash
apiVersion: v1
kind: Pod
metadata:
  name: pass-policy
  namespace: test
  labels:
    name: mypod
spec:
  containers:
    - name: mypod
      image: nginx

```

As in the above yaml file there presences of namespace in cluster's yaml file.

## 2. Check existance for more container than expected (Important for System Stablity)

### 1. Use case scenario of policy

Consider, the scenario is there is an bug in the code which
creates the unnessarily multiple container in the cluster.
Due to this it eventually use more resourse of the cluster and
after that if cluster hits the resourse limit it will stop working
So it necessary to have policy to restrict the number of container than expected, it will increase stablity of system.

### 2. When the policy failing?

When there is presences of more containers than expected

```bash
kind: Pod
apiVersion: v1
metadata:
  name: multi-container
spec:
  containers:
    - name: container-1
      image: nginx
    - name: container-2
      image: ubuntu

```

In this perticular case i have restrict the cluster to only one container
but in the policy i have created we can change the number of container as we want

### 2. When the policy Passing

In this perticular case it has only one container and i had set to one only so it will pass.

```bash
kind: Pod
apiVersion: v1
metadata:
  name: multi-container
spec:
  containers:
    - name: container-1
      image: nginx

```

## 3. Ensure existance for exposed port of container

### 1. Use case scenario of policy

When in single Kubernetes cluster if there is multi-container
there should be and way so that each container can communicate with each other
and exchange data so it is necessary to have the port exposed to communicate

So i have created the custom policy to check whether in Kubernetes cluster, the container port are exposed or not

### 2. When the policy failing?

Whenever in Kubernetes cluster if container port is not exposed the policy will fail.

```bash
apiVersion: v1
metadata:
  name: fail-policy
spec:
  containers:
    - name: container-exposed-port
      image: nginx
```

### 2. When the policy Passing

Whenever in Kubernetes cluster if container port
is exposed the policy will pass.

Look here in this perticular case in this yaml file the container
port are exposed

```bash
kind: Pod
apiVersion: v1
metadata:
  name: pod-exposed-port
spec:
  containers:
    - name: container-exposed-port
      image: nginx
      ports:
        - containerPort: 80

```

## 4. Ensure existance for correct OS or image of the container

Whenever in Kubernetes cluster if the code is code is written in perticularily
taken into consideration of operating system they working but while at the deloyment
due to some reason the container's operating system or image get change there it will create misconfigration
while deloyment

So it overcome this problem i have created custom policy to ensure the correct image of the Kubernetes cluster's container.

### 2. When the policy failing?

Whenever in Kubernetes cluster's container image is not what it expected to be it will fail.

Here in this perticular case the expected image of container
must be ubuntu but it not so it will fail.

```bash
kind: Pod
apiVersion: v1
metadata:
  name: pass-image-container
spec:
  containers:
    - name: container-1
      image: nginx

```

### 2. When the policy Passing

Whenever in Kubernetes cluster's container image meet criteria of image it will pass.

In this perticular case we expected to Kubernetes cluster's container image
to be ubuntu and it is there so it will pass.

```bash
kind: Pod
apiVersion: v1
metadata:
  name: pass-image-container
spec:
  containers:
    - name: container-1
      image: ubuntu

```

## Lessons Learned

It was very Wonderful experiance while build these custom policy
Initially for it was very challenging task for me
because begin an beginner. But eventually as go deeper and
deeper in Learning yaml syntax and more about cloud technology
and different use case. I Learned new thing such as what is
kubernetes,docker,cloud and cluster and also about git and
github. while Learning new syntax i have done many mistakes and gone through
n-th number of problem but very thanks to the team of the datree
to help me whem i was stuck. They guided very much they gave me various references from
where i can learn and solve the problem that i encontered.

Inspit of the all this i finally able to manage to create 4 custom policy.
Thanks You very much for this Wonderful oppurtunity.

## 🚀 About Me

- I'm a Beginner Front End Developer and Open Source Enthusiast

## Authors

- [@AdeshKhandait](https://github.com/AdeshKhandait)
