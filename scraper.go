package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

func DownloadJSON() {
	os.Mkdir("json", 0755)
	file, err := os.Open("sites.csv")
	if err != nil {
		fmt.Printf("Error opening the CSV file: %s", err)
	}

	defer file.Close()
	csvRead := csv.NewReader(file)

	data, err := csvRead.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %s", err)
	}

	const baseURL string = "https://"
	const product string = "/products.json"
	var wg sync.WaitGroup

	for k, v := range data {
		if k == 0 {
			continue
		}

		fmt.Println(k, v[0], v[1])
		wg.Add(1)

		go func(v []string) {
			file, err := http.Get(baseURL + v[1] + product)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Body.Close()
			filename := fmt.Sprintf("json/" + v[0] + ".json")

			out, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			if file.StatusCode != http.StatusOK {
				fmt.Printf("Error writing: %s", file.Status)
			}

			_, err = io.Copy(out, file.Body)
			if err != nil {
				fmt.Printf("Error writing to the file: %s", err)
			}
			defer wg.Done()
		}(v)
	}
	wg.Wait()
	os.Chdir("json/")
	cmd := exec.Command("prettier", "-w", "*")
	stdOut, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running program: ", err)
	} else {
		fmt.Println(string(stdOut))
	}

	fmt.Println("All tasks completed.")
}
