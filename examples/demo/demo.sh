#!/usr/bin/env bash
# Showcase kentik-cloudexport Terraform provider against live Kentik API.

SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
REPO_DIR=$(cd -- "$SCRIPT_DIR" && cd ../../ && pwd)

source "$REPO_DIR/tools/utility_functions.sh" || exit 1

run() {
    pushd "$SCRIPT_DIR" > /dev/null || die

    check_prerequisites
    check_env
    cleanup_tf_files

    stage "Build & install kentik-cloudexport Terraform provider"
    pause_and_run make --directory "$REPO_DIR" install || die
    pause

    stage "Initialize Terraform"
    pause_and_run pygmentize ./providers.tf
    pause

    pause_and_run terraform init || die
    pause

    stage "Create AWS cloud export"
    pause_and_run pygmentize ./aws.tf
    pause

    pause_and_run terraform plan || die
    pause
    pause_and_run terraform apply -auto-approve || die
    pause

    list_cloud_exports

    stage "Update AWS cloud export"
    pause_and_run sed -i 's/terraform aws cloud export/updated description/g' ./aws.tf
    pause_and_run pygmentize ./aws.tf
    pause

    pause_and_run terraform plan || die
    pause
    pause_and_run terraform apply -auto-approve || die
    pause
    sed -i 's/updated description/terraform aws cloud export/g' ./aws.tf

    list_cloud_exports

    stage "Delete AWS cloud export"
    pause_and_run terraform destroy -auto-approve

    list_cloud_exports

    popd > /dev/null || exit
}

check_prerequisites() {
    if ! terraform -v > /dev/null 2>&1; then
        die "Please install Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli"
    fi

    if ! pygmentize -V > /dev/null 2>&1; then
        die "Please install Pygments: https://pygments.org/"
    fi

    if ! curl -V > /dev/null 2>&1; then
        die "Please install cURL: https://curl.se/"
    fi

    if ! jq -V > /dev/null 2>&1; then
        die "Please install jq: https://stedolan.github.io/jq/"
    fi
}

cleanup_tf_files() {
    rm -rf .terraform .terraform.lock.hcl terraform.tfstate
}

check_env() {
    stage "Check auth env variables"

    if [[ -z "$KTAPI_AUTH_EMAIL" ]]; then
        die "KTAPI_AUTH_EMAIL env variable must be set to Kentik API account email"
    fi

    if [[ -z "$KTAPI_AUTH_TOKEN" ]]; then
        die "KTAPI_AUTH_TOKEN env variable must be set to Kentik API authorization token"
    fi

    echo "Print KTAPI_AUTH_EMAIL"
    echo "$KTAPI_AUTH_EMAIL"
    echo "Print KTAPI_AUTH_TOKEN (first 10 chars)"
    echo "${KTAPI_AUTH_TOKEN:0:10}"

    pause
}

list_cloud_exports() {
    read -r -p "Press any key to list Cloud Exports with cURL on https://cloudexports.api.kentik.com/cloud_export/v202101beta1/exports"
    curl --location --request GET --max-time 30 "https://cloudexports.api.kentik.com/cloud_export/v202101beta1/exports" \
        --header "X-CH-Auth-Email: $KTAPI_AUTH_EMAIL" \
        --header "X-CH-Auth-API-Token: $KTAPI_AUTH_TOKEN" | jq
    pause
}

run
