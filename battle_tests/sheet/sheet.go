package sheet

import (
	"context"
	"fmt"
	"time"

	"github.com/bearer/curio/battle_tests/config"
	battletests "github.com/bearer/curio/battle_tests/config"
	metricsscan "github.com/bearer/curio/battle_tests/metrics_scan"
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/rs/zerolog/log"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var SpreadSheetId = int64(1)
var MaxStringLength = 35000

type GoogleSheets struct {
	sheet *sheets.Service
	drive *drive.Service
}

type datatypeHolder struct {
	name        string
	occurrences *float64
}

func New() *GoogleSheets {
	ctx := context.Background()
	client := battletests.Runtime.Sheets.AppConfig.Client(
		context.Background(),
		battletests.Runtime.Sheets.UserToken,
	)

	sheet, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Err(fmt.Errorf("unable to retrieve Docs client: %w", err)).Send()
	}

	drives, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Err(fmt.Errorf("unable to retrieve drive client %w", err)).Send()
	}

	return &GoogleSheets{
		sheet: sheet,
		drive: drives,
	}
}

type Document struct {
	ID string
}

func (client *GoogleSheets) CreateDocument(tagName string, parentFolderId string) *Document {
	var spreadSheet *sheets.Spreadsheet
	var counter = 0
	for {
		counter++
		req := client.sheet.Spreadsheets.Create(&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: "Battle Report - " + tagName,
			},
			Sheets: []*sheets.Sheet{
				{
					Properties: &sheets.SheetProperties{
						Title:   "Full Report",
						SheetId: SpreadSheetId,
						GridProperties: &sheets.GridProperties{
							FrozenRowCount: 1,
						},
					},
				},
			},
		})

		var err error
		spreadSheet, err = req.Do()
		if err != nil {
			log.Err(fmt.Errorf("failed to create document %w", err)).Send()
			time.Sleep(1 * time.Second)
			if counter >= config.Runtime.MaxAttempt {
				return nil
			}
			continue
		}

		break
	}

	counter = 0
	for {
		columnHeaders := []string{
			"Repo URL",
			"Repo Size (KB)",
			"Scan timing (seconds)",
			"Memory consumption",
			"Number of Data Types",
			"Number of Line of Code",
		}

		var headerCells []*sheets.CellData
		for i := 0; i < len(columnHeaders); i++ {
			headerCells = append(headerCells, &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					StringValue: &columnHeaders[i],
				},
			})
		}

		data_types := db.Default().DataTypes
		for index := range data_types {
			headerCells = append(headerCells, &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					StringValue: &data_types[index].DataCategoryName,
				},
			})
		}

		call := client.sheet.Spreadsheets.BatchUpdate(spreadSheet.SpreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					AppendDimension: &sheets.AppendDimensionRequest{
						SheetId:   SpreadSheetId, // sheet id
						Dimension: "COLUMNS",
						Length:    int64(len(data_types)),
					},
				},
			},
		})

		_, err := call.Do()
		if err != nil {
			log.Err(fmt.Errorf("failed to update google sheets with dimensions %w", err)).Send()
			time.Sleep(1 * time.Second)
			if counter >= config.Runtime.MaxAttempt {
				return nil
			}
			continue
		}

		call = client.sheet.Spreadsheets.BatchUpdate(spreadSheet.SpreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					AppendCells: &sheets.AppendCellsRequest{
						SheetId: SpreadSheetId, // sheet id
						Fields:  "*",
						Rows: []*sheets.RowData{
							{
								Values: headerCells,
							},
						},
					},
				},
			},
		})

		_, err = call.Do()
		if err != nil {
			log.Err(fmt.Errorf("failed to update google sheets with headers %w", err)).Send()
			time.Sleep(1 * time.Second)
			if counter >= config.Runtime.MaxAttempt {
				return nil
			}
			continue
		}

		break
	}

	counter = 0
	for {
		driveCall := client.drive.Files.Update(spreadSheet.SpreadsheetId, &drive.File{}).AddParents(parentFolderId).SupportsAllDrives(true)

		_, err := driveCall.Do()
		if err != nil {
			log.Err(fmt.Errorf("failed to update google file parent folder %w", err)).Send()
			time.Sleep(1 * time.Second)
			if counter >= config.Runtime.MaxAttempt {
				return nil
			}
			continue
		}

		break
	}

	return &Document{
		ID: spreadSheet.SpreadsheetId,
	}
}

// inserts metrics with exponential backoff

func (client *GoogleSheets) InsertMetricsMustPass(documentID string, metrics *metricsscan.MetricsReport) {
	for {
		err := client.InsertMetrics(documentID, metrics)
		if err == nil {
			break
		}

		time.Sleep(10 * time.Second)
	}
}

func (client *GoogleSheets) InsertMetrics(documentID string, metrics *metricsscan.MetricsReport) error {
	values := []*sheets.CellData{
		{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &metrics.URL,
			},
		},
		{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: &metrics.RepoSizeKB,
			},
		},
		{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: &metrics.Time,
			},
		},
		{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: &metrics.Memory,
			},
		},
		{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: &metrics.NumberOfDataTypes,
			},
		},
		{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: &metrics.NumberOfLineOfCode,
			},
		},
	}

	holder := make(map[string]datatypeHolder)
	for index, data_type := range metrics.DataTypes {
		if _, exists := holder[data_type.Name]; !exists {
			holder[data_type.Name] = datatypeHolder{
				name:        data_type.Name,
				occurrences: &metrics.DataTypes[index].Occurrences,
			}
		}
	}

	empty_value := float64(0)
	for _, data_type := range db.Default().DataTypes {
		if _, exists := holder[data_type.DataCategoryName]; exists {
			values = append(values, &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					NumberValue: holder[data_type.DataCategoryName].occurrences,
				},
			})
		} else {
			values = append(values, &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					NumberValue: &empty_value,
				},
			})
		}
	}

	call := client.sheet.Spreadsheets.BatchUpdate(documentID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AppendCells: &sheets.AppendCellsRequest{
					SheetId: SpreadSheetId, // sheet id
					Fields:  "*",
					Rows: []*sheets.RowData{
						{
							Values: values,
						},
					},
				},
			},
		},
	})

	_, err := call.Do()
	if err != nil {
		log.Printf("error inserting metrics %s", err)
		return err
	}

	return nil
}
