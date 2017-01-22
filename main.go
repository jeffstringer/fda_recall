package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type jsonProduct struct {
	Recall struct {
		Date               string `json:"release_date"`
		BrandName          string `json:"name"`
		ProductDescription string `json:"product_description"`
		Reason             string `json:"reason"`
		CompanyReleaseLink string `json:"company_release_link"`
	}	`json:"recall"`
}

type Product struct {
	Date               string `xml:"DATE"`
	BrandName          string `xml:"BRAND_NAME"`
	ProductDescription string `xml:"PRODUCT_DESCRIPTION"`
	Reason             string `xml:"REASON"`
	Company            string `xml:"COMPANY"`
	CompanyReleaseLink string `xml:"COMPANY_RELEASE_LINK"`
}

type Recall struct {
	Products []Product `xml:"PRODUCT"`
}

// GET to fda dataset site
func getFdaJson() {
	fdaUrl := "http://www.fda.gov/DataSets/Recalls/Food/Food.xml"
	response, err := http.Get(fdaUrl)
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if response.Status == "200 OK" {
		var r Recall
		err = xml.Unmarshal(body, &r)
		if err != nil {
			panic(err.Error())
		}

		// convert to JSON
		var oneProduct jsonProduct
		var allProducts []jsonProduct

		for _, value := range r.Products {
			oneProduct.Recall.Date = value.Date
			oneProduct.Recall.BrandName = value.BrandName
			oneProduct.Recall.ProductDescription = value.ProductDescription
			oneProduct.Recall.Reason = value.Reason
			oneProduct.Recall.CompanyReleaseLink = value.CompanyReleaseLink
			allProducts = append(allProducts, oneProduct)
		}

		jsonData, err := json.Marshal(allProducts)
		if err != nil {
			fmt.Println(err)
		}

		// POST to rails app
		railsUrl := "http://localhost:3000/recalls"
		var jsonStr = string(jsonData)
		request, err := http.Post(railsUrl, "application/json", strings.NewReader(jsonStr))
		fmt.Println("I am getting fda data...", time.Now())
		if err != nil {
			println(err)
		}
		defer request.Body.Close()
	}
	if err != nil {
		panic(err.Error())
	}
}

// initiate process with chron job
func main() {
    ch := gocron.Start()
    gocron.Every(10).Seconds().Do(getFdaJson)

    <-ch
}
