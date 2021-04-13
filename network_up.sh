  LOCAL_VERSION=$(configtxlator version | sed -ne 's/ Version: //p')
  DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-tools:$IMAGETAG peer version | sed -ne 's/ Version: //p' | head -1)

  echo "LOCAL_VERSION=$LOCAL_VERSION"
  echo "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

if [ "$LOCAL_VERSION" != "$DOCKER_IMAGE_VERSION" ]; then
    echo "=================== WARNING ==================="
    echo "  Local fabric binaries and docker images are  "
    echo "  out of  sync. This may cause problems.       "
    echo "==============================================="
  fi
./host1up.sh
./host2up.sh
./host6up.sh
#. organizations/fabric-ca/registerEnroll1.sh
#createOrg1

#. organizations/fabric-ca/registerEnroll2.sh
#createorg2

#. organizations/fabric-ca/registerEnroll3.sh
#createorg3

#. organizations/fabric-ca/registerEnroll4.sh
#createorg4

#. organizations/fabric-ca/registerEnroll5.sh
#createorg5

#. organizations/fabric-ca/registerOrderer.sh
#createOrderer



