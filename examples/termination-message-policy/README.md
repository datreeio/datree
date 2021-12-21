# Policy: deployment_best_practices (example policy)

Kubernetes retrieves termination messages from the termination message file specified in the terminationMessagePath field of a container, which has a default value of /dev/termination-log. Users can customize this field to tell kubernetes use a different file. Kubernetes uses the contents from the specified file to populate the Container's status message on both success and failure. Users cannot set the termination message path after the Pod is lauched.

Moreover, users can set the terminationMessagePolicy filed of a Container for further customization. This field defaults to File which means the termination messages are retrieved only from the termination message file. By setting the terminationMessagePolicy to FallbackToLogsOnError, you can tell kubernetes to use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller.

**This policy helps to enforce the following deployment best practices:**

- [Ensure termination message fields are provided for k8s to set container status on fail/pass scenarios](#Ensure-termination-message-fields-are-provided-for-k8s=-to-set-container-status-on-fail/pass-scenarios)
- [Ensure termination policy is provided to use last chunk of container logs writter under terminationMessagePath](#Ensure-termination-policy-is-provided-to-use-last-chunk-of-container-logs-writter-under-terminationMessagePath)

## Ensure terminationMessagePath and path follows the linux path conventions

Having an terminationMessagePath tells kubernetes to use a different file. Kubernetes uses the contents from the specified file to populate the Container's status message on both success and failure.

- `Provided Path tells kubernetes to use the specified path for logging container logs`

### When this rule is failing?

If the `terminationMessagePath` key is missing from the container spec section:

```
apiVersion: v1
kind: Pod
metadata:
  name: termination-msg
spec:
  containers:
  - name: termination-msg-example
    image: termination-msg-example
```

**OR** a value of `terminationMessagePath` does not follow linux conventions:

```
apiVersion: v1
kind: Pod
metadata:
  name: termination-msg
spec:
  containers:
  - name: termination-msg-example
    image: termination-msg-example
    terminationMessagePath: "C:\Users\root"
```

## Ensure terminationMessagePolicy is defined under container spec

Having terminationMessagePolicy tells kubernetes to use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller.

### When this rule is failing?

If the `terminationMessagePolicy` key is missing from the container spec section:

```
apiVersion: v1
kind: Pod
metadata:
  name: termination-msg
spec:
  containers:
  - name: termination-msg-example
    image: termination-msg-example
    terminationMessagePath: "/tmp/my-log"
```

## Policy author

Meghal Chhabria \\ [emjay010](https://github.com/emjay010)
