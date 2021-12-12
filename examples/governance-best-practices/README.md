# Policy: stability_best_practices
Ingress host name should be a valid organization name. Invalid host names can cause problems with access to the API.

__This policy helps to enforce the following labels best practices:__
* [Ensure Ingress only uses approved domain names for hostnames](#ensure-ingress-only-uses-approved-domain-names-for-hostnames)

## Ensure Ingress only uses approved domain names for hostnames
### When this rule is failing?
If `host` is missing:
```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: fail-governance-best-practices-for-ingress
spec:
  rules:
    - http:
        paths:
          - pathType: ImplementationSpecific
            path: /
            backend:
              service:
                name: nginx
                port:
                  number: 80
```

__OR__ the value of `host` is not a valid organization name (*.example.com):

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pass-governance-best-practices-for-ingress
spec:
  rules:
    - host: test.com
      http:
        paths:
          - pathType: ImplementationSpecific
            path: /
            backend:
              service:
                name: nginx
                port:
                  number: 80

```

## Policy author
Brijesh Shah \\ [brijeshshah13](https://github.com/brijeshshah13)