//Filename: cmd/api/currencyrate.go

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (app *application) currencyRate(w http.ResponseWriter, r *http.Request) {

	//get id from param
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	url := "https://api.apilayer.com/exchangerates_data/convert?to=BZD&from=USD&amount=" + strconv.Itoa(int(id))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "uTbJBEic4v1QoiCMxStWsudXZ4d4ijBg")

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)

	fmt.Fprintf(w, string(body))
}
