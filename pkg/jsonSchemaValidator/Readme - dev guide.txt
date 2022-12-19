Creating custom keys allows us to create more intricate rules involving multiple fields in the same manifest. 
For example, we can create a rule that requires a field to be present if another field is present. 
This is supported in the "github.com/santhosh-tekuri/jsonschema/v5" library.


How to create a custom key
--------------------------
1. Create a new .go file inside the "extensions" directory (copy-paste the "requestLimitEq.go" file and rename it)
2. 

How does it work?
-----------------
When adding a custom key to a schema, 