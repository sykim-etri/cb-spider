
echo "####################################################################"
echo "## SecurityGroup Test Scripts for CB-Spider IID Working Version     "
echo "##   SecurityGroup: Delete"
echo "####################################################################"

$CBSPIDER_ROOT/interface/spider security delete --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n sg-01
