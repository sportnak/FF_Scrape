package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var client_id = "571566073293-v12sves6fsh9p14lcq6bvj9rh3qkg1pv.apps.googleusercontent.com"
var client_secrect = "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"

var api_key = "AIzaSyCYNfbReExWxFLxF8ChLIL1olhuThZv4Cg"
var sheet_id string = "1U92n1OVEure98o1TVbuvDMq2PxG09J4uqW0YqfXj4Ks"
var sheet_name string = "'Conference Stats'!"
var start_col string = "A"
var start_write_col string = "E"
var end_col string = "O"

type AthleteRow struct {
	Name         string
	Row          int
	Position     string
	Team         string
	Points       float64
	PASS_YDS     float64
	PASS_TD      float64
	INT          float64
	Receptions   float64
	Receiving    float64
	Receiving_TD float64
	Rushing      float64
	Rushing_TD   float64
	Fumble       float64
	Game_started float64
	Two_pt       float64
}

var sheet map[string]*AthleteRow
var field_positions map[int]string = make(map[int]string)

func UpdatePlayer(athlete *Athlete) {
	row_athlete := sheet[athlete.Name]

	if row_athlete == nil {
		return
	}
	//*/3 * * * * /Users/MichaelNakayama/FantasyFootballScraper/go_version/main
	has_update := false
	if athlete.PASS_YDS != 0 {
		has_update = has_update || row_athlete.PASS_YDS != athlete.PASS_YDS
		row_athlete.PASS_YDS = athlete.PASS_YDS
	}
	if athlete.PASS_TD != 0 {
		has_update = has_update || row_athlete.PASS_TD != athlete.PASS_TD
		row_athlete.PASS_TD = athlete.PASS_TD
	}
	if athlete.INT != 0 {
		has_update = has_update || row_athlete.INT != athlete.INT

		row_athlete.INT = athlete.INT
	}
	if athlete.REC != 0 {
		has_update = has_update || row_athlete.Receptions != athlete.REC

		row_athlete.Receptions = athlete.REC
	}
	if athlete.REC_YDS != 0 {
		has_update = has_update || row_athlete.Receiving != athlete.REC_YDS

		row_athlete.Receiving = athlete.REC_YDS
	}
	if athlete.REC_TD != 0 {
		has_update = has_update || row_athlete.Receiving_TD != athlete.REC_TD

		row_athlete.Receiving_TD = athlete.REC_TD
	}
	if athlete.RUSH_YDS != 0 {
		has_update = has_update || row_athlete.Rushing != athlete.RUSH_YDS

		row_athlete.Rushing = athlete.RUSH_YDS
	}
	if athlete.RUSH_TD != 0 {
		has_update = has_update || row_athlete.Rushing_TD != athlete.RUSH_TD

		row_athlete.Rushing_TD = athlete.RUSH_TD
	}
	if athlete.LOST != 0 {
		has_update = has_update || row_athlete.Fumble != athlete.LOST

		row_athlete.Fumble = athlete.LOST
	}

	// Replace with your Google Sheets ID, range, and API key
	// sheetRange := sheet_name + start_col + "1:" + end_col + "36"

	if !has_update {
		fmt.Println("No changes for ", row_athlete.Name)
		return
	}

	// Create a Google Sheets API client
	sheetsService, err := createSheetsService()
	if err != nil {
		log.Fatalf("Unable to create Sheets API client: %v", err)
	}

	var result []interface{}
	// result := []interface{/}
	for i := 0; i < len(field_positions); i++ {
		field := field_positions[i]
		switch field {
		case "Passing Yards":
			result = append(result, row_athlete.PASS_YDS)
		case "Passing TDs":
			result = append(result, row_athlete.PASS_TD)
		case "INT":
			result = append(result, row_athlete.INT)
		case "Receptions":
			result = append(result, row_athlete.Receptions)
		case "Receiving":
			result = append(result, row_athlete.Receiving)
		case "Receiving TDs":
			result = append(result, row_athlete.Receiving_TD)
		case "Rushing":
			result = append(result, row_athlete.Rushing)
		case "Rushing TDs":
			result = append(result, row_athlete.Rushing_TD)
		case "Fumbles":
			result = append(result, row_athlete.Fumble)
		case "Game Started":
			result = append(result, row_athlete.Game_started)
		case "2pt":
			result = append(result, row_athlete.Two_pt)
		}
	}

	// Update the sheet
	// values := [][]interface{}{
	// 	{"New Value 1", "New Value 2", "New Value 3"},
	// 	// Add more rows as needed
	// }

	// fmt.Println(field_positions)
	// fmt.Println(result)
	fmt.Println("Updating player:", row_athlete.Name)
	var sheetRange = sheet_name + start_write_col + strconv.Itoa(row_athlete.Row) + ":" + end_col + strconv.Itoa(row_athlete.Row)
	updateSheet(sheetsService, sheet_id, sheetRange, [][]interface{}{result})

}

func InitSheet() {
	// Replace with your Google Sheets ID, range, and API key
	sheetRange := sheet_name + start_col + "1:" + end_col + "100"

	// Create a Google Sheets API client
	sheetsService, err := createSheetsService()
	if err != nil {
		log.Fatalf("Unable to create Sheets API client: %v", err)
	}

	// Update the sheet
	readSheet(sheetsService, sheet_id, sheetRange)
}

func createSheetsService() (*sheets.Service, error) {
	ctx := context.Background()

	// Use API key for authentication
	// client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to create Google Sheets client: %v", err)
	// }

	// Create a Sheets service
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile("/Users/MichaelNakayama/FantasyFootballScraper/go_version/ff-scrape-72d1f2dd0192.json"))
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	return sheetsService, nil
}

func updateSheet(sheetsService *sheets.Service, sheetID, sheetRange string, values [][]interface{}) {
	data := &sheets.ValueRange{
		Values: values,
	}

	_, err := sheetsService.Spreadsheets.Values.Update(sheetID, sheetRange, data).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Fatalf("Unable to update sheet: %v", err)
	}

	fmt.Println("Sheet updated successfully!")
}

func readSheet(sheetsService *sheets.Service, sheetID, sheetRange string) {
	// Read data from the sheet
	response, err := sheetsService.Spreadsheets.Values.Get(sheetID, sheetRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the rows
	if len(response.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		sheet = make(map[string]*AthleteRow)

		var header []string
		for index, row := range response.Values {
			var data []string
			for _, value := range row {
				data = append(data, fmt.Sprintf("%v", value))
			}

			if index == 0 {
				fmt.Printf("%v\n", header)
				for index, value := range data {
					field_positions[index] = value
				}

				header = data
				continue
			}
			// fmt.Println(row)
			athlete_row := mapRowToStrings(header, data, index+1)
			sheet[athlete_row.Name] = athlete_row // fmt.Printf("%v\n", row)
		}
	}

}

func mapRowToStrings(header []string, row []string, position int) *AthleteRow {
	athlete := &AthleteRow{
		Row: position,
	}
	for index, header_name := range header {
		if len(row) == index {
			return athlete
		}
		var err error
		switch header_name {
		case "Position":
			athlete.Position = row[index]
		case "Player":
			athlete.Name = row[index]
		case "Team":
			athlete.Team = row[index]
		case "Points":
			athlete.Points, err = strconv.ParseFloat(row[index], 64)
		case "Passing Yards":
			athlete.PASS_YDS, err = strconv.ParseFloat(row[index], 64)
		case "Passing TDs":
			athlete.PASS_TD, err = strconv.ParseFloat(row[index], 64)
		case "INT":
			athlete.INT, err = strconv.ParseFloat(row[index], 64)
		case "Receptions":
			athlete.Receptions, err = strconv.ParseFloat(row[index], 64)
		case "Receiving":
			athlete.Receiving, err = strconv.ParseFloat(row[index], 64)
		case "Receiving TDs":
			athlete.Receiving_TD, err = strconv.ParseFloat(row[index], 64)
		case "Rushing":
			athlete.Rushing, err = strconv.ParseFloat(row[index], 64)
		case "Rushing TDs":
			athlete.Rushing_TD, err = strconv.ParseFloat(row[index], 64)
		case "Fumbles":
			athlete.Fumble, err = strconv.ParseFloat(row[index], 64)
		case "Game Started":
			athlete.Game_started, err = strconv.ParseFloat(row[index], 64)
		case "2pt":
			athlete.Two_pt, err = strconv.ParseFloat(row[index], 64)
		}
		if err != nil {
			// log.Println(err)
		}
	}

	return athlete
}
