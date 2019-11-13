package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Goods struct{
	Id string `json:"id"`
	Name string `json:"name"`
	EndPrice string `json:"endPrice"`
	WinUser string `json:"winUser"`
	BidInfo []BidInfo `json:"bidInfos"`
}
type BidInfo struct{
	UserId string  `json:"userId"`
	Price string `json:"price"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "addGoods" {
		return s.addGoods(APIstub, args)
	} else if function == "queryAllGoods" {
		return s.queryAllGoods(APIstub)
	} else if function == "queryGoods" {
		return s.queryGoods(APIstub, args)
	} else if function == "updateBidInfo" {
		return s.updateBidInfo(APIstub, args)
	} else if function == "updateWinUser" {
		return s.updateWinUser(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name.")
}
 
func (s *SmartContract) addGoods(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("fail!")
	}
	var goods = Goods{Id: args[0], Name: args[1], EndPrice: "", BidInfo:[]BidInfo{}}
	goodsAsBytes, _ := json.Marshal(goods)
	APIstub.PutState(args[0], goodsAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllGoods(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	startKey := "GOODS0"
	endKey := "GOODS999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllGoods:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryGoods(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	goodsAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(goodsAsBytes)
}

func (s *SmartContract) updateBidInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	
	goodsAsBytes, err := APIstub.GetState(args[0])
	if err != nil{
		jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	} else if goodsAsBytes == nil{ // no State! error
		jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	}
	
	goods := Goods{}
	err = json.Unmarshal(goodsAsBytes, &goods)
	if err != nil {
		return shim.Error(err.Error())
	}
	// create rate structure
	var BidInfo = BidInfo{UserId: args[1], Price: args[2]}


	goods.BidInfo=append(goods.BidInfo,BidInfo)

	// update to User World state
	goodsAsBytes, err = json.Marshal(goods);

	APIstub.PutState(args[0], goodsAsBytes)

	return shim.Success([]byte("bidInfo is updated"))
}

func (s *SmartContract) updateWinUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	goodsAsBytes, _ := APIstub.GetState(args[0])
	goods := Goods{}

	json.Unmarshal(goodsAsBytes, &goods)
	goods.WinUser = args[1]
	goods.EndPrice = args[2]

	goodsAsBytes, _ = json.Marshal(goods)
	APIstub.PutState(args[0], goodsAsBytes)

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}