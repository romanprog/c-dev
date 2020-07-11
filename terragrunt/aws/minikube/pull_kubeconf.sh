#!/bin/bash

set -x
CLUSTER_FULLNAME=$1
WAIT_TIME=5

export KUBECONFIG=~/.kube/kubeconfig_${CLUSTER_FULLNAME}

until kubectl version --request-timeout=5s >/dev/null 2>&1; do
    sleep ${WAIT_TIME}
    aws s3 cp "s3://${CLUSTER_FULLNAME}/kubeconfig_${CLUSTER_FULLNAME}" "${HOME}/.kube/kubeconfig_${CLUSTER_FULLNAME}" 2>/dev/null
    cp ${HOME}/.kube/kubeconfig_${CLUSTER_FULLNAME} ${HOME}/.kube/config 2>/dev/null
done
