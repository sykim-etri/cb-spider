#!/bin/bash

if [ "$1" = "" ]; then
	echo
	echo -e 'usage: '$0' mock|aws|azure|gcp|alibaba|tencent|ibm|openstack|ncp|nhncloud'
	echo -e '\n\tex) '$0' aws'
	echo
	exit 0;
fi

# common setup.env path
SETUP_PATH=$CBSPIDER_ROOT/test/vm-cb-user-validation-cli/common
source $SETUP_PATH/setup.env $1

echo "============== before create VPC/Subnet: '${VPC_NAME}'"
$CLIPATH/spctl  vpc create -d \
    '{
      "ConnectionName":"'${CONN_CONFIG}'",
      "ReqInfo": {
        "Name": "'${VPC_NAME}'",
        "IPv4_CIDR": "'${VPC_CIDR}'",
        "SubnetInfoList": [
          {
            "Name": "'${SUBNET_NAME}'",
            "IPv4_CIDR": "'${SUBNET_CIDR}'"
          }
        ]
      }
    }' 2> /dev/null

echo "============== after create VPC/Subnet: '${VPC_NAME}'"

echo -e "\n\n"

