# Policy: industry_best_practices
Kuberenets configuration files are made out of a specified structure. Although this structure can be tweeked by the developers according to their usecase. Kuberenetes labels helps engineers a lot to structure these configuration file into easy and human readable format.

__This policy helps to enforce the following labels best practices:__
* Ensure each configuration file has deployment name.

## Ensure each configuration file has apiVersion label
While starting the configuration it's essential to make sure that there is `name` specified for the deployment. Not having this might not raise errors by the native development kits (kubectl etc) but naming the deployment configurations while working with several clusters in a team is essential and cosidered a good practice. 
* `name: demo-application-deployment`

### When this rule is failing?
If the `name` tag is not included in the file (Generally under `metadata`)
```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    owner: hello
```

### When this rule pass?
The policy passes when provided with `name` of the deployment.
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-application-deployment
  labels:
    owner: hello
```

```
## Policy author
Kaiwalya Koparkar \\ [kaiwalyakoparkar](https://github.com/kaiwalyakoparkar)
Ashwin Kumar Uppala \\ [ashwinexe](https://github.com/ashwinexe)
Karuna Tata \\ [starlightknown](https://github.com/starlightknown)
Abhishek Choudhary \\ [shreemaan-abhishek](https://github.com/shreemaan-abhishek)