package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

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
			// Print person and their sex
			log.Printf("Found: %s (%s)", record.Name[0].Name, record.Sex)

			// Add person (and sex) to string of people peopleList
			peopleList += fmt.Sprintf(
				"%s (%s) ",
				record.Name[0].Name,
				record.Sex)
		}
	}
	peopleList = strings.ReplaceAll(peopleList, "/", "")
	return peopleList, nil
}
