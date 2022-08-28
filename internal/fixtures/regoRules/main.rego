package main
import compute

deny[msg] {
    input.kind != "Deployment"
    msg := { "message": "this rule failed, here is why 1", "ruleID": "THIS_IS_THE_RULE_ID_1" }
}

deny[msg] {
    compute.isKindDeployment
    msg := { "message": "this rule failed, here is why 2", "ruleID": "THIS_IS_THE_RULE_ID_2" }
}
