docker-compose -f docker/docker-compose-ca1.yaml down
docker-compose -f docker/docker-compose-ca2.yaml down
docker-compose -f docker/docker-compose-ca-orderer.yaml down
docker-compose -f host1.yaml down
docker-compose -f host2.yaml down

docker-compose -f host6.yaml down --volumes --remove-orphans


docker rm -f $(docker ps -a -q)
docker volume rm $(docker volume ls -q)


docker rmi $(docker images dev* -q)
rm -rf /tmp/hfc-*

sudo rm -rf organizations/fabric-ca/org1/msp organizations/fabric-ca/org1/tls-cert.pem organizations/fabric-ca/org1/ca-cert.pem organizations/fabric-ca/org1/IssuerPublicKey organizations/fabric-ca/org1/IssuerRevocationPublicKey organizations/fabric-ca/org1/fabric-ca-server.db

sudo rm -rf organizations/fabric-ca/org2/msp organizations/fabric-ca/org2/tls-cert.pem organizations/fabric-ca/org2/ca-cert.pem organizations/fabric-ca/org2/IssuerPublicKey organizations/fabric-ca/org2/IssuerRevocationPublicKey organizations/fabric-ca/org2/fabric-ca-server.db



sudo rm -rf organizations/fabric-ca/ordererOrg/msp organizations/fabric-ca/ordererOrg/tls-cert.pem organizations/fabric-ca/ordererOrg/ca-cert.pem organizations/fabric-ca/ordererOrg/IssuerPublicKey organizations/fabric-ca/ordererOrg/IssuerRevocationPublicKey organizations/fabric-ca/ordererOrg/fabric-ca-server.db

rm -rf channel-artifacts/*
