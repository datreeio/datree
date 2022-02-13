# 2. Implement local evaluation

Date: 2022-02-13

## Status

Accepted

## Context

We want to allow our users to run datree locally without the need of sending the data to our backend
Users want some more security features - not send data to datree


## Decision

We will make the cli should work without backend datree / without internet
We will implement local evaluation:
# We will add the ability to run datree offline
# We will support --no-record flag

## Consequences

Custom rules is not executed using golang json schema package instead of npm package. This might affect potential users
Same goes for yamlschemavalidator.datree.io that might give bad result to a potential user
