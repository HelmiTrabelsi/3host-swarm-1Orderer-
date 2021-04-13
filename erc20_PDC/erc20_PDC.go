/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

const (
	BalancePrefix   = `BALANCE`
	AllowancePrefix = `APPROVE`
)

var (
	ErrNotEnoughFunds                   = errors.New(`not enough funds`)
	ErrForbiddenToTransferToSameAccount = errors.New(`forbidden to transfer to same account`)
	ErrSpenderNotHaveAllowance          = errors.New(`spender not have allowance for amount`)
)

type Balance struct {
	Amount float64 `json:"Amount"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabcar_cc")

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "Init" { //create a new marble
		return t.Init(stub)
	}
	if function == "invokeTransfer" { //transfer money
		return t.invokeTransfer(stub, args)
	}
	if function == "BalanceOf" { //get balance
		return t.BalanceOf(stub, args)
	}
	if function == "SetBalance" { //set balance
		return t.SetBalance(stub, args)
	}
	if function == "GetUserIdentity" { //get user identity atrribute
		return t.GetUserIdentity(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// Init initializes the chaincode
func (t *SimpleChaincode) invokeTransfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	type AmountTransientInput struct {
		Amount string `json:"amount"`
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	AmountAsBytes, ok := transMap["amount"]
	if !ok {
		return shim.Error("Amount must be a key in the transient map")
	}

	if len(AmountAsBytes) == 0 {
		return shim.Error("Amount value in the transient map must be a non-empty JSON string")
	}

	var AmountInput AmountTransientInput
	err = json.Unmarshal(AmountAsBytes, &AmountInput)
	if err != nil {
		return shim.Error("Failed to decode JSON " + err.Error())
	}

	// transfer target
	ReceiverMspId := args[0]
	ReceiverCertId := args[1]

	//transfer amount
	Amount, _ := strconv.ParseFloat(AmountInput.Amount, 32)

	SenderMspId := args[2]
	SenderCertId := args[3]

	// Disallow to transfer Receiverken to same account
	if SenderMspId == ReceiverMspId && SenderCertId == ReceiverCertId {
		return shim.Error("forbidden to transfer to same account")
	}

	ReceiverBalance, err := getBalance(stub, ReceiverMspId, ReceiverCertId)
	if err != nil {
		return shim.Error("Can not get Receiver balance" + err.Error())
		//return shim.Error(collectionName)
	}

	SenderBalance, err := getBalance(stub, SenderMspId, SenderCertId)
	if err != nil {
		return shim.Error("Can not get Sender balance")
	}
	out, _ := json.Marshal(SenderBalance)
	// Check the funds sufficiency
	if SenderBalance-Amount < 0 {

		return shim.Error(string(out))
	}

	// Update payer and Receiver balance
	if err = setBalance(stub, SenderMspId, SenderCertId, SenderBalance-Amount); err != nil {
		return shim.Error("Can not update Sender balance")
	}

	if err = setBalance(stub, ReceiverMspId, ReceiverCertId, ReceiverBalance+Amount); err != nil {
		return shim.Error("Can not update receipient balance")
	}

	return shim.Success([]byte(out))
}

func (t *SimpleChaincode) BalanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	MSPID := args[0]
	CertID := args[1]

	InvokerMspId, _ := cid.GetMSPID(stub)
	InvokerCertId, _ := cid.GetID(stub)

	if MSPID != InvokerMspId || CertID != InvokerCertId {
		return shim.Error("You are not allowed to get this balance")
	}

	Balance, err := getBalance(stub, MSPID, CertID)
	if err != nil {
		return shim.Error(err.Error())
	}
	BalanceAsBytes, _ := json.Marshal(Balance)

	return shim.Success(BalanceAsBytes)

}

func (t *SimpleChaincode) SetBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	MSPID := args[0]
	CertID := args[1]
	Amount2, _ := strconv.ParseFloat(args[2], 32)
	//if err != nil {
	//	return shim.Error("Parse failed" + err.Error())
	//}
	err2 := setBalance(stub, MSPID, CertID, Amount2)
	out, _ := json.Marshal(Amount2)
	if err2 != nil {
		return shim.Error("Error: Balance not setted: " + err2.Error())
	}
	return shim.Success([]byte(out))

}

func (t *SimpleChaincode) GetUserIdentity(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	MSPID, _ := cid.GetMSPID(stub)
	CertID, _ := cid.GetID(stub)

	identity := fmt.Sprintf(MSPID, CertID)

	return shim.Success([]byte(identity))

}

// Internal Function
// setBalance puts balance value to state
func balanceKey(ownerMspId, ownerCertId string) string {
	s := fmt.Sprintf(BalancePrefix, ownerMspId, ownerCertId)
	return s
}

func getBalance(stub shim.ChaincodeStubInterface, mspId, certId string) (float64, error) {
	var balance Balance

	collectionName := getCollectionName(stub, mspId)
	BalanceAsBytes, err := stub.GetPrivateData(collectionName, balanceKey(mspId, certId))
	if err != nil {
		return -1, err
	}
	json.Unmarshal(BalanceAsBytes, &balance)
	return balance.Amount, nil
}

// setBalance puts balance value to state
func setBalance(stub shim.ChaincodeStubInterface, mspId, certId string, amount float64) error {
	balance := Balance{
		Amount: amount,
	}
	BalanceAsBytes, _ := json.Marshal(balance)

	collectionName := getCollectionName(stub, mspId)
	return stub.PutPrivateData(collectionName, balanceKey(mspId, certId), BalanceAsBytes)

}

//get collection name
func getCollectionName(stub shim.ChaincodeStubInterface, MSPID string) string {
	collection_name := ""
	if MSPID == "Org1MSP" {
		collection_name = "BOC_collection"
	} else if MSPID == "Org2MSP" {
		collection_name = "org2_collection"
	} else if MSPID == "Org3MSP" {
		collection_name = "org3_collection"
	} else if MSPID == "Org4MSP" {
		collection_name = "org4_collection"
	} else if MSPID == "Org5MSP" {
		collection_name = "org5_collection"
	} else {
		fmt.Printf("No Organization collection")
	}
	return collection_name
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
