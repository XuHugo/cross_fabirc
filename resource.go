package main

import (
	"fmt"
	"strconv"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const MSG_DATA_KEY = "msg_data"
const INIT_METHOD = "init_method"

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("resource Init")
	var allMethodList []string

	var method_storeCrossData =`{
        'inputs' : [{
			'name':'crossData',
			'type':'string'
		}],
		'outputs': [],
		'name': 'storeCrossData',
		'type': 'function'
	}`
	var method_callCOntractStoreMsg = `{
        'inputs' : [{
			'name':'msgKey',
			'type':'string'
		},{
			'name':'msgData',
			'type':'string'
		}],
		'outputs': [],
		'name': 'callContractStoreMsg',
		'type': 'function'
	}`

	allMethodList.push(method_callCOntractStoreMsg)
	allMethodList.push(method_storeCrossData)
	err = stub.PutState(INIT_METHOD, []byte(allMethodList))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	if function == "storeCrossData" {
		return t.storeCrossData(stub, args)
	} else if function == "callContractStoreMsg" {
		return t.callContractStoreMsg(stub, args)
	}else if function == "queryCrossData" {
		return t.queryCrossData(stub, args)
	}else if function == "queryAllMethod" {
		return t.queryAllMethod(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"apply\" \"query\"")
}

func (t *SimpleChaincode) storeCrossData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	crossData := args[0]

	data, err := stub.GetState(MSG_DATA_KEY)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	var dataList []string

	dataList.push(crossData)
	err = stub.PutState("msg_data", []byte(dataList))
	if err != nil {
		return shim.Error(err.Error())
	}
}

func (t *SimpleChaincode) callContractStoreMsg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	msgKey := args[0]
	msgData := args[1]

	err = stub.PutState(msgKey, []byte(msgData))
	if err != nil {
		return shim.Error(err.Error())
	}
}

func (t *SimpleChaincode) queryCrossData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	data, err := stub.GetState(MSG_DATA_KEY)
	if err != nil {
		return shim.Error("Failed to get state")
	}

	return shim.Success(data)

}

func (t *SimpleChaincode) queryAllMethod(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	data, err := stub.GetState(INIT_METHOD)
	if err != nil {
		return shim.Error("Failed to get state")
	}

	return shim.Success(data)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


