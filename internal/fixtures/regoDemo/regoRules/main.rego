package main
import data.compute

deny[error] {
    input.kind == "Deployment"
    error := { "ruleID": "REGO_RULE_1", "message": "this rule failed, here is why 1" }
}

deny[error] {
    compute.isKindDeployment
    error := { "ruleID": "REGO_RULE_1", "message": "this rule failed, here is why 2" }
}

deny[error] {
    compute.isKindDeployment
    error := { "ruleID": "REGO_RULE_3" }
}
