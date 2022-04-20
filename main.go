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
        Email:        "hajo-ride-sheets@secret-274010.iam.gserviceaccount.com",
        PrivateKey:   []byte("-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDBKml/4fSda4Zb\nGD2XsHxxs7Vc4cxCXJJe9jAszLcozA6epNrdupKqNLJ3KGJPG0iqfoXVm3pwRpVi\nTvrXWKlRqHNOnzHTzQTNn+cQ3BbyLlsfH8hkm9/JEeAnlDowg2/bOOioZmN82AvR\nKRrFlSD6IxIHVVAiCneM9IJnUrM2DAWg8o0fN1YMrcy40WaCdel57WhUf56sSLaq\nR2HJpO9k+Mg2w41bdZD6mJRFGgZuJIma3JU2KBSR+pnLwRymptYqTuPNFCOaxQuR\n8BX1Adi+6byXmxD57eCpZWz7Jw46zC95epuk9P9jZCBkJ0cpisZEY/mWUxU/WeQR\nJYzw54enAgMBAAECggEAGdVO5VgWBh+P3mBFjHiBNsDulmIUQSPEPN2hNRHub48A\nEt4VRJC5IKOKN+liDGunsjWC4ehjMt9zBwIXsJ5zi9otYZcPa4tVrceezWBsfoMN\n6Kxazn8SKf/cAdm5q9dIKayuC3zLUrJKEkpJrExKFr8t9WZxdbv/SmNA4AntXGLP\naBGhjzEhRGR9wMbEwFD44BE+8IfuTO1bSvV3OVzejavHN8D6zvvuI7eYxC6A9B05\nzPrPx5zsWGxZXa3SFy2BU+8iPd4UaRt60QSNqlMX/oyDUQXPw6Oj1MshJ5HZTm0E\nG905xmcNvjQvRZGaFXiXyhNuiDtTUWX/k3gqaUjJoQKBgQD21PxalQ9l1Q1Jt0kx\nbjJ08bUYp7HdsOIxN1lJOn73+t3RE0+/5N8pkDwyf+t9ALskbT7IG9Nab4Od8aT9\n8QbGVBErRvE+qlrfEEfTdtXBGS0qF5E9lNouqP7793pg409TqZMoSVnObRhRPLIg\n4KzlRBYuKqw5BOLBYw3XXQ6GfwKBgQDIVyM6Uu6yRVylqxR/kTlkavMn4WlUzw2Q\n8qAkyG68ed45sDrjfL0U7CF77YKphNTYSpevKWu5kc45neMI0Ch4Anfg/ypnBGV2\nvTivGikkjGN4xYVaeDJNLoViWdQuEF26AJWZsJjopRR8F9+2Dg+D9pwr7ya1OVsg\nJQ4r4wF62QKBgDxilK4QY2j1Oz/6gASAwgVMEJLMi74BhQSnVseMcP13uQZqJAt9\nEX9YkvGDy/eT0zSxKTfJXuZz+44LTOL3MUIZ1O1yeJhtErgZgojgPKD16M+wXk7j\neJELtP2taJoFYiXEiYTElRzkIZvWLnnOnV2lT4vR0rpMULEg64Txdhk/AoGBAMIY\nWhL5ojSpKrd+5l8pTddSHr42jxwRde1Gg4zfexmzYG+3o3Yetfn9Q/uTApJGFQhw\nTEDx+mJRNuFvbDRvp6WNjqM309b21W6SqMZMrCuyr/SkOSyLYwCP376HlSADzyqT\nwlyiUaKjvRL09lKlkLajSG7wRioxpW2hOPef+SERAoGBAPaj+Fdya3bmuWJfZwnx\nep/EPWd9HG6iPkc+60YTTZanxHKK54rzZh1yk9QI1B+MRFXD0o2HMoYiZzbR3mC/\nYgzbIwK/Dky5XuDDdE8a43+gI1rT5y80jpmspCT9FhZ7eEHjXiI6PVQP1mDRG6st\n6iaBukb1sJVnFmKQufAqDme4\n-----END PRIVATE KEY-----\n"),
        PrivateKeyID: "caefee48194f0eeedbb05a51c8d74f4faf8df421",
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