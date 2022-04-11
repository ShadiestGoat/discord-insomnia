package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	TMP_JSON_EG = "_'_EG_JSON_HERE_'_"
)

var (
	RegShortenTable = regexp.MustCompile(`-{2,}`)
	RegShortenSpace = regexp.MustCompile(`(\t| ){2,}`)
	RegRemoveWarningType = regexp.MustCompile(`\n> (info)|(warn)`)
	RegEntryStart = regexp.MustCompile(`## .+?%.+?\n`)
	RegJSONEg = regexp.MustCompile("(?s)```json.+?```")
	RegJSONEgTMP = regexp.MustCompile(TMP_JSON_EG)
)

func PrepareDoc(inp []byte) []byte {
	entLoc := RegEntryStart.FindIndex(inp)
	if len(entLoc) == 0 {
		return []byte{}
	}
	inp = inp[entLoc[0]:]

	inp = RegShortenTable.ReplaceAll(inp, []byte(":-:"))
	inp = RegRemoveWarningType.ReplaceAll(inp, []byte(""))
	
	
	jsonEgs := [][]byte{}
	
	for RegJSONEg.Match(inp) {
		loc := RegJSONEg.FindIndex(inp)
		curJSON := []byte{}

		for _, b := range inp[loc[0]:loc[1]] {
			curJSON = append(curJSON, b)
		}

		jsonEgs = append(jsonEgs, curJSON)
		
		tmpInp := inp[:loc[0]]
		tmpInp = append(tmpInp, []byte(TMP_JSON_EG)...)
		tmpInp = append(tmpInp, inp[loc[1]:]...)
		inp = tmpInp
	}
	
	inp = RegShortenSpace.ReplaceAll(inp, []byte(" "))
	
	i := 0

	for RegJSONEgTMP.Match(inp) {
		loc := RegJSONEgTMP.FindIndex(inp)
		tmpInp := []byte{}
		
		for _, b := range inp[:loc[0]] {
			tmpInp = append(tmpInp, b)
		}

		for _, b := range jsonEgs[i] {
			tmpInp = append(tmpInp, b)
		}

		for _, b := range inp[loc[1]:] {
			tmpInp = append(tmpInp, b)
		}

		inp = tmpInp
		i++
	}

	return inp
}

func main() {
	files, err := ioutil.ReadDir("resources")
	PanicIfErr(err)

	allOutputs := []RequestGroup{}

	for _, file := range files {
		if file.IsDir() {continue}
		name := file.Name()
		if name[len(name)-3:] != ".md" {
			continue
		}

		fmt.Printf("Parsing %v...\n", name)

		content, err := ioutil.ReadFile(filepath.Join("resources", name))
		PanicIfErr(err)
		fmt.Println("preparing...")
		content = PrepareDoc(content)
		if len(content) == 0 {
			continue
		}
		
		fmt.Println("Parsing requests...")

		parsed := RequestGroup{
			Name:     name[:len(name)-3],
			Requests: ParseDoc(content),
		}

		allOutputs = append(allOutputs, parsed)
	}

	output, err := json.Marshal(GenerateExport(allOutputs))
	PanicIfErr(err)
	outputF, err := os.Create("Output.json")
	PanicIfErr(err)
	outputF.Write(output)
	outputF.Close()
	fmt.Println("Donezo!")
}

