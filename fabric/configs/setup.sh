WORKSPACE=`pwd`
FABRIC_HOME=$GOPATH"/src/github.com/hyperledger/fabric"
FABRIC_GO_SDK=$GOPATH"/src/github.com/hyperledger/fabric-sdk-go"
CHANNEL_ID="traceability-system"
ORG1_MSP="Org1MSP"
ORG2_MSP="Org2MSP"
FABRIC_VERSION="master"
clear(){
    rm -rf crypto-config artifacts
    cd $FABRIC_HOME
    echo "切换fabric版本为"$FABRIC_VERSION
    git checkout $FABRIC_VERSION
    docker ps -qa | xargs docker stop
    docker ps -qa | xargs docker rm
    docker rm -f $(docker ps -aq)
    docker network prune
    docker volume prune
}


BLACKLISTED_VERSIONS="^1\.0\. ^1\.1\.0-preview ^1\.1\.0-alpha"
IMAGETAG="latest"
# Do some basic sanity checking to make sure that the appropriate versions of fabric
# binaries/images are available.  In the future, additional checking for the presence
# of go or other items could be added.
function checkPrereqs() {
  # Note, we check configtxlator externally because it does not require a config file, and peer in the
  # docker image because of FAB-8551 that makes configtxlator return 'development version' in docker
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

  for UNSUPPORTED_VERSION in $BLACKLISTED_VERSIONS; do
    echo "$LOCAL_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      echo "ERROR! Local Fabric binary version of $LOCAL_VERSION does not match this newer version of BYFN and is unsupported. Either move to a later version of Fabric or checkout an earlier version of fabric-samples."
      exit 1
    fi

    echo "$DOCKER_IMAGE_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      echo "ERROR! Fabric Docker image version of $DOCKER_IMAGE_VERSION does not match this newer version of BYFN and is unsupported. Either move to a later version of Fabric or checkout an earlier version of fabric-samples."
      exit 1
    fi
  done
}

cryptoGenerate(){
    echo "进入"$WORKSPACE
    cd $WORKSPACE
    cryptogen generate --config=crypto-config.yaml
}

configtxGenerate(){
    cd $WORKSPACE
    mkdir artifacts
    configtxgen --profile TwoOrgsOrdererGenesis -outputBlock ./artifacts/genesis.block 
    echo "生成创世块./artifacts/genesis.block"

    configtxgen --profile TwoOrgsChannel -outputCreateChannelTx ./artifacts/${CHANNEL_ID}.tx -channelID $CHANNEL_ID
    echo "创建通道配置交易./artifacts/${CHANNEL_ID}.tx"

    configtxgen --profile TwoOrgsChannel -outputAnchorPeersUpdate ./artifacts/${ORG1_MSP}Anchors.tx -channelID $CHANNEL_ID -asOrg $ORG1_MSP
    echo "创建瞄节点配置交易./artifacts/${ORG1_MSP}Anchors.tx"
     configtxgen --profile TwoOrgsChannel -outputAnchorPeersUpdate ./artifacts/${ORG2_MSP}Anchors.tx -channelID $CHANNEL_ID -asOrg $ORG2_MSP
    echo "创建瞄节点配置交易./artifacts/${ORG2_MSP}Anchors.tx"
}

startNetwork(){
    
    docker-compose -f docker-compose-cli.yaml up -d
    checkPrereqs
}

downNetwork(){
    echo "downNetwork() nothing to do"
}

ARG1=$1
if ["${ARG1}"=="start"];then
    startNetwork
elif ["${ARG1}"=="generate"];then
    cryptoGenerate
    configtxGenerate
elif ["{$ARG1}"=="down"];then
    downNetwork
else
    downNetwork
    clear
    cryptoGenerate
    configtxGenerate
    startNetwork
fi