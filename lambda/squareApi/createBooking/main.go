package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/brianvoe/gofakeit/v6"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Build request to send to Square Create Bookings API
	verVarIdStr := os.Getenv("SERVICE_VARIATION_VERSION")
	verVarIdInt, err := strconv.Atoi(verVarIdStr)
	if err != nil {
		log.Println("err parsing version variation id: ", err)
		return events.APIGatewayProxyResponse{Body: string("Error parsing square response"), StatusCode: 500}, err
	}
	req := requestBody{
		IdempotencyKey: gofakeit.UUID(),
		Booking: booking{
			Version:      1,
			StartAt:      time.Now().Add(time.Duration(time.Duration(gofakeit.IntRange(1, 24)) * time.Hour)),
			LocationID:   os.Getenv("LOCATION_ID"),
			CustomerID:   os.Getenv("CUSTOMER_ID"),
			CustomerNote: request.QueryStringParameters["message"],
			SellerNote:   "Thank you for your business!",
			AppointmentSegments: []segments{
				{
					DurationMinutes:         45,
					ServiceVariationID:      os.Getenv("SERVICE_VARIATION_ID"),
					TeamMemberID:            os.Getenv("TEAM_MEMBER_ID"),
					ServiceVariationVersion: int64(verVarIdInt),
				},
			},
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error parsing square payload"), StatusCode: 400}, err
	}

	postReq, err := http.NewRequest("POST", "https://connect.squareupsandbox.com/v2/bookings", bytes.NewBuffer(reqBytes))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error building square request"), StatusCode: 400}, err
	}
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Square-Version", "2023-08-16")
	postReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SQUARE_API_KEY")))

	// Do the POST request
	client := http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error calling square"), StatusCode: resp.StatusCode}, err
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error parsing square response"), StatusCode: 500}, err
	}

	body, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error parsing payload"), StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

type requestBody struct {
	IdempotencyKey string  `json:"idempotency_key,omitempty"`
	Booking        booking `json:"booking,omitempty"`
}

type booking struct {
	Version             int        `json:"version,omitempty"`
	StartAt             time.Time  `json:"start_at,omitempty"`
	LocationID          string     `json:"location_id,omitempty"`
	CustomerID          string     `json:"customer_id,omitempty"`
	CustomerNote        string     `json:"customer_note,omitempty"`
	SellerNote          string     `json:"seller_note,omitempty"`
	AppointmentSegments []segments `json:"appointment_segments,omitempty"`
	LocationType        string     `json:"location_type,omitempty"`
}

type segments struct {
	DurationMinutes         int    `json:"duration_minutes,omitempty"`
	ServiceVariationID      string `json:"service_variation_id,omitempty"`
	TeamMemberID            string `json:"team_member_id,omitempty"`
	ServiceVariationVersion int64  `json:"service_variation_version,omitempty"`
}

type response struct {
	Booking struct {
		ID                  string    `json:"id,omitempty"`
		Version             int       `json:"version,omitempty"`
		Status              string    `json:"status,omitempty"`
		CreatedAt           time.Time `json:"created_at,omitempty"`
		UpdatedAt           time.Time `json:"updated_at,omitempty"`
		LocationID          string    `json:"location_id,omitempty"`
		CustomerID          string    `json:"customer_id,omitempty"`
		CustomerNote        string    `json:"customer_note,omitempty"`
		StartAt             time.Time `json:"start_at,omitempty"`
		AllDay              bool      `json:"all_day,omitempty"`
		AppointmentSegments []struct {
			DurationMinutes         int    `json:"duration_minutes,omitempty"`
			ServiceVariationID      string `json:"service_variation_id,omitempty"`
			TeamMemberID            string `json:"team_member_id,omitempty"`
			ServiceVariationVersion int64  `json:"service_variation_version,omitempty"`
			AnyTeamMember           bool   `json:"any_team_member,omitempty"`
			IntermissionMinutes     int    `json:"intermission_minutes,omitempty"`
		} `json:"appointment_segments,omitempty"`
		SellerNote            string `json:"seller_note,omitempty"`
		TransitionTimeMinutes int    `json:"transition_time_minutes,omitempty"`
		CreatorDetails        struct {
			CreatorType  string `json:"creator_type,omitempty"`
			TeamMemberID string `json:"team_member_id,omitempty"`
		} `json:"creator_details,omitempty"`
		Source       string `json:"source,omitempty"`
		LocationType string `json:"location_type,omitempty"`
	} `json:"booking,omitempty"`
	Errors []any `json:"errors,omitempty"`
}

func main() {
	lambda.Start(handleRequest)
}
