#!/bin/bash
instruction=$1
version=$2

set -ev

#chaincode install
docker exec cli peer chaincode install -n auction -v $version -p github.com
#chaincode instatiate
docker exec cli peer chaincode $instruction -n auction -v $version -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 5
#chaincode invoke user1
#docker exec cli peer chaincode invoke -n auction -C mychannel -c '{"Args":["addUser","user1"]}'
#sleep 5
#chaincode query user1
#docker exec cli peer chaincode query -n auction -C mychannel -c '{"Args":["readRating","user1"]}'

#chaincode invoke add rating
#docker exec cli peer chaincode invoke -n auction -C mychannel -c '{"Args":["addRating","user1","p1","5.0"]}'
#sleep 5

echo '-------------------------------------END-------------------------------------'
