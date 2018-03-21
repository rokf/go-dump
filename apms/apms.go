package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type data struct {
	Odhod  string
	Prihod string
	Voznja string
}

func main() {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://www.apms.si/response.ajax.php?", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) < 4 {
		fmt.Println("Usage: apms FROM TO DATE")
		os.Exit(1)
	}

	query := request.URL.Query()
	query.Add("com", "voznired")
	query.Add("task", "get")
	query.Add("datum", os.Args[3])
	query.Add("postaja_od", os.Args[1])
	query.Add("postaja_do", os.Args[2])
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if len(body) < 100 {
		fmt.Println("Incorrect params. An error message has been received.")
		os.Exit(1)
	}

	var data []data
	err = json.Unmarshal(body, &data)

	for _, x := range data {
		fmt.Println(x.Odhod, x.Prihod, x.Voznja)
	}
}
