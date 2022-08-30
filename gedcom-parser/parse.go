package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/iand/gedcom"
)

func parseFile(data []byte, header *multipart.FileHeader) (string, error) {
	gedcomDecoder := gedcom.NewDecoder(bytes.NewReader(data))
	var peopleList string

	decodedGedcom, err := gedcomDecoder.Decode()
	if err != nil {
		return peopleList, err
	}

	for _, record := range decodedGedcom.Individual {
		if len(record.Name) > 0 {
			log.Printf("Found name: %s", record.Name[0].Name)
			peopleList += fmt.Sprintf("%s\n", record.Name[0].Name)
		}
	}
	return peopleList, nil
}
