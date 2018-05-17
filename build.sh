#!/bin/sh

docker build -t ssm-parent -f Dockerfile.build .
docker run --rm -v "$(pwd):/tmp/builder" ssm-parent cp ssm-parent /tmp/builder
