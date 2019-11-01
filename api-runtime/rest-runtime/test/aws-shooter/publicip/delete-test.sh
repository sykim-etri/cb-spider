RESTSERVER=localhost

ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-ohio-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-ohio-config

ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-oregon-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-oregon-config

ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-singapore-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-singapore-config

ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-paris-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-paris-config

ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-saopaulo-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-saopaulo-config


ID=`curl -X GET http://$RESTSERVER:1024/publicip?connection_name=aws-tokyo-config |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
curl -X GET http://$RESTSERVER:1024/publicip/${ID}?connection_name=aws-tokyo-config |json_pp
