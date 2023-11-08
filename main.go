package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	inputFilePath  string
	outputFilePath string
)

const (
	defaultInputFilePath  = "./channels.json"
	defaultOutputFilePath = "./channels_parsed.csv"
)

type Channel struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Members []string       `json:"members"`
	Created int            `json:"created"`
	Topic   TopicOrPurpose `json:"topic"`
	Purpose TopicOrPurpose `json:"purpose"`
}

type TopicOrPurpose struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int    `json:"last_set"`
}

func main() {

	// Command line flags for alternate input/output file paths
	inp := flag.String("input", defaultInputFilePath, "The full path of the file.json to be parsed. Default is channels.json in the current folder.")
	outp := flag.String("output", defaultOutputFilePath, "The full path of the file_parsed.csv to be created. Default is channels_parsed.csv in the current folder.")
	flag.Parse()
	if *inp != "" {
		inputFilePath = *inp
	}
	if *outp != "" {
		outputFilePath = *outp
	}

	// Reads the input file using the "os" package in the directory this package was run in unless flagged
	b, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmarshals the .json file into a list of Channel structs
	channels := []Channel{}
	err = json.Unmarshal(b, &channels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Creates the output file using the "os" package in the directory this package was run in unless flagged
	output, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer output.Close()            // closes the file once the script ends
	writer := csv.NewWriter(output) // creates a new CSV Writer

	// Writes the 1st Row of Data as Headers
	headers := []string{}
	for _, channel := range channels {
		headers = append(headers, fmt.Sprintf("%s (%s)", channel.Name, channel.ID))
	}
	writer.Write(headers)

	// Calculate the maximum number of rows necessary to iterate
	var currentRow, maxRows int
	for _, channel := range channels {
		if len(channel.Members) > maxRows {
			maxRows = len(channel.Members)
		}
	}

	// Iterate the maximum rows
	for currentRow <= maxRows {
		row := []string{}
		for _, channel := range channels {
			if len(channel.Members) > currentRow {
				row = append(row, channel.Members[currentRow])
			} else {
				row = append(row, "")
			}
		}

		// Write the Row of Members under Channels
		writer.Write(row)
		currentRow += 1
	}

	// Finalize the output
	writer.Flush()

}
