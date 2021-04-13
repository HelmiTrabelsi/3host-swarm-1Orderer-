docker-compose -f host1.yaml up -d
docker-compose -f docker/docker-compose-ca1.yaml up -d
sleep 5
. organizations/fabric-ca/registerEnroll1.sh
createOrg1
