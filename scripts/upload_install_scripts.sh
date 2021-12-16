#!/bin/bash
set -ex

aws test-connection
aws s3 cp windows_install.ps1 s3://get.datree.io/windows_install.ps1 --acl public-read
aws s3 cp install.sh s3://get.datree.io/install.sh --acl public-read

echo "Cloudfront: Invalidating /*"
aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID --paths "/*"
