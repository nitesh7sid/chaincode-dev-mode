package main

/* Imports
 * 1 utility libraries for formatting
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"fmt"
		"crypto/x509"
		"encoding/pem"
		"bytes"
		"github.com/hyperledger/fabric/core/chaincode/shim"
		sc "github.com/hyperledger/fabric/protos/peer"

)

// Define the Smart Contract structure
type SmartContract struct {
}

var logger = shim.NewLogger("hello_world")

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("########### hello_world Init ###########")
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	logger.Info("########### hello_world Invoke ###########")
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "readFunc1" {
		return s.readFunc1(APIstub, args)
	} else if function == "writeFunc1" {
		return s.writeFunc1(APIstub,args)
	}

logger.Errorf("Unknown action, check the first argument, must be one of 'readFunc1', 'writeFunc1'. But got: %v", args[0])
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) readFunc1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	logger.Infof("key = %d", args[0])

	getSigner(APIstub)

	valAsBytes, _ := APIstub.GetState(args[0])

	logger.Infof("Query Response:%s\n", valAsBytes)
	return shim.Success(valAsBytes)
}

func (s *SmartContract) writeFunc1(APIstub shim.ChaincodeStubInterface,args []string) sc.Response {

	logger.Infof("key = %d, value = %d\n", args[0], args[1])

		getSigner(APIstub)
		APIstub.PutState(args[0], []byte(args[1]))

		return shim.Success(nil)
}

func getSigner(APIstub shim.ChaincodeStubInterface)sc.Response {

	creator, err := APIstub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}
	certStart := bytes.Index(creator, []byte("-----BEGIN CERTIFICATE-----"))
	if certStart == -1 {
		return shim.Error("No Certificate found")
	}
	certText := creator[certStart:]
	block, _ := pem.Decode(certText)
	if block == nil {
		return shim.Error("Error received on pem.Decode of certificate:" + string(certText))
	}

	ucert, err := x509.ParseCertificate(block.Bytes)
	fmt.Println("Creator: ", ucert.Subject.CommonName)
	return shim.Success([]byte (ucert.Subject.CommonName))
}


func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
