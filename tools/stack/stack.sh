#!/bin/bash

# Variables
STACK_NAME="MyBogieStack"
TEMPLATE_FILE="stack.json"
REGION="us-west-2" # Specify your AWS region

# Update CloudFormation stack
aws cloudformation update-stack \
  --stack-name "$STACK_NAME" \
  --template-body file://"$TEMPLATE_FILE" \
  --region "$REGION" \
  --capabilities CAPABILITY_IAM

# Check if the update command was successful
if [ $? -eq 0 ]; then
  echo "Successfully initiated update for stack $STACK_NAME."
else
  echo "Failed to update stack $STACK_NAME."
  exit 1
fi

# Optionally, you can wait for the stack update to complete
aws cloudformation wait stack-update-complete \
  --stack-name "$STACK_NAME" \
  --region "$REGION"

if [ $? -eq 0 ]; then
  echo "Stack $STACK_NAME update complete."
else
  echo "Failed to complete stack update for $STACK_NAME."
  exit 1
fi