#!/bin/bash

OUTDIR=${1:-build}
mkdir -p "$OUTDIR"

OUTPUT=$OUTDIR/glusterblockrestd.service

cat >"$OUTPUT" <<EOF
[Unit]
Description=ReST server for Gluster Block Volumes management

[Service]
ExecStart=${SBINDIR}/glusterblockrestd --config=${SYSCONFDIR}/gluster-block-restapi/config.toml
KillMode=process

[Install]
WantedBy=multi-user.target

EOF
