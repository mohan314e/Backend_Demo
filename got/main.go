package main

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

var allInvestmentOpportunity []InvestmentOpportunity

func main() {
	allInvestmentOpportunity = getAllInvestmentOpportunity()
	router := gin.Default()
	router.GET("/investmentopportunity", getInvestmentOpportunity)
	router.GET("/getportfolio", getPortFolio)

	router.Run("localhost:8080")
}

func getInvestmentOpportunity(c *gin.Context) {

	Type := c.Request.URL.Query().Get("type")
	switch Type {
	case "stocks":
		c.IndentedJSON(http.StatusOK, stocks)
	case "gold":
		c.IndentedJSON(http.StatusOK, gold)
	case "mutualfunds":
		c.IndentedJSON(http.StatusOK, MutualFunds)
	case "fixeddeposite":
		c.IndentedJSON(http.StatusOK, FixedDeposite)
	default:
		c.IndentedJSON(http.StatusOK, allInvestmentOpportunity)
	}
}

func getPortFolio(c *gin.Context) {
	Type := c.Request.URL.Query().Get("type")
	InvestmentAmountStr := c.Request.URL.Query().Get("investmentamount")
	ExpectedAnnualReturnStr := c.Request.URL.Query().Get("expectedannualreturn")
	if InvestmentAmountStr == "" && ExpectedAnnualReturnStr == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Investment Amount and Expected Annual Return"})
		return
	}
	if InvestmentAmountStr == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Investment Amount"})
		return
	}
	if ExpectedAnnualReturnStr == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Expected Annual Return"})
		return
	}
	InvestmentAmount, _ := strconv.Atoi(InvestmentAmountStr)
	ExpectedAnnualReturn, _ := strconv.ParseFloat(ExpectedAnnualReturnStr, 32)
	switch Type {
	case "stocks":
		if InvestmentAmount < 2290 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Minimum Investment 2290"})
			return
		}
		if ExpectedAnnualReturn < 27 {
			ExpectedAnnualReturn = 27
		}
		if ExpectedAnnualReturn > 33 {
			ExpectedAnnualReturn = 33
		}
		resultportfolio := portfoliowrapper(stocks, InvestmentAmount, ExpectedAnnualReturn)
		c.IndentedJSON(http.StatusOK, resultportfolio)
	case "gold":
		if InvestmentAmount < 5100 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Minimum Investment 5100"})
			return
		}
		if ExpectedAnnualReturn < 25 {
			ExpectedAnnualReturn = 25
		}
		if ExpectedAnnualReturn > 32 {
			ExpectedAnnualReturn = 32
		}
		resultportfolio := portfoliowrapper(gold, InvestmentAmount, ExpectedAnnualReturn)
		c.IndentedJSON(http.StatusOK, resultportfolio)
	case "mutualfunds":
		if InvestmentAmount < 2000 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Minimum Investment 2000"})
			return
		}
		if ExpectedAnnualReturn < 17 {
			ExpectedAnnualReturn = 17
		}
		if ExpectedAnnualReturn > 25 {
			ExpectedAnnualReturn = 25
		}
		resultportfolio := portfoliowrapper(MutualFunds, InvestmentAmount, ExpectedAnnualReturn)
		c.IndentedJSON(http.StatusOK, resultportfolio)
	case "fixeddeposite":
		if InvestmentAmount < 3000 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Minimum Investment 3000"})
			return
		}
		if ExpectedAnnualReturn < 8.5 {
			ExpectedAnnualReturn = 8.5
		}
		if ExpectedAnnualReturn > 16 {
			ExpectedAnnualReturn = 16
		}
		resultportfolio := portfoliowrapper(FixedDeposite, InvestmentAmount, ExpectedAnnualReturn)
		c.IndentedJSON(http.StatusOK, resultportfolio)
	default:
		if InvestmentAmount < 2000 {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Minimum Investment 2000"})
			return
		}
		if ExpectedAnnualReturn < 8.5 {
			ExpectedAnnualReturn = 8.5
		}
		if ExpectedAnnualReturn > 33 {
			ExpectedAnnualReturn = 33
		}
		resultportfolio := portfoliowrapper(allInvestmentOpportunity, InvestmentAmount, ExpectedAnnualReturn)
		c.IndentedJSON(http.StatusOK, resultportfolio)
	}
}

func portfoliowrapper(invopparry []InvestmentOpportunity, InvestmentAmount int, ExpectedAnnualReturn float64) PortFolio {
	var resultportfolio PortFolio
	var portfolioslice []PortFolioInvestmentOpportunity
	portfoliocore(invopparry, 0, 0, InvestmentAmount, ExpectedAnnualReturn, portfolioslice, &resultportfolio)
	resultportfolio.OneYearForecastValue = int(float64(InvestmentAmount) * (1 + resultportfolio.AnnualReturnForecast/100))
	resultportfolio.ThreeYearForecastValue = int(float64(InvestmentAmount) * math.Pow(1+resultportfolio.AnnualReturnForecast/100, 3))
	resultportfolio.FiveYearForecastValue = int(float64(InvestmentAmount) * math.Pow(1+resultportfolio.AnnualReturnForecast/100, 5))
	return resultportfolio
}

func portfoliocore(invopparry []InvestmentOpportunity, idx int, depth int, InvestmentAmount int, ExpectedAnnualReturn float64, portfolioslice []PortFolioInvestmentOpportunity, resultportfolio *PortFolio) bool {
	var TotalUnit int
	var TotalForecast float64
	for _, portfoliostruct := range portfolioslice {
		TotalUnit += portfoliostruct.NoOfUnit
		TotalForecast += portfoliostruct.Next12MonthReturnForecast * float64(portfoliostruct.NoOfUnit)
	}
	AnnualReturnForecast := TotalForecast / float64(TotalUnit)
	if math.Abs(AnnualReturnForecast-ExpectedAnnualReturn) < 0.5 {
		resultportfolio.AnnualReturnForecast = AnnualReturnForecast
		resultportfolio.PortFolioInvestmentList = portfolioslice
		return true
	}
	if math.Abs(AnnualReturnForecast-ExpectedAnnualReturn) < math.Abs(ExpectedAnnualReturn-resultportfolio.AnnualReturnForecast) {
		resultportfolio.AnnualReturnForecast = AnnualReturnForecast
		resultportfolio.PortFolioInvestmentList = portfolioslice
	}
	if depth == 5 {
		return false
	}
	for i := idx; i < len(invopparry); i += 1 {
		var portfolio PortFolioInvestmentOpportunity
		portfolio.Type = invopparry[i].Type
		portfolio.Name = invopparry[i].Name
		portfolio.Last12MonthReturn = invopparry[i].Last12MonthReturn
		portfolio.Next12MonthReturnForecast = invopparry[i].Next12MonthReturnForecast
		portfolio.UnitPrice = invopparry[i].UnitPrice
		for j := 1; j <= InvestmentAmount/invopparry[i].UnitPrice; j += 1 {
			portfolio.NoOfUnit = j
			portfolioslice = append(portfolioslice, portfolio)
			retval := portfoliocore(invopparry, i+1, depth+1, InvestmentAmount-j*invopparry[i].UnitPrice, ExpectedAnnualReturn, portfolioslice, resultportfolio)
			portfolioslice = portfolioslice[:len(portfolioslice)-1]
			if retval {
				return true
			}
		}
	}
	return false
}

type InvestmentOpportunity struct {
	Type                      string
	Name                      string
	Last12MonthReturn         float64
	Next12MonthReturnForecast float64
	UnitPrice                 int
}
type PortFolioInvestmentOpportunity struct {
	Type                      string
	Name                      string
	Last12MonthReturn         float64
	Next12MonthReturnForecast float64
	UnitPrice                 int
	NoOfUnit                  int
}
type PortFolio struct {
	AnnualReturnForecast    float64
	OneYearForecastValue    int
	ThreeYearForecastValue  int
	FiveYearForecastValue   int
	PortFolioInvestmentList []PortFolioInvestmentOpportunity
}

func getAllInvestmentOpportunity() []InvestmentOpportunity {
	var allInvestmentOpportunity []InvestmentOpportunity
	var tempallInvestmentOpportunity = [][]InvestmentOpportunity{
		MutualFunds, gold, FixedDeposite, stocks,
	}
	for _, eachInvestmentOpportunity := range tempallInvestmentOpportunity {
		allInvestmentOpportunity = append(allInvestmentOpportunity, eachInvestmentOpportunity...)
	}
	return allInvestmentOpportunity
}

// Stocks slice
var stocks = []InvestmentOpportunity{
	{Type: "Stocks", Name: "Apple", Last12MonthReturn: 33, Next12MonthReturnForecast: 27, UnitPrice: 2290},
	{Type: "Stocks", Name: "Microsoft", Last12MonthReturn: 50, Next12MonthReturnForecast: 27, UnitPrice: 2300},
	{Type: "Stocks", Name: "Alphabet", Last12MonthReturn: 65, Next12MonthReturnForecast: 28, UnitPrice: 2350},
	{Type: "Stocks", Name: "Meta", Last12MonthReturn: 23, Next12MonthReturnForecast: 28, UnitPrice: 2390},
	{Type: "Stocks", Name: "Netflix", Last12MonthReturn: 19, Next12MonthReturnForecast: 29, UnitPrice: 2440},
	{Type: "Stocks", Name: "Amazon", Last12MonthReturn: 20, Next12MonthReturnForecast: 29, UnitPrice: 2510},
	{Type: "Stocks", Name: "Realiance", Last12MonthReturn: 22, Next12MonthReturnForecast: 29, UnitPrice: 2580},
	{Type: "Stocks", Name: "Adani", Last12MonthReturn: 28, Next12MonthReturnForecast: 30, UnitPrice: 2680},
	{Type: "Stocks", Name: "Airtel", Last12MonthReturn: 20, Next12MonthReturnForecast: 30, UnitPrice: 2710},
	{Type: "Stocks", Name: "P&G", Last12MonthReturn: 22, Next12MonthReturnForecast: 30, UnitPrice: 2710},
	{Type: "Stocks", Name: "Tata", Last12MonthReturn: 22, Next12MonthReturnForecast: 30, UnitPrice: 2710},
	{Type: "Stocks", Name: "Nike", Last12MonthReturn: 24, Next12MonthReturnForecast: 31, UnitPrice: 2710},
	{Type: "Stocks", Name: "Flipkart", Last12MonthReturn: 22, Next12MonthReturnForecast: 31, UnitPrice: 2800},
	{Type: "Stocks", Name: "Walmart", Last12MonthReturn: 26, Next12MonthReturnForecast: 31, UnitPrice: 2880},
	{Type: "Stocks", Name: "Dell", Last12MonthReturn: 27, Next12MonthReturnForecast: 32, UnitPrice: 2890},
	{Type: "Stocks", Name: "HP", Last12MonthReturn: 28, Next12MonthReturnForecast: 32, UnitPrice: 2895},
	{Type: "Stocks", Name: "Lenovo", Last12MonthReturn: 21, Next12MonthReturnForecast: 32, UnitPrice: 2895},
	{Type: "Stocks", Name: "DishTV", Last12MonthReturn: 21, Next12MonthReturnForecast: 33, UnitPrice: 2899},
	{Type: "Stocks", Name: "BMW", Last12MonthReturn: 23, Next12MonthReturnForecast: 33, UnitPrice: 3010},
	{Type: "Stocks", Name: "TVS", Last12MonthReturn: 24, Next12MonthReturnForecast: 33, UnitPrice: 3290},
}

// Gold slice
var gold = []InvestmentOpportunity{
	{Type: "Gold", Name: "PhysicalGold", Last12MonthReturn: 33, Next12MonthReturnForecast: 25, UnitPrice: 5100},
	{Type: "Gold", Name: "GoldETF", Last12MonthReturn: 37, Next12MonthReturnForecast: 25, UnitPrice: 5130},
	{Type: "Gold", Name: "GoldMutualFund", Last12MonthReturn: 37, Next12MonthReturnForecast: 25, UnitPrice: 5180},
	{Type: "Gold", Name: "SovereignGoldBonds", Last12MonthReturn: 39, Next12MonthReturnForecast: 26, UnitPrice: 5190},
	{Type: "Gold", Name: "DigitalGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 26, UnitPrice: 5200},
	{Type: "Gold", Name: "DurgaGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 26, UnitPrice: 5230},
	{Type: "Gold", Name: "FutureGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 27, UnitPrice: 5260},
	{Type: "Gold", Name: "WhiteGold", Last12MonthReturn: 30, Next12MonthReturnForecast: 27, UnitPrice: 5280},
	{Type: "Gold", Name: "BlackGold", Last12MonthReturn: 33, Next12MonthReturnForecast: 27, UnitPrice: 5300},
	{Type: "Gold", Name: "DigitalGold", Last12MonthReturn: 32, Next12MonthReturnForecast: 27, UnitPrice: 5340},
	{Type: "Gold", Name: "NightGold", Last12MonthReturn: 31, Next12MonthReturnForecast: 28, UnitPrice: 5380},
	{Type: "Gold", Name: "HawkGold", Last12MonthReturn: 28, Next12MonthReturnForecast: 28, UnitPrice: 5400},
	{Type: "Gold", Name: "PureGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 38, UnitPrice: 5430},
	{Type: "Gold", Name: "O24kGold", Last12MonthReturn: 27, Next12MonthReturnForecast: 28, UnitPrice: 5480},
	{Type: "Gold", Name: "MangalamGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 29, UnitPrice: 5490},
	{Type: "Gold", Name: "KonarkGold", Last12MonthReturn: 26, Next12MonthReturnForecast: 29, UnitPrice: 5500},
	{Type: "Gold", Name: "KalaMandir", Last12MonthReturn: 32, Next12MonthReturnForecast: 30, UnitPrice: 5550},
	{Type: "Gold", Name: "PhoenixGold", Last12MonthReturn: 29, Next12MonthReturnForecast: 30, UnitPrice: 5590},
	{Type: "Gold", Name: "FineGold", Last12MonthReturn: 33, Next12MonthReturnForecast: 31, UnitPrice: 6100},
	{Type: "Gold", Name: "ShinyGold", Last12MonthReturn: 30, Next12MonthReturnForecast: 31, UnitPrice: 6120},
	{Type: "Gold", Name: "StrikeGold", Last12MonthReturn: 30, Next12MonthReturnForecast: 32, UnitPrice: 6190},
	{Type: "Gold", Name: "BlueGold", Last12MonthReturn: 28, Next12MonthReturnForecast: 32, UnitPrice: 6200},
}

// MutualFunds slice
var MutualFunds = []InvestmentOpportunity{
	{Type: "MutualFunds", Name: "EquityFund", Last12MonthReturn: 20, Next12MonthReturnForecast: 17, UnitPrice: 2000},
	{Type: "MutualFunds", Name: "DebtFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 17, UnitPrice: 2000},
	{Type: "MutualFunds", Name: "HybridFund", Last12MonthReturn: 23, Next12MonthReturnForecast: 17, UnitPrice: 2010},
	{Type: "MutualFunds", Name: "GrowthFund", Last12MonthReturn: 24, Next12MonthReturnForecast: 17, UnitPrice: 2050},
	{Type: "MutualFunds", Name: "LiquidFund", Last12MonthReturn: 15, Next12MonthReturnForecast: 17, UnitPrice: 2090},
	{Type: "MutualFunds", Name: "PensionFund", Last12MonthReturn: 14, Next12MonthReturnForecast: 18, UnitPrice: 2100},
	{Type: "MutualFunds", Name: "LiquidFund", Last12MonthReturn: 15, Next12MonthReturnForecast: 18, UnitPrice: 2100},
	{Type: "MutualFunds", Name: "DebtFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 18, UnitPrice: 2140},
	{Type: "MutualFunds", Name: "NitrogenFund", Last12MonthReturn: 28, Next12MonthReturnForecast: 18, UnitPrice: 2200},
	{Type: "MutualFunds", Name: "CRFund", Last12MonthReturn: 23, Next12MonthReturnForecast: 29, UnitPrice: 2300},
	{Type: "MutualFunds", Name: "DRFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 19, UnitPrice: 2450},
	{Type: "MutualFunds", Name: "ETFFund", Last12MonthReturn: 18, Next12MonthReturnForecast: 29, UnitPrice: 2460},
	{Type: "MutualFunds", Name: "LifeFund", Last12MonthReturn: 19, Next12MonthReturnForecast: 20, UnitPrice: 2490},
	{Type: "MutualFunds", Name: "CrazyFund", Last12MonthReturn: 17, Next12MonthReturnForecast: 20, UnitPrice: 2500},
	{Type: "MutualFunds", Name: "RichFund", Last12MonthReturn: 19, Next12MonthReturnForecast: 20, UnitPrice: 2520},
	{Type: "MutualFunds", Name: "MillionaireFund", Last12MonthReturn: 20, Next12MonthReturnForecast: 21, UnitPrice: 2530},
	{Type: "MutualFunds", Name: "BillionaireFund", Last12MonthReturn: 20, Next12MonthReturnForecast: 21, UnitPrice: 2560},
	{Type: "MutualFunds", Name: "YouthFund", Last12MonthReturn: 19, Next12MonthReturnForecast: 22, UnitPrice: 2590},
	{Type: "MutualFunds", Name: "OldFund", Last12MonthReturn: 22, Next12MonthReturnForecast: 22, UnitPrice: 3000},
	{Type: "MutualFunds", Name: "FutureFund", Last12MonthReturn: 23, Next12MonthReturnForecast: 22, UnitPrice: 3100},
	{Type: "MutualFunds", Name: "PastFund", Last12MonthReturn: 24, Next12MonthReturnForecast: 23, UnitPrice: 3120},
	{Type: "MutualFunds", Name: "BreFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 23, UnitPrice: 3150},
	{Type: "MutualFunds", Name: "ColdFund", Last12MonthReturn: 22, Next12MonthReturnForecast: 23, UnitPrice: 3190},
	{Type: "MutualFunds", Name: "FusionFund", Last12MonthReturn: 27, Next12MonthReturnForecast: 23, UnitPrice: 3200},
	{Type: "MutualFunds", Name: "StarFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 24, UnitPrice: 3400},
	{Type: "MutualFunds", Name: "RedFund", Last12MonthReturn: 22, Next12MonthReturnForecast: 24, UnitPrice: 3500},
	{Type: "MutualFunds", Name: "SkyFund", Last12MonthReturn: 23, Next12MonthReturnForecast: 24, UnitPrice: 3590},
	{Type: "MutualFunds", Name: "RocketFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 25, UnitPrice: 3600},
	{Type: "MutualFunds", Name: "DebtFund", Last12MonthReturn: 20, Next12MonthReturnForecast: 25, UnitPrice: 3690},
	{Type: "MutualFunds", Name: "FuelFund", Last12MonthReturn: 21, Next12MonthReturnForecast: 25, UnitPrice: 3700},
}

// FixedDeposite slice
var FixedDeposite = []InvestmentOpportunity{
	{Type: "FixedDeposite", Name: "SBI", Last12MonthReturn: 7, Next12MonthReturnForecast: 8.5, UnitPrice: 3000},
	{Type: "FixedDeposite", Name: "ICICI", Last12MonthReturn: 7.5, Next12MonthReturnForecast: 8.8, UnitPrice: 3100},
	{Type: "FixedDeposite", Name: "HDFC", Last12MonthReturn: 7.5, Next12MonthReturnForecast: 8.8, UnitPrice: 3150},
	{Type: "FixedDeposite", Name: "AXIS", Last12MonthReturn: 7.4, Next12MonthReturnForecast: 8.8, UnitPrice: 3170},
	{Type: "FixedDeposite", Name: "HSBC", Last12MonthReturn: 8, Next12MonthReturnForecast: 8.9, UnitPrice: 3180},
	{Type: "FixedDeposite", Name: "PNB", Last12MonthReturn: 7.6, Next12MonthReturnForecast: 8.9, UnitPrice: 3185},
	{Type: "FixedDeposite", Name: "Bank of Baroda", Last12MonthReturn: 7.2, Next12MonthReturnForecast: 9.7, UnitPrice: 3190},
	{Type: "FixedDeposite", Name: "Bank of India", Last12MonthReturn: 7.3, Next12MonthReturnForecast: 9.9, UnitPrice: 3190},
	{Type: "FixedDeposite", Name: "Bank of Maharashtra", Last12MonthReturn: 7.5, Next12MonthReturnForecast: 9.9, UnitPrice: 3195},
	{Type: "FixedDeposite", Name: "Canara Bank", Last12MonthReturn: 7.6, Next12MonthReturnForecast: 10.5, UnitPrice: 3200},
	{Type: "FixedDeposite", Name: "Central Bank", Last12MonthReturn: 7.9, Next12MonthReturnForecast: 10.9, UnitPrice: 3259},
	{Type: "FixedDeposite", Name: "Indian Bank", Last12MonthReturn: 9.6, Next12MonthReturnForecast: 11, UnitPrice: 3270},
	{Type: "FixedDeposite", Name: "Indian Overseas Bank", Last12MonthReturn: 8.6, Next12MonthReturnForecast: 11.7, UnitPrice: 3300},
	{Type: "FixedDeposite", Name: "P&SB", Last12MonthReturn: 7.6, Next12MonthReturnForecast: 11.7, UnitPrice: 3350},
	{Type: "FixedDeposite", Name: "Union Bank of India", Last12MonthReturn: 7.6, Next12MonthReturnForecast: 12, UnitPrice: 3400},
	{Type: "FixedDeposite", Name: "UCOBank", Last12MonthReturn: 7.8, Next12MonthReturnForecast: 12, UnitPrice: 3500},
	{Type: "FixedDeposite", Name: "CSBBank", Last12MonthReturn: 8.6, Next12MonthReturnForecast: 12.4, UnitPrice: 3500},
	{Type: "FixedDeposite", Name: "City Union Bank", Last12MonthReturn: 8.6, Next12MonthReturnForecast: 12.5, UnitPrice: 3500},
	{Type: "FixedDeposite", Name: "DCBBank", Last12MonthReturn: 9.6, Next12MonthReturnForecast: 12.8, UnitPrice: 3700},
	{Type: "FixedDeposite", Name: "FederalBank", Last12MonthReturn: 9.6, Next12MonthReturnForecast: 13, UnitPrice: 3800},
	{Type: "FixedDeposite", Name: "IDBIBank", Last12MonthReturn: 9.7, Next12MonthReturnForecast: 13.7, UnitPrice: 3900},
	{Type: "FixedDeposite", Name: "IDFCBank", Last12MonthReturn: 10, Next12MonthReturnForecast: 13.9, UnitPrice: 3950},
	{Type: "FixedDeposite", Name: "IndusInd Bank", Last12MonthReturn: 10.6, Next12MonthReturnForecast: 14, UnitPrice: 4100},
	{Type: "FixedDeposite", Name: "RBLBank", Last12MonthReturn: 10.6, Next12MonthReturnForecast: 14.5, UnitPrice: 4200},
	{Type: "FixedDeposite", Name: "South Indian Bank", Last12MonthReturn: 10.9, Next12MonthReturnForecast: 15, UnitPrice: 4300},
	{Type: "FixedDeposite", Name: "Yes Bank", Last12MonthReturn: 11.6, Next12MonthReturnForecast: 15, UnitPrice: 4400},
	{Type: "FixedDeposite", Name: "UGB Bank", Last12MonthReturn: 12.6, Next12MonthReturnForecast: 16, UnitPrice: 3500},
	{Type: "FixedDeposite", Name: "Bank of America", Last12MonthReturn: 12.6, Next12MonthReturnForecast: 16, UnitPrice: 4500},
	{Type: "FixedDeposite", Name: "CitiBank India", Last12MonthReturn: 12.6, Next12MonthReturnForecast: 15.5, UnitPrice: 4500},
	{Type: "FixedDeposite", Name: "DNBBank", Last12MonthReturn: 12.6, Next12MonthReturnForecast: 16.5, UnitPrice: 4500},
}
