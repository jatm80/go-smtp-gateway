#!/bin/bash

set -e

## Project configuration

NAME="go-smtp-gateway"
NAMESPACE="home-server"

## Terraform configuration

eval `pass JT-01/exports/homelab`
BUCKET="home-server-tf-backend"
KEY="app/homeserver/${NAME}.tfstate"
TF="terraform -chdir=infra"
#export TF_LOG_PROVIDER=TRACE

if [[ "$NAME" == "example" ]]; then
 echo "NAME and NAMESPACE must be defined before continue"
 exit 255
fi

${TF} init -backend=false -upgrade
${TF} validate
${TF} fmt
${TF} init -backend-config="bucket=${BUCKET}" \
  -backend-config="key=${KEY}" \
  -backend-config="encrypt=true"
${TF} plan -var="name=${NAME}" \
    -var="namespace=${NAMESPACE}"

PS3='Please enter your choice: '
options=("Apply [1]" "Destroy [2]" "Quit")
select opt in "${options[@]}"
do
    case $opt in
        "Apply [1]")
            echo "Applying changes"
            ${TF} apply -auto-approve -var="name=${NAME}" \
            -var="namespace=${NAMESPACE}"
            sleep 1
            break
            ;;
        "Destroy [2]")
            echo "Destroying changes"
             ${TF} destroy -auto-approve -var="name=${NAME}" \
            -var="namespace=${NAMESPACE}"
             sleep 1
             break
            ;;
        "Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
