// Copyright 2020 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func GenerateCSVFile(path string, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create csv file %s", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(rows)
	csvWriter.Flush()
	return nil
}

func GenerateHTMLFile(sourceCSV string, targetHTML string) error {
	data, err := ioutil.ReadFile(sourceCSV)
	if err != nil {
		return fmt.Errorf("failed to read csv file %s", err)
	}
	htmlTemplate, err := Asset("templates/single_chart.html")
	if err != nil {
		return fmt.Errorf("failed to load asset: %s", err)
	}
	viewTemplate, err := template.New("chart").Parse(string(htmlTemplate))
	if err != nil {
		return fmt.Errorf("failed to parse html template %s", err)
	}
	htmlFile, err := os.OpenFile(targetHTML, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open html file %s", err)
	}
	defer htmlFile.Close()
	return viewTemplate.Execute(htmlFile, map[string]interface{}{
		"Data": string(data),
	})
}

func GenerateJSONFile(jsonData []byte, targetJSON string) error {
	jsonFile, err := os.OpenFile(targetJSON, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create json file %s", err)
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("failed to write json data", err)
		return err
	}
	return nil
}

func CheckOutputLocation(outputLocation string) (string, error) {
	dirInfo, err := os.Stat(outputLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return outputLocation, fmt.Errorf("output location (%s) is not existed: %s", outputLocation, err)
		}
		return outputLocation, fmt.Errorf("output location (%s) has error: %s", outputLocation, err)
	} else {
		if !dirInfo.IsDir() {
			return outputLocation, fmt.Errorf("output location (%s) is not directory", outputLocation)
		}
		if dirInfo.Mode().Perm()&(1<<(uint(7))) == 0 {
			return outputLocation, fmt.Errorf("output location (%s) is not writable", outputLocation)
		}
	}
	return outputLocation, nil
}
