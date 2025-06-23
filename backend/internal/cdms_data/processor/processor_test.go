package processor

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/model"
	"github.com/shopspring/decimal"
)

func TestConvertRecordToChargeback(t *testing.T) {
	// --- Test Setup ---
	// Create a standard header map that all test cases can use.
	headerMap := map[string]int{
		"Fund":                 0,
		"Business Line":        1,
		"Region":               2,
		"Location/System":      3,
		"Program":              4,
		"AL Num":               5,
		"Source Num":           6,
		"Agreement Num":        7,
		"Title":                8,
		"ALC":                  9,
		"Customer TAS":         10,
		"Task/Subtask":         11,
		"Class ID":             12,
		"Vendor Name":          13,
		"Org Code":             14,
		"Doc Date":             15,
		"Accomp Date":          16,
		"Assigned Rebill DRN":  17,
		"Chargeback Amount":    18,
		"Statement":            19,
		"BD Doc Num":           20,
		"Vendor":               21,
		"Articles or Services": 22,
	}

	// Helper to create a valid decimal for the expected struct
	validAmount, _ := decimal.NewFromString("123.45")

	// Helper to create valid time objects
	docDate, _ := time.Parse("2006-01-02", "2025-06-23")
	accompDate, _ := time.Parse("2006-01-02", "2025-06-24")

	// --- Test Cases ---
	testCases := []struct {
		name          string
		inputRecord   []string
		reportType    string
		expectError   bool
		expectedError string
	}{
		{
			name: "Happy Path - Valid BC1048 Record",
			inputRecord: []string{
				"F-100", "IT Services", "10", "LOC", "PROG", "123", "S123", "AGR1", "Test Title", "12345678",
				"TAS123", "TASK1", "C123", "Test Customer", "ORG123", "2025-06-23", "2025-06-24", "DRN123", "123.45",
				"STMT1", "BD123", "VEND123", "Test Services",
			},
			reportType:  "BC1048",
			expectError: false,
		},
		{
			name: "Validation Fail - ALC has letters",
			inputRecord: []string{
				"F-100", "IT Services", "10", "LOC", "PROG", "123", "S123", "AGR1", "Test Title", "ABCDE", // Invalid ALC
				"TAS123", "TASK1", "C123", "Test Customer", "ORG123", "2025-06-23", "2025-06-24", "DRN123", "123.45",
				"STMT1", "BD123", "VEND123", "Test Services",
			},
			reportType:    "BC1048",
			expectError:   true,
			expectedError: "value for 'ALC' does not match required format (4-8 digits): 'ABCDE'",
		},
		{
			name: "Validation Fail - Vendor code too long",
			inputRecord: []string{
				"F-100", "IT Services", "10", "LOC", "PROG", "123", "S123", "AGR1", "Test Title", "12345678",
				"TAS123", "TASK1", "C123", "Test Customer", "ORG123", "2025-06-23", "2025-06-24", "DRN123", "123.45",
				"STMT1", "BD123", "THISISWAYTOOLONG", "Test Services", // Invalid Vendor
			},
			reportType:    "BC1048",
			expectError:   true,
			expectedError: "value for 'Vendor' exceeds 8 characters: 'THISISWAYTOOLONG'",
		},
		{
			name: "Validation Fail - Missing Required Field (Fund)",
			inputRecord: []string{
				"", "IT Services", "10", "LOC", "PROG", "123", "S123", "AGR1", "Test Title", "12345678", // Missing Fund
				"TAS123", "TASK1", "C123", "Test Customer", "ORG123", "2025-06-23", "2025-06-24", "DRN123", "123.45",
				"STMT1", "BD123", "VEND123", "Test Services",
			},
			reportType:    "BC1048",
			expectError:   true,
			expectedError: "missing 'Fund'",
		},
		{
			name: "Validation Fail - Invalid Business Line",
			inputRecord: []string{
				"F-100", "INVALID LINE", "10", "LOC", "PROG", "123", "S123", "AGR1", "Test Title", "12345678",
				"TAS123", "TASK1", "C123", "Test Customer", "ORG123", "2025-06-23", "2025-06-24", "DRN123", "123.45",
				"STMT1", "BD123", "VEND123", "Test Services",
			},
			reportType:    "BC1048",
			expectError:   true,
			expectedError: "invalid 'Business Line' value: INVALID LINE",
		},
	}

	// --- Test Execution ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function we are testing
			result, err := convertRecordToChargeback(tc.inputRecord, headerMap, tc.reportType)

			if tc.expectError {
				// We expected an error
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("expected error message to contain '%s', but got '%v'", tc.expectedError, err)
				}

			} else {
				// We did not expect an error
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}

				// If no error, we can check if the returned struct is correct
				// We'll create the expected struct for the "Happy Path" case
				title := "Test Title"
				loc := "LOC"
				agrNum := "AGR1"
				classID := "C123"
				drn := "DRN123"
				articles := "Test Services"

				expected := model.Chargeback{
					ReportingSource:   "BC1048",
					IsActive:          true,
					CurrentStatus:     "Open",
					Fund:              "F-100",
					BusinessLine:      "IT Services",
					Region:            10,
					LocationSystem:    &loc,
					Program:           "PROG",
					ALNum:             123,
					SourceNum:         "S123",
					AgreementNum:      &agrNum,
					Title:             &title,
					ALC:               "12345678",
					CustomerTAS:       "TAS123",
					TaskSubtask:       "TASK1",
					ClassID:           &classID,
					CustomerName:      "Test Customer",
					OrgCode:           "ORG123",
					DocumentDate:      docDate,
					AccompDate:        accompDate,
					AssignedRebillDRN: &drn,
					ChargebackAmount:  validAmount,
					Statement:         "STMT1",
					BDDocNum:          "BD123",
					Vendor:            "VEND123",
					ArticlesServices:  &articles,
				}

				if !reflect.DeepEqual(result, expected) {
					t.Errorf("struct mismatch: got %+v, want %+v", result, expected)
				}
			}
		})
	}
}
