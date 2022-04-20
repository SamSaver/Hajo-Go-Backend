package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/joho/godotenv"
)

type Ride struct {
    Name     string `form:"name" binding:"required"`
    Email string `form:"email" binding:"required"`
	Phone string `form:"phone" binding:"required"`
	From string `form:"from" binding:"required"`
	To string `form:"to" binding:"required"`
	Dep string `form:"dep" binding:"required"`
	Vehicle string `form:"vehicle" binding:"required"`
	HasGroup string `form:"hasGroup"`
	GroupSize string `form:"groupSize"`
	ShareRides string `form:"shareRides"`
}

type Trip struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	From string `json:"from" binding:"required"`
	To string `json:"to" binding:"required"`
	Dep string `json:"dep" binding:"required"`
	Vehicle string `json:"vehicle" binding:"required"`
	Interests string `json:"int" binding:"required"`
}

func HomePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Server is up and running :)"})
}

// POST /book_trip
func BookTrip(c *gin.Context) {

	var json Trip
    c.Bind(&json)

	name := json.Name
	phone := json.Phone
	email := json.Email
	from := json.From
	to := json.To
	dep := json.Dep
	vehicle := json.Vehicle
	interests := json.Interests

	date:= strings.Split(dep, " ")[0]
	time:= strings.Split(dep, " ")[1]

	fmt.Println("Interests:")
	fmt.Println(interests)

	ctx := context.Background()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Create a JWT configurations object for the Google service account
    conf := &jwt.Config{
        Email:        os.Getenv("EMAIL"),
        PrivateKey:   []byte(os.Getenv("PRIVATEKEY")),
        PrivateKeyID: os.Getenv("PRIVATEKEYID"),
        TokenURL:     "https://oauth2.googleapis.com/token",
        Scopes: []string{
            "https://www.googleapis.com/auth/spreadsheets",
        },
    }

    client := conf.Client(ctx)

    // Create a service object for Google sheets
    srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets service: %v", err)
	}
	
	// The ID of the spreadsheet to update.
	spreadsheetId := "1mSdXHbPkgKG237AocbGEvDLucaXpoAwpvNnYFSbT6Mw" // TODO: Update placeholder value.

	// The A1 notation of a range to search for a logical table of data.
	// Values will be appended after the last row of the table.
	range2 := "Sheet2!A:J"

	// How the input data should be interpreted.
	valueInputOption := "USER_ENTERED"

	// How the input data should be inserted.
	insertDataOption := "INSERT_ROWS"

	rb := &sheets.ValueRange{
			// TODO: Add desired fields of the request body.
			Range: range2,
			Values: [][]interface{}{
				{name, email, phone, from, to, date, time, vehicle, interests},
			},
	}

	resp, err := srv.Spreadsheets.Values.Append(spreadsheetId, range2, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
			log.Fatal(err)
	}
	fmt.Println(resp)
	// TODO: Change code below to process the `resp` object:
	c.JSON(http.StatusOK, gin.H{"message": "Successfully booked a trip"})

	
}

// POST book_ride
func BookRide (c *gin.Context){

	var json Ride
    c.Bind(&json)

	name := json.Name
	phone := json.Phone
	email := json.Email
	from := json.From
	to := json.To
	dep := json.Dep
	vehicle := json.Vehicle
	has_group := json.HasGroup
	group_size := json.GroupSize
	share_rides := json.ShareRides

	date:= strings.Split(dep, " ")[0]
	time:= strings.Split(dep, " ")[1]

	ctx := context.Background()

	// Create a JWT configurations object for the Google service account
    conf := &jwt.Config{
        Email:        os.Getenv("EMAIL"),
        PrivateKey:   []byte(os.Getenv("PRIVATEKEY")),
        PrivateKeyID: os.Getenv("PRIVATEKEYID"),
        TokenURL:     "https://oauth2.googleapis.com/token",
        Scopes: []string{
            "https://www.googleapis.com/auth/spreadsheets",
        },
    }

    client := conf.Client(ctx)

    // Create a service object for Google sheets
    srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets service: %v", err)
	}
	
	// The ID of the spreadsheet to update.
	spreadsheetId := "1mSdXHbPkgKG237AocbGEvDLucaXpoAwpvNnYFSbT6Mw" // TODO: Update placeholder value.

	// The A1 notation of a range to search for a logical table of data.
	// Values will be appended after the last row of the table.
	range2 := "Sheet1!A:K"

	// How the input data should be interpreted.
	valueInputOption := "USER_ENTERED"

	// How the input data should be inserted.
	insertDataOption := "INSERT_ROWS"

	rb := &sheets.ValueRange{
			// TODO: Add desired fields of the request body.
			Range: range2,
			Values: [][]interface{}{
				{name, phone, email, from, to, date, time, vehicle, has_group, group_size, share_rides},
			},
	}

	resp, err := srv.Spreadsheets.Values.Append(spreadsheetId, range2, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
			log.Fatal(err)
	}
	fmt.Println("Response:")
	fmt.Println(resp)
	// TODO: Change code below to process the `resp` object:
	c.JSON(http.StatusOK, gin.H{"message": "Successfully booked a ride"})
}


func main() {
	
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/", HomePage)
	r.POST("/book_trip", BookTrip)
	r.POST("/book_ride", BookRide)
  	r.Run(":3000")
}