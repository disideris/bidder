# bidder

Folder structure:
The implementation of the bidder avocarrot assignment lies in the bidder.go file whereas the tests in the bidder_test.go file. The response_json and response_json2 files are necessary for mocking the campaign api response.

Overview:
The bidder implementation leverages the built-in net/http package of go and a well known library of gorilla/mux for making the infrastructure for http routing. When bidder takes a bid in our bidder api, the appropriate handler, handles this request and after quering and retrieving the available campaigns from the campaign api it process them and constructs the minimum price from same country response.

Must be installed:
Golang 1.9+
git

How to run:
go get -u "github.com/gorilla/mux"
go get -u "github.com/stretchr/testify/assert"
go get -u "gopkg.in/h2non/gock.v1"

Inside the folder of bidder:
go build bidder.go
./bidder

Post json body to http://localhost:8000/bid
