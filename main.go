package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type CurrencyDay struct {
	ID         string
	Date       time.Time
	DayNo      string
	Currencies []Currency
}

type Currency struct {
	Code            string
	CrossOrder      int
	Unit            int
	CurrencyNameTR  string
	CurrencyName    string
	ForexBuying     float64
	ForexSelling    float64
	BanknoteBuying  float64
	BanknoteSelling float64
	CrossRateUSD    float64
	CrossRateOther  float64
}

type tarih_Date struct {
	XMLName   xml.Name   `xml:"Tarih_Date"`
	Tarih     string     `xml:"Tarih,attr"`
	Date      string     `xml:"Date,attr"`
	Bulten_No string     `xml:"Bulten_No,attr"`
	Currency  []currency `xml:"Currency"`
}

type currency struct {
	Kod             string `xml:"Kod,attr"`
	CrossOrder      string `xml:"CrossOrder,attr"`
	CurrencyKod     string `xml:"CurrencyCode,attr"`
	Unit            string `xml:"Unit"`
	Isim            string `xml:"Isim"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`
	CrossRateUSD    string `xml:"CrossRateUSD"`
	CrossRateOther  string `xml:"CrossRateOther"`
}

func (c *CurrencyDay) GetData(CurrencyDate time.Time) {
	xDate := CurrencyDate
	t := new(tarih_Date)
	currDay := t.getData(CurrencyDate, xDate)
	for {
		if currDay == nil {
			CurrencyDate = CurrencyDate.AddDate(0, 0, -1)
			currDay = t.getData(CurrencyDate, xDate)
			if currDay != nil {
				break
			}
		} else {
			break
		}
	}

	if currDay != nil {
		*c = *currDay
	}
}

func (c Currency) Print() {
	fmt.Printf("Code: %s\n", c.Code)
	fmt.Printf("CurrencyName: %s\n", c.CurrencyName)
	fmt.Printf("CurrencyNameTR: %s\n", c.CurrencyNameTR)
	fmt.Printf("Unit: %d\n", c.Unit)
	fmt.Printf("ForexBuying: %.4f\n", c.ForexBuying)
	fmt.Printf("ForexSelling: %.4f\n", c.ForexSelling)
	fmt.Printf("BanknoteBuying: %.4f\n", c.BanknoteBuying)
	fmt.Printf("BanknoteSelling: %.4f\n", c.BanknoteSelling)
	fmt.Printf("CrossOrder: %d\n", c.CrossOrder)
	fmt.Printf("CrossRateUSD: %.4f\n", c.CrossRateUSD)
	fmt.Printf("CrossRateOther: %.4f\n", c.CrossRateOther)
}

func (c *CurrencyDay) Print() {
	fmt.Printf("ID: %s\n", c.ID)
	fmt.Printf("Date: %s\n", c.Date.Format("2006-01-02"))
	fmt.Printf("DayNo: %s\n", c.DayNo)

	for i, currency := range c.Currencies {
		fmt.Printf("\nCurrency %d\n", i+1)
		currency.Print()
	}
}

func (c *tarih_Date) getData(CurrencyDate time.Time, xDate time.Time) *CurrencyDay {
	currDay := new(CurrencyDay)
	var resp *http.Response
	var err error
	var url string
	currDay = new(CurrencyDay)
	url = "http://www.tcmb.gov.tr/kurlar/" + CurrencyDate.Format("200601") + "/" + CurrencyDate.Format("02012006") + ".xml"
	resp, err = http.Get(url)

	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			tarih := new(tarih_Date)
			d := xml.NewDecoder(resp.Body)
			marshalError := d.Decode(tarih)

			if marshalError != nil {
				log.Printf("err: %v", marshalError)

			}
			currDay.ID = xDate.Format("20060102")
			currDay.Date = xDate
			currDay.DayNo = tarih.Bulten_No
			currDay.Currencies = make([]Currency, len(tarih.Currency))

			for i, curr := range tarih.Currency {
				currDay.Currencies[i].Code = curr.CurrencyKod
				currDay.Currencies[i].CurrencyName = curr.CurrencyName
				currDay.Currencies[i].CurrencyNameTR = curr.Isim
				currDay.Currencies[i].BanknoteBuying = parseFloat(curr.BanknoteBuying)
				currDay.Currencies[i].BanknoteSelling = parseFloat(curr.BanknoteSelling)
				currDay.Currencies[i].ForexBuying = parseFloat(curr.ForexBuying)
				currDay.Currencies[i].ForexSelling = parseFloat(curr.ForexSelling)
				currDay.Currencies[i].CrossOrder = parseInt(curr.CrossOrder)
				currDay.Currencies[i].CrossRateOther = parseFloat(curr.CrossRateOther)
				currDay.Currencies[i].CrossRateUSD = parseFloat(curr.CrossRateUSD)

			}
			currDay.Print()
		} else {
			currDay = nil
		}
	}
	return currDay
}

func parseFloat(value string) float64 {
	if value == "" {
		return 0
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return parsed
}

func parseInt(value string) int {
	if value == "" {
		return 0
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return parsed
}

func main() {
	runtime.GOMAXPROCS(2)
	startTime := time.Now()
	currencyDay := new(CurrencyDay)
	currencyDate := time.Now()
	currencyDay.GetData(currencyDate)
	elepsedTime := time.Since(startTime)
	fmt.Printf("Execution Time : %s", elepsedTime)
}
