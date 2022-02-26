#!/bin/bash

KO_DOCKER_REPO=kind.local ko resolve -f controller.yaml | kubectl apply -f - --prune -l app=generated-secrets