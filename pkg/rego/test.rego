package example.rules

allow {
    input.body.kind == "Deployment"
}
