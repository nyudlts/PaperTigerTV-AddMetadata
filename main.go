package main

import (
	"bufio"
	"fmt"
	"strconv"

	go_aspace "github.com/nyudlts/go-aspace/lib"
	"os"
	"strings"
)

type Recording struct {
	ResourceId	string
	RefID		string
	URI			string
	Indicator1	string
	Indicator2	string
	Indicator3	string
	Title		string
	Component 	string
	Scope 		string
	StartDate	string
}

func (r Recording) Fill (resourceID string, refId string, uri string, indicator1 string, indicator2 string, indicator3 string, title string, cuid string) Recording {
	r.ResourceId = resourceID
	r.RefID = refId
	r.URI = uri
	r.Indicator1 = indicator1
	r.Indicator2 =  indicator2
	r.Indicator3 = indicator3
	r.Title = title
	r.Component = cuid
	return r
}

func (r Recording) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", r.RefID, r.ResourceId, r.URI, r.Indicator1, r.Indicator2, r.Indicator3, r.Title, r.Component, r.Scope, r.StartDate)
}

func main() {
	var aspace = go_aspace.Client
	tsv, err := os.Open("pttv_original.tsv")
	if err != nil {
		panic(err)
	}
	defer tsv.Close()


	outputFile, err := os.Create("pttv_md.tsv")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	fw := bufio.NewWriter(outputFile)
	fw.WriteString("Resource Id\tRef Id\tURI\tIndicator 1\tIndicator 2\tIndicator 3\tTitle\tCUID\tScopeContent\tStart Date\n")
	fw.Flush()


	scanner := bufio.NewScanner(tsv)
	recordings := []Recording{}

	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "\t")
		recording := Recording{}
		recording = recording.Fill(cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7])
		recordings = append(recordings, recording)
	}

	for i, r := range recordings {
		if i == 0 { continue }
		asuriSplit := strings.Split(r.URI, "/")
		rid, _ := strconv.Atoi(asuriSplit[2])
		aid, _ := strconv.Atoi(asuriSplit[4])
		ao, err := aspace.GetArchivalObjectById(rid, aid)
		if err != nil {
			panic(err)
		}
		r.StartDate = ao.Dates[0].Begin
		scope, err := GetScope(ao.Notes)
		if err != nil {
			r.Scope = err.Error()
		} else {
			r.Scope = scope
		}

		fw.WriteString(r.String())
		fw.Flush()
	}
	fw.Flush()
	os.Exit(0)
}

func GetScope(notes []*go_aspace.Note) (string, error) {
	for _,n := range notes {
		if n.Type == "scopecontent" {
			return n.Subnotes[0].Content, nil
		}
	}
	return "", fmt.Errorf("No Scope Note Found")
}
