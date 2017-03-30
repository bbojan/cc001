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
	} else if function == "push" {
		return t.push(stub, args)
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
	}else if function == "pull" { 
		return t.pull(stub, args)
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

// push - invoke function to push values
func (t *SimpleChaincode) push(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running push()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var countKey, commandKeyPrefix, values, separator string

	countKey = args[0] 
	commandKeyPrefix = args[1]
	values = args[2]
	separator = args[3]
	
	var count, countIndex uint64
	var commands []string
	var countBytes []byte
	var err error

	countBytes, err = stub.GetState(countKey)
	if err != nil {
		count = 0
	}else{
		var countString = string(countBytes)
		count, err = strconv.ParseUint(countString, 10, 64)
		if err != nil{
			count = 0
		}
	}	

	commands = strings.Split(values, separator)

	countIndex = count
    for _, command := range commands {
        if command != "" {
            //
			var key = commandKeyPrefix + strconv.FormatUint(countIndex, 10)
			err = stub.PutState(key, []byte(command)) 
			if err != nil {
				fmt.Println("err stub.PutState(key, []byte(command))")			
			}	
			//
			countIndex = countIndex + 1
        }
    }

	err = stub.PutState(countKey, []byte(strconv.FormatUint(countIndex, 10))) 
	if err != nil {
		return nil, err
	}	

	return nil, nil
}

// pull - invoke function to pull values
func (t *SimpleChaincode) pull(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running pull()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var countKey, commandKeyPrefix, pozition string

	countKey = args[0] 
	commandKeyPrefix = args[1]
	pozition = args[2]
	
	var count uint64
	var countBytes []byte
	var err error

	countBytes, err = stub.GetState(countKey)
	if err != nil {
		count = 0
		//return nil, errors.New("count = 0, err != nil")
	}else{
		var countString = string(countBytes)
		count, err = strconv.ParseUint(countString, 10, 64)
		if err != nil{
			count = 0
			//return nil, errors.New("count = 00, err != nil")
		}else{
			//return nil, errors.New("err == nil : " + strconv.FormatUint(count, 10))
		}		
	}	

	var position, outIndex uint64
	var result string = ""

	position, err = strconv.ParseUint(pozition, 10, 64)
	if err != nil {
		position = 0
	}

	var commandBytes []byte
	var command string

	outIndex = position

	result = "{\"commands\":["

	for i := position; i < count; i++ {
		var key = commandKeyPrefix + strconv.FormatUint(i, 10)
		commandBytes, err = stub.GetState(key)
		if err != nil {
			fmt.Println("err stub.GetState(key)")		
		}else {
			command = string(commandBytes)
			if command != ""{
				if i == position {
					result = result + command;
				}else if command != ""{				
					result = result + "," + command
				}
				outIndex = outIndex + 1
			}	
		}
	}

	result = result + "],\"position\":" + strconv.FormatUint(count, 10) + "}"

	return []byte(result), nil
}