#!/usr/bin/env bash

readonly DEVELOPMENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

usage () {
    echo "Usage: \$ ${BASH_SOURCE[1]} /path/to/config.yaml"
    exit 1
}

readonly CONFIG_PATH="${1}"
if [[ -z "${CONFIG_PATH}" ]]; then
    usage
fi

readonly UPLOADER="${DEVELOPMENT_DIR}/tools/cmd/configuploader/main.go"
if [[ ! -d "${DEVELOPMENT_DIR}/tools/vendor/github.com" ]]; then
    (cd "${DEVELOPMENT_DIR}/tools" && dep ensure -v -vendor-only)
fi

readonly CONFIG="${HOME}/.kube/config"

go run "${UPLOADER}" --kubeconfig "${CONFIG}" --config-path "${CONFIG_PATH}"