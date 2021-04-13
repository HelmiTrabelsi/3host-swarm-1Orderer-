docker-compose -f host2.yaml up -d
docker-compose -f docker/docker-compose-ca2.yaml up -d
sleep 5
. organizations/fabric-ca/registerEnroll2.sh
createorg2
