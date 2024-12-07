package main

import (
	"burl/utils"
	"flag"
	"fmt"
)

const versionBurl = "1.0.0"

func main() {
	version := flag.Bool("version", false, "Show version")
	help := flag.Bool("help", false, "Show help")
	xmlToJson := flag.Bool("xtj", false, "Convert XML to JSON")
	flow := flag.Bool("flow", false, "Enable flow mode to run multiple files")
	isRunner := flag.Bool("runner", false, "Enable runner file to run test cases")

	flag.Parse()

	if *version {
		utils.PrintGreen(fmt.Sprintf("\n\tBurl version: %s\n\n", versionBurl))
		return
	}
	if *help {
		showHelp()
		return
	}
	if *xmlToJson {
		jsonStr := utils.XmlStrToJson(flag.Args()[0])
		fmt.Println(jsonStr)
		return
	}
	if len(flag.Args()) == 0 {
		utils.PrintRed("Error: Missing file argument...\n\n")
		return
	}
	runner := utils.GetRunnerByArgs(flag.Args())
	testing := utils.Testting{}
	switch {
	case *flow && *isRunner:
		utils.PrintRed("Error: Cannot use both flow and runner mode at the same time...\n\n")
	case *isRunner:
		testing.RunnerFile(flag.Args()[0])
	case *flow:
		testing.FlowMode(runner.Tests)
	default:
		testing.UnitestMode(runner.Tests)
	}

}

func showHelp() {
	fmt.Print("\n\nUsage: burl [options]\n\n")
	flag.PrintDefaults()
	fmt.Print("\n\n\tExample: \n\t\tburl  sample.json\n\t\tburl -flow  sample.json-1 sample.json-2\n\t\tburl -runner runner.sample.json\n\n")
	// print description
	fmt.Print("\n\n\tDescription: \n\n")
	fmt.Print("\t\tBurl is a tool for testing APIs using JSON configuration files.\n")
	fmt.Print("\t\tThe tool can be used in three modes:\n")
	fmt.Print("\t\t- Flow mode: to test multiple files in a sequence.\n")
	fmt.Print("\t\t- Runner mode: to run test cases from a file.\n\n")
	fmt.Print("\t\t]\n\n")

}
