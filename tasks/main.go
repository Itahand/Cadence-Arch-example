package main

import (
	"fmt"
	"io/ioutil"
	"log"

	//if you imports this with .  you do not have to repeat overflow everywhere
	. "github.com/bjartek/overflow/v2"
	"github.com/fatih/color"
)

func readJSFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {

	// Specify the path to your JavaScript file
	filePath := "bytecode/CadenceArch.js"

	// Read the content of the JavaScript file
	jsContent, err := readJSFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JavaScript file: %v", err)
	}

	//start an in memory emulator by default
	o := Overflow(
		WithGlobalPrintOptions(),
		// WithNetwork("testnet"),
	)

	fmt.Println("Testing Cadence Arch on Tesnet")
	// fmt.Println("Press any key to continue")
	// fmt.Scanln()

	//
	///// TESTING Cadence Arch /////
	//

	// color.Red("Should be able to create a COA")
	// Create COA inside Random's account
	o.Tx("create_COA",
		WithSigner("account"),
	).Print()
	// Fund COA inside the first slot
	color.Cyan("Fund a COA inside Random's account with Flow tokens")
	o.Tx("fund_COA",
		WithSigner("account"),
		WithArg("amount", "15.0"),
		WithArg("pathId", 0),
	).Print()

	// Get balance
	color.Cyan("Fetch balance from the COA inside Random's account")
	o.Script("get_coa_balance",
		WithArg("address", "random"),
		WithArg("pathId", 3),
	).Print()
	// Deploy a Solidity contract to the COA
	color.Cyan("Deploy a Solidity contract to Random's COA")
	events := o.Tx("deploy_contract_COA",
		WithSigner("random"),
		WithArg("code", jsContent),
		WithArg("pathId", 3),
	).Print().GetEventsWithName("TransactionExecuted")

	//fmt.Print(contractAddress, error)
	contractAddress := events[0].RawEvent.FieldsMappedByName()["contractAddress"]
	hexEncodedAddress := contractAddress.String()[3:]
	address := hexEncodedAddress[:len(hexEncodedAddress)-1]
	fmt.Print(address)

	// Confirm contract was deployed
	// by calling one of its functions
	color.Cyan("Confirm Contract exist by calling a function from it")
	o.Script("call_query_COA",
		WithArg("hexEncodedData", ""),
		WithArg("hexEncodedAddress", address),
		WithArg("address", "random"),
		WithArg("pathId", 3),
	).Print()

}
