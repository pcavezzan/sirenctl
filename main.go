package main

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/pcavezzan/sirenctl/pkg/services"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"time"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

const csvCharSeparator = ';'
const apiTemporisationInSeconds = 200
const apiTemporisation = apiTemporisationInSeconds * time.Millisecond

func main() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.String("siren", "", "Un numéro de SIREN")
	f.String("codepostal", "", "Un code postal")
	f.StringP("file", "f", "", "Le chemin vers le fichier d'entrée des SIREN/Code Postal.")
	f.StringP("output", "o", "", "Le chemin vers le fichier de sortie des établissements liés au SIREN/Code Postal.")
	f.BoolP("verbose" ,"v",false, "Mode verbeux")
	f.Parse(os.Args[1:])
	// "time" and "type" may have been loaded from the config file, but
	// they can still be overridden with the values from the command line.
	// The bundled posflag.Provider takes a flagset from the spf13/pflag lib.
	// Passing the Koanf instance to posflag helps it deal with default command
	// line flag values that are not present in conf maps from previously loaded
	// providers.
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("error loading config from flags: %v", err)
	}

	verbose := k.Bool("verbose")
	siren := k.String("siren")
	zipCode := k.String("codepostal")
	file := k.String("file")
	outputFilePath := k.String("output")
	if verbose {
		fmt.Println(" SIREN: ", siren)
		fmt.Println(" Code Postal: ", zipCode)
		fmt.Println("File: ", file)
	}

	var sirenZipCodeParser services.SearchSirenApiParser
	if isSet(zipCode) && isSet(siren) {
		sirenZipCodeParser = services.NewArgumentZipCodeParser(siren, zipCode)
	} else if isSet(file) {
		sirenZipCodeParser = services.NewCsvFileSirenZipCodeParser(file, csvCharSeparator)
	}

	if sirenZipCodeParser == nil {
		fmt.Println("Please, set either file or siren and codepostal.")
		os.Exit(1)
	}


	searchesSirenApi, err := sirenZipCodeParser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	formatter := services.NewCsvAgencyFormatter(csvCharSeparator)
	service := services.HttpSirenService{}
	var agencies []services.Agency
	for _, searchSirenApi := range searchesSirenApi {
		agenciesFromApi, err := service.GetAgencies(searchSirenApi)
		if err != nil {
			log.Fatalln(err)
		}
		agencies = append(agencies, agenciesFromApi...)
		if verbose {
			fmt.Printf("Sleep %d seconds for temporise api call from data.gouv.fr.\n", apiTemporisationInSeconds)
		}
		time.Sleep(apiTemporisation)
	}

	format, err := formatter.Format(agencies)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	if outputFilePath == "" {
		fmt.Println(format)
	} else {
		if err = os.WriteFile(outputFilePath, []byte(format), 0644); err != nil {
			fmt.Printf("Could not write file %s because: %v", outputFilePath, err)
			os.Exit(1)
		}
	}
}


func isSet(code string) bool {
	return code != ""
}