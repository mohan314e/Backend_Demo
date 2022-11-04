# Backend_Demo

get gin and all other depedency by running "go get" commnad

run got/main.go

can use postman or curl or browser

GET method for list investment opportunity
http://localhost:8080/investmentopportunity

GET method for list investment opportunity in a specific type
//Stocks
http://localhost:8080/investmentopportunity?type=stocks
//Gold
http://localhost:8080/investmentopportunity?type=gold
//Mutual Funds
http://localhost:8080/investmentopportunity?type=mutualfunds
//Fixed Deposite
http://localhost:8080/investmentopportunity?type=fixeddeposite

GET method for PortFolio
Query Param --investmentamount (Mandatory)
Query Param --expectedannualreturn (Mandatory)
Query Param --type (Optional) --vaues = [gold, stocks, mutualfunds, fixeddeposite]

http://localhost:8080/getportfolio?investmentamount=100000&expectedannualreturn=19
