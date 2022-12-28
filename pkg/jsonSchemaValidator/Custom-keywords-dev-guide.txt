Creating custom keys allows us to create more intricate rules involving multiple fields in the same manifest. 
For example, we can create a rule that requires a field to be present if another field is present. 
This is supported in the "github.com/santhosh-tekuri/jsonschema/v5" library.


How to create a custom key
--------------------------
1. Create a new .go file inside the "extensions" directory with the name "customKeyRule<number>" (copy-paste the "customKeyRule81.go" file and rename it)
2. For each key, create a compiled string of its schema using the "MustCompileString" method. Make sure to set the key's type to "string".
3. For each key, implement the "Compile" and "Validate" methods.
4. Create structs as needed, containing the expected keys below your custom key.

How does it work?
-----------------
When testing a schema with a custom key against a manifest, the "Validate" function's "dataValue" parameter will contain the manifest's data starting from the custom key all the way down to the end.

Take rule 81 for example. Say we have the following manifest:

apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: front-end
          image: nginx:latest
          resources:
            requests:
              cpu: "64m"
            limits:
              cpu: "64m"


And we use the following schema:

properties:
    spec:
        properties:
            containers:
                items:
                properties:
                    resources:
                    customKeyRule81:
                        type: string

In this case the "dataValue" parameter will contain everything below the "resources" key - the "requests" and "limits" and their children.