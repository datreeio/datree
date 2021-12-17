# Policy: governance_best_practices
Ingress host name should be a valid organization name. Invalid host names can cause problems with access to the API.

__This policy helps to enforce the following best practices:__
* [Ensure Ingress only uses approved domain names for hostnames](#ensure-ingress-only-uses-approved-domain-names-for-hostnames)

## Ensure Ingress only uses approved domain names for hostnames
### When this rule is failing?
If `host` is missing:
```
apiVersion: networking.k8s.io/v1
kind: Ingress
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

`--ignore-missing-schemas` will have to be used to ignore the missing schema.

## Policy author
Brijesh Shah \\ [brijeshshah13](https://github.com/brijeshshah13)