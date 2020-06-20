#!/bin/bash

echo 'Build go project'
go build -o kubectl-ingress_refresh

echo 'Remove old'
rm -rf /usr/local/bin/kubectl-ingress_refresh

echo 'Install new'
cp kubectl-ingress_refresh  /usr/local/bin