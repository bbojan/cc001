package main

import (
	"errors"
	"fmt"
	"strings"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "append" {
		return t.append(stub, args)
	} else if function == "sync" {
		return t.sync(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// append - invoke function to append value to key/value pair
func (t *SimpleChaincode) append(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running append()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to append")
	}

	key = args[0] //rename for funsies
	value = args[1]

	var oldValue, newValue string
	var err2 error

	oldVal, err2 := stub.GetState(key)
	if err2 == nil {
		oldValue = string(oldVal)
		newValue = oldValue + "|" + value

		err = stub.PutState(key, []byte(newValue)) //write the variable into the chaincode state
		if err != nil {
			return nil, err
		}
	}else{
		err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
		if err != nil {
			return nil, err
		}	
	}

	return nil, nil
}

// sync - invoke function to sync values
func (t *SimpleChaincode) sync(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running sync()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	var countKey, commandKeyPrefix, pozition, values string

	countKey = args[0] 
	commandKeyPrefix = args[1]
	pozition = args[1]
	values = args[3]
	
	var count, countIndex uint64
	var commands []string
	var err error

	count, err = stub.GetState(countKey)
	if err != nil {
		count = 0
	}	

	values = strings.Replace(values, "{commands:[", "", 1)
	values = strings.Replace(values, "],end:42}", "", 1)

	commands = strings.Split(values, ",")

	countIndex = count
    for i, command := range commands {
        if command != "" {
            //
			key := commandKeyPrefix + string(countIndex)
			err = stub.PutState(key, []byte(command)) 
			if err != nil {
				fmt.Println("err stub.PutState(key, []byte(command))")			
			}	
			//
			countIndex = countIndex + 1
        }
    }

	err = stub.PutState(countKey, []byte(string(countIndex))) 
	if err != nil {
		return nil, err
	}	

	var position, outIndex uint64
	var result string = ""

	position, err = strconv.ParseUint(pozition, 10, 64)
	if err != nil {
		position = 0
	}

	outIndex = position

	result = "{commands:["

	for i := position; i < countIndex; i++ {
		key := commandKeyPrefix + string(i)
		command, err := stub.GetState(key)
		if err != nil {
			fmt.Println("err stub.GetState(key)")		
		}else if command != ""{
			if i == position {
				result = result + command;
			}else if command != ""{				
				result = result + "," + command
			}
		}	
	}

	result = result + "],position:" + string(countIndex) + "}"

	return []byte(result), nil
}