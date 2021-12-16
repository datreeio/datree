#!/bin/bash
set -ex

mkdir -p ~/.aws

cat > ~/.aws/credentials << EOL
[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}
EOL

aws configure list
aws s3 cp windows_install.ps1 s3://get.datree.io/windows_install.ps1 --acl public-read
aws s3 cp install.sh s3://get.datree.io/install.sh --acl public-read

echo "Cloudfront: Invalidating /*"
aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID --paths "/*"
