#!/bin/sh
set -eo pipefail

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

chart_dir="$1"
gs_bucket="$2"

if [ -z "$chart_dir" ] || [ -z "$gs_bucket" ] ; then
    echo "Usage: publish-helm-chart.sh CHART_PATH GS_BUCKET" >&2
    exit 2
fi

tmp="$(mktemp -d)"
chart_full_dir="$(realpath "$chart_dir")"

(cd "$chart_dir" && helm dep update)
cd "$tmp"
helm package "$chart_full_dir"
gsutil -q cp gs://$gs_bucket/index.yaml ./index.yaml
helm repo index --url https://$gs_bucket.storage.googleapis.com/ --merge ./index.yaml ./
gsutil -q rsync ./ gs://kluster-charts
