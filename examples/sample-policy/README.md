## Policy #1 - startingDeadlineSeconds_min_accepted_value
> Ensure min value is set for `startingDeadlineSeconds`
---
If `startingDeadlineSeconds` is set to a value less than 10 seconds, the CronJob may not be scheduled. This is because the CronJob controller checks things every 10 seconds.

### When this rule is failing?
If the `startingDeadlineSeconds` is set to less than 10 seconds
```
kind: CronJob
metadata:
  name: omiCronJob
spec:
  schedule: "*/1 * * * *"
  startingDeadlineSeconds: 8
```

## Policy #2 - ttlSecondsAfterFinished_set
> Ensure value for `ttlSecondsAfterFinished` is set
---
Add `ttlSecondsAfterFinished` to ensure unmanaged jobs are not left around after Job is fully deleted

### When this rule is failing?
If the `ttlSecondsAfterFinished` is not present in the job manifest *.spec*
```
kind: Job
metadata:
  name: job-with-ttl
spec:
  template:
    spec:
      containers:
      - name: myjob
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
```

## Policy #3 - revisionHistoryLimit_min_accepted_value
> Ensure `revisionHistoryLimit` is not explicitly set to zero
---
Specifies how many old ReplicaSets for this Deployment you want to retain. Explicitly setting `revisionHistoryLimit` to 0, will result in cleaning up all the history of your Deployment preventing you from rolling back.

### When this rule is failing?
If the `revisionHistoryLimit` is set to 0
```
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  revisionHistoryLimit: 0
  replicas: 3
  selector:
    matchLabels:
      app: nginx
```

## Policy #4 - namespace_present
> Ensure workload has a configured `namespace`
---
Add a proper `namespace` from accidently creating pods in the active namespace

### When this rule is failing?
If `namespace` is not present under *.metadata*
```
kind: Pod
metadata:
  name: mypod
  labels:
    name: mypod
spec:
  containers:
  - name: mypod
    image: nginx
```

## Policy #5 - poddisruptionbudget_minAvailable_set
> Ensure `minAvailable` is set to protect from voluntary disruption
---
Add `minAvailable` for a PodDisruptionBudget to ensure that the number of replicas running is never brought below the number needed

### When this rule is failing?
If `minAvailable` is not present under *.spec*
```
kind: PodDisruptionBudget
metadata:
  name: pdb
spec:
  selector:
    matchLabels:
      app: zookeeper
```

## Policy #6 - startup_probe_configured
> Ensure `startupProbe` options is used to determine app start state within the container
---
If you want to wait before executing a liveness probe you should use initialDelaySeconds or a `startupProbe`

### When this rule is failing?
If `startupProbe` is not present under *.spec.containers*
```
kind: Pod
metadata:
  name: etcd-pod
spec:
  containers:
  - name: etcd
    image: k8s.gcr.io/etcd:3.5.1-0
    command: [ "/usr/local/bin/etcd", "--data-dir",  "/var/lib/etcd", "--listen-client-urls", "http://0.0.0.0:2379", "--advertise-client-urls", "http://127.0.0.1:2379", "--log-level", "debug"]
    ports:
    - name: liveness-port
      containerPort: 8080
      hostPort: 8080
    livenessProbe:
      httpGet:
        path: /healthz
        port: liveness-port
      failureThreshold: 1
      periodSeconds: 10
```

## Policy #7 - image_avoid_latest_tag
> Ensure images do not use the `:latest` tag
---
Avoid using the `:latest` tag for images as its difficult to track and rollback properly. More of a production use-case.

### When this rule is failing?
If `:latest` tag is present under *.spec.containers.image* value
```
Kind: Pod
metadata:
  name: fail-policy
  labels:
    owner: Runcy
    environment: prod
    app: web
spec:
  containers:
  - name: nginx
    image: nginx:latest
    imagePullPolicy: Sometimes
    ports:
    - containerPort: 80
```

## Policy #8 - ipfamilypolicy_accepted_options
> Ensure correct `ipFamilyPolicy` options are used
---
The below values will help ensure to determine support for IPv4 and/or IPv6 addresses
* `SingleStack`
* `PreferDualStack`

### When this rule is failing?
If invalid `ipFamilyPolicy` option is provided under *.spec.ipFamilyPolicy*
```
kind: Service
metadata:
  name: my-service
  labels:
    app: MyApp
spec:
  ipFamilyPolicy: DualStack
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
```

## Policy #9 - concurrencyPolicy_accepted_options
> Ensure correct `concurrencyPolicy` options are used
---
The below values will help ensure to determine `concurrencyPolicy` values while scheduling a CronJob
* `Allow`
* `Forbid`
* `Replace`

### When this rule is failing?
If invalid `concurrencyPolicy` option is provided under *.spec.concurrencyPolicy*
```
kind: CronJob
metadata:
  name: OmiCronJob
spec:
  schedule: "*/1 * * * *"
  concurrencyPolicy: Deny
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello world
          restartPolicy: OnFailure
```

## Policy #10 - limitrangetype_only_container
> Options used for LimitRange kind are for Container only
---
Enforce `type` values set for LimitRange only allows the below value
* `Container`

### When this rule is failing?
If `type` option is set for Pod under *.spec.concurrencyPolicy*
```
kind: LimitRange
metadata:
  name: mem-limit-range
spec:
  limits:
  - default:
      memory: 512Mi
    defaultRequest:
      memory: 256Mi
    type: Pod
```

## Policy #11 - deploymentStrategy_only_recreate
> Options used for Deployment kind only allows Recreate
---
Enforce `type` values set for *.spec.strategy* only allows the below value
* `Recreate`

### When this rule is failing?
If `type` option is set for RollingUpdate under *.spec.strategy*
```
kind: Deployment
metadata:
  name: rollingupdate-strategy
spec:
  strategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: web-app-rollingupdate-strategy
      version: nanoserver-1809
```

## Policy #12 - imagePullPolicy_accepted_options
> Ensure correct `imagePullPolicy` options are used
---
The below values will help ensure to determine `imagePullPolicy` values for a Pod
* `Always`
* `IfNotPresent`
* `Never`

### When this rule is failing?
If invalid `imagePullPolicy` option is provided under *.spec.containers*
```
Kind: Pod
metadata:
  name: fail-policy
  labels:
    owner: Runcy
    environment: prod
    app: web
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    imagePullPolicy: Sometimes
    ports:
    - containerPort: 80
```

## Policy #13 - imagePullSecrets_contains_name
> Ensure secrets `name` is present
---
Add proper secrets `name` to make sure the image pull happens without failure while using private registry

### When this rule is failing?
If `name` option is missing under *spec.imagePullSecrets*
```
Kind: Pod
metadata:
  name: fail-policy
  labels:
    owner: Runcy
    environment: prod
    app: web
spec:
  containers:
    - name: run
      image: roommen/my-awesome-app:v1
  imagePullSecrets:
```

## Policy #14 - seLinuxOptions_contains_level
> Ensure valid `level` is assigned
---
Add proper SELinux `level` for your container while using securityContext

### When this rule is failing?
If `level` option is missing under *spec.securityContext.seLinuxOptions*
```
kind: Pod
metadata:
  name: security-context
spec:
  containers:
  - name: sec-ctx
    image: gcr.io/google-samples/node-hello:1.0
    securityContext:
      seLinuxOptions:
        role: "admin"
```

## Policy #15 - stringdata_present
> Ensure only `data` is provided
---
Secreated should only be created with `data` that accepts base64 encoded format. It should not have *stringData* that accepts plain-text values
### When this rule is failing?
If *stringData* option is present for Secret kind
```
kind: Secret
metadata:
  name: secret-basic-auth
type: kubernetes.io/basic-auth
stringData:
  username: admin
  password: t0p-Secret
```


# Policy authored by
Runcy Oommen \\ [https://runcy.me](https://runcy.me)
