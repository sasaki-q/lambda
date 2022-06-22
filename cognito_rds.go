// cognito → lambda → rds

package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func demo1A(e events.CognitoEventUserPoolsPostConfirmation) error {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("USER"),
		os.Getenv("PASS"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("NAME"),
	)

	db, err := gorm.Open("postgres", url)

	if err != nil {
		fmt.Print("ERROR CONNECT === ", err)
		return err
	}

	res := db.Exec("INSERT INTO users (email) VALUES($1)", e.Request.UserAttributes["email"])

	if res.Error != nil {
		fmt.Print("ERROR SQL === ", res.Error)
		return res.Error
	}

	return nil
}

func demo1() {
	lambda.Start(demo1A)
}
