if deny is empty
then
pass
else
fail, take message from deny if there are any. (fallback to failSuggestion)

ignore rules that are not enabled

```yaml
regoRules: ./path-to-rego-rules/**.rego
policies:
  - policy1
    rules:
    ID: rule_id_1
    isEnabled: false
```

```rego
package main

deny[res] {
    input.kind == "Deployment"
    res := {"message":"bad${index}", "ruleID": "MY_SPECIAL_RULE_ID"}
}
```

compile the rego upon publish (maybe)

fail for invalid deny struct upon run

"data.main.deny"

group deny messages by RuleID, concat messages
