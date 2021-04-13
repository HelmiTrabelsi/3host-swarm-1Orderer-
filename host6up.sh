docker-compose -f host6.yaml up -d
docker-compose -f docker/docker-compose-ca-orderer.yaml up -d
sleep 5
. organizations/fabric-ca/registerOrderer.sh
createOrderer
