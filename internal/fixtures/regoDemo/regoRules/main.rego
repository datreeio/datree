package main
import data.compute

deny[error] {
    input.kind == "Deployment"
    error := { "ruleID": "THIS_IS_THE_RULE_ID_1", "message": "this rule failed, here is why 1" }
}

deny[error] {
    compute.isKindDeployment
    error := { "ruleID": "THIS_IS_THE_RULE_ID_1", "message": "this rule failed, here is why 2" }
}

deny[error] {
    compute.isKindDeployment
    error := { "ruleID": "THIS_IS_THE_RULE_ID_3", "message": "this rule failed, here is why 3" }
}
