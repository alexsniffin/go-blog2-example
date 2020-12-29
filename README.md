# Example Prometheus Evaluation Service

This is an example from the article I wrote over on [medium](https://alexsniffin.medium.com/building-a-prometheus-expression-evaluation-service-in-go-ea58f0cc406).

## What It Does

A configurable service which takes a rule set for evaluating Prometheus query responses against expressions and then templating the expression response into a message to send to Slack.