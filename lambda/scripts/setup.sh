#!/bin/bash

# Runs initial setup of the infrastructure:
# - IAM policies for lambdas
# - Creates and configures DynamoDB tables

. scripts/_config.sh

# Create IAM Policy and associated Trusted relationship
read -r -d '' POLICY_TEMPLATE <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
			"Effect": "Allow",
			"Action": [
				"dynamodb:*"
			],
			"Resource": [
        "arn:aws:dynamodb:*:*:table/*"
      ]
		}
	]
}
EOF

# Remove whitespaces
POLICY=`echo $POLICY_TEMPLATE | jq -c -M .`

echo "Create Policy"
POLICY_ARN=`aws iam create-policy --policy-name $POLICY_NAME --policy-document "$POLICY" | jq -M -r .Policy.Arn`

echo $POLICY > /tmp/POLICY.json

echo "AWS IAM role trust relationsip"
read -r -d '' TRUST_RELATIONSHIP <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
echo $TRUST_RELATIONSHIP > /tmp/TRUST_RELATIONSHIP.json

echo "Create Role"
ROLE_ARN=`aws iam create-role \
              --role-name $ROLE_NAME \
              --assume-role-policy-document file:///tmp/TRUST_RELATIONSHIP.json | jq -M -r .Role.Arn`

echo "Attach policy to the role"
aws iam attach-role-policy \
  --role-name $ROLE_NAME \
  --policy-arn $POLICY_ARN

rm /tmp/POLICY.json
rm /tmp/TRUST_RELATIONSHIP.json

# Setup DynamoDB tables
echo "Setup DynamoDB tables: Captchas"
aws dynamodb create-table \
    --table-name $DYNAMO_TABLE_CAPTCHAS \
    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=secret,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=secret,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST

echo "Setup DynamoDB tables: Forms"
aws dynamodb create-table \
    --table-name $DYNAMO_TABLE_FORMS \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST
