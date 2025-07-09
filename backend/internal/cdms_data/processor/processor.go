package processor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/model"
	"github.com/jjckrbbt/cdms/backend/internal/config"
	"github.com/jjckrbbt/cdms/backend/internal/database"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/shopspring/decimal"
	"google.golang.org/api/option"
)

var alcFormatRegex = regexp.MustCompile(`^[0-9]{4,8}$`)

type ProcessingResult struct {
	Status       string
	Error        error
	RowsRemoved  int
	RowsUpserted int64
}

var ExpectedHeaders1048 = []string{"Fund", "Business Line", "Region", "Location/System", "Program", "Statement", "BD Doc Num", "AL Num", "Source Num", "Agreement Num", "Agreement Line Number", "Title", "ALC", "Customer TAS", "Task/Subtask", "Class ID", "Vendor", "Vendor Name", "Org Code", "Agency", "Bureau Code", "Chargeback Amount", "Doc Date", "Days Old", "Accomp Date", "Assigned Rebill DRN", "Articles or Services"}
var ExpectedHeaders1300 = []string{"Fund", "Business Line", "Region", "Location/System", "Program", "Statement", "BD Doc Num", "AL Num", "Source Num", "Agreement Line Number", "Title", "ALC", "Customer TAS", "Task/Subtask", "Class ID", "Vendor", "Vendor Name", "Org Code", "Agency", "Bureau Code", "Chargeback Amount", "Doc Date", "Days Old", "Accomp Date", "Assigned Rebill DRN", "Articles or Services"}
var ExpectedHeadersOutstandingBills = []string{"G_Inv_IPAC_Indicator", "Business_Application_Type", "Business_Application_Code", "Document_Type", "BD Doc Num", "Billing_Reference_Number", "Statement", "Requester_Servicer_Type", "GTC_Num", "G_Invoicing_Order_Number", "Order_Line_Num", "Order_Schedule_Num", "G_Invoicing_Line_Type", "Chargeback Amount", "Principal_Amount", "Interest_Amount", "Penalty_Amount", "System_Generated_Bill_Reduction_Amount", "Total_Write_Off_Amount", "Administration_Charges_Amount", "Outstanding_Amount", "Credit_Total_Amount", "Credit_Outstanding_Amount", "Title", "Doc Date", "Collection_Due_Date", "Debt_Age_Category", "User_ID", "Vendor", "Address_Code", "Vendor Name", "Business Line", "Debt_Appeal_Forebearance", "Rebill_Flag", "Selected_For_G_Inv_IPAC", "Chargeback_End_Date", "Chargeback_Age"}
var ExpectedHeadersVendorCode = []string{"Vendor Agency Code", "Bureau Code", "Agency Location Code", "Vendor Code", "Vendor Address Code", "Name", "Address Line 1", "Address Line 2", "Address Line 3", "City", "State", "Zip", "Status", "Vendor Type", "Reporting Attribute", "Security Org", "Transmit to VCSS Flag"}

var validChargebackFunds = map[string]bool{
	"F-100": true, "F-201": true, "F-305": true, "F-410": true, "F-501": true,
}
var validBusinessLines = map[string]bool{
	"Procurement": true, "Operations": true, "Research & Dev": true, "IT Services": true,
	"Logistics": true, "Admin": true, "Cars": true, "Rent": true, "Credit": true,
	"Hotels": true, "Grocery": true,
}

type Processor struct {
	db           *database.DBClient
	logger       *slog.Logger
	cfg          *config.Config
	gcsClient    *storage.Client
	gcsBucket    string
	systemUserID int64
}

func NewProcessor(dbClient *database.DBClient, logger *slog.Logger, cfg *config.Config) (*Processor, error) {
	ctx := context.Background()
	var clientOptions []option.ClientOption
	if keyFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); keyFilePath != "" {
		clientOptions = append(clientOptions, option.WithCredentialsFile(keyFilePath))
	}
	gcsClient, err := storage.NewClient(ctx, clientOptions...)
	if err != nil {
		logger.Error("Failed to create GCS client for processor", "error", err)
		return nil, fmt.Errorf("failed to create GCS client for processor: %w", err)
	}

	queries := db.New(dbClient.Pool)
	systemUser, err := queries.GetUserByEmail(ctx, "system@cdms.local")
	if err != nil {
		logger.Error("CRITICAL: Failed to fetch system user ID on startup", "error", err)
		return nil, fmt.Errorf("failed to fetch system user ID, processor cannot start: %w", err)
	}

	return &Processor{
		db:           dbClient,
		logger:       logger,
		cfg:          cfg,
		gcsClient:    gcsClient,
		gcsBucket:    cfg.GCSBucketName,
		systemUserID: systemUser.ID,
	}, nil
}

func (p *Processor) ProcessFileFromCloudStorage(ctx context.Context, uploadID string, storageKey string, reportType string) *ProcessingResult {
	procLogger := p.logger.With("upload_id", uploadID, "storage_key", storageKey, "report_type", reportType)
	procLogger.InfoContext(ctx, "Starting asynchronous report processing from cloud storage")

	reader, err := p.gcsClient.Bucket(p.gcsBucket).Object(storageKey).NewReader(ctx)
	if err != nil {
		procLogger.ErrorContext(ctx, "Failed to download file from GCS", "error", err)
		return &ProcessingResult{Status: "FAILED_GENERIC", Error: fmt.Errorf("failed to download from GCS: %w", err)}
	}
	defer reader.Close()

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	headers, err := csvReader.Read()
	if err != nil {
		procLogger.ErrorContext(ctx, "Error reading header row", "error", err)
		return &ProcessingResult{Status: "FAILED_INVALID_FORMAT", Error: fmt.Errorf("error reading header row: %w", err)}
	}
	var expectedHeaders []string
	switch reportType {
	case "BC1300":
		expectedHeaders = ExpectedHeaders1300
	case "BC1048":
		expectedHeaders = ExpectedHeaders1048
	case "OUTSTANDING_BILLS":
		expectedHeaders = ExpectedHeadersOutstandingBills
	case "VENDOR_CODE":
		expectedHeaders = ExpectedHeadersVendorCode
	default:
		return &ProcessingResult{Status: "FAILED_GENERIC", Error: fmt.Errorf("unknown report type: %s", reportType)}
	}
	if !areHeadersValid(headers, expectedHeaders) {
		err := fmt.Errorf("header validation failed")
		procLogger.ErrorContext(ctx, err.Error(), "expected", expectedHeaders, "actual", headers)
		return &ProcessingResult{Status: "FAILED_HEADERS_MISMATCH", Error: err}
	}
	procLogger.InfoContext(ctx, "Headers validated successfully")

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(h)] = i
	}

	allRecords, err := csvReader.ReadAll()
	if err != nil {
		procLogger.ErrorContext(ctx, "Failed to read all CSV records", "error", err)
		return &ProcessingResult{Status: "FAILED_INVALID_FORMAT", Error: fmt.Errorf("failed to read all CSV records: %w", err)}
	}

	existingChargebackSources := make(map[string]db.ChargebackReportingSource)

	if reportType == "BC1048" || reportType == "BC1300" {
		bdDocNums := make([]string, 0, len(allRecords))
		bdDocNumIndex, ok := headerMap["BD Doc Num"]
		if !ok {
			return &ProcessingResult{Status: "FAILED_INVALID_FORMAT", Error: errors.New("missing 'BD Doc Num' header")}
		}
		for _, record := range allRecords {
			if len(record) > bdDocNumIndex {
				bdDocNums = append(bdDocNums, record[bdDocNumIndex])
			}
		}

		if len(bdDocNums) > 0 {
			queries := db.New(p.db.Pool)
			existingRows, err := queries.GetChargebackSourcesByBDDocNums(ctx, bdDocNums)
			if err != nil {
				return &ProcessingResult{Status: "FAILED_GENERIC", Error: fmt.Errorf("failed to get existing chargebacks: %w", err)}
			}

			for _, row := range existingRows {
				key := fmt.Sprintf("%s-%d", row.BdDocNum, row.AlNum)
				existingChargebackSources[key] = row.ReportingSource
			}
		}
	}

	processedChargebacks := []model.Chargeback{}
	processedNonIpacs := []model.NonIpac{}
	processedAgencyBureaus := []model.AgencyBureau{}
	removedRows := []model.RemovedRow{}
	processedKeys := make(map[string]bool)

	for _, record := range allRecords {
		switch reportType {
		case "BC1300", "BC1048":
			chargeback, convErr := convertRecordToChargeback(record, headerMap, reportType)
			if convErr != nil {
				removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, fmt.Sprintf("Data conversion/validation error: %v", convErr), reportType))
				continue
			}

			businessKey := fmt.Sprintf("%s-%d", chargeback.BDDocNum, chargeback.ALNum)

			if existingSource, found := existingChargebackSources[businessKey]; found {
				if string(existingSource) != reportType {
					reason := fmt.Sprintf("Conflict: Record exists but belongs to a different report source ('%s')", existingSource)
					removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, reason, reportType))
					continue
				}
			}

			if _, found := processedKeys[businessKey]; found {
				removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, "Duplicate within current report (Chargeback)", reportType))
			} else {
				processedChargebacks = append(processedChargebacks, chargeback)
				processedKeys[businessKey] = true
			}

		case "OUTSTANDING_BILLS":
			nonipac, convErr := convertRecordToNonIpac(record, headerMap, reportType)
			if convErr != nil {
				removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, fmt.Sprintf("Data conversion/validation error: %v", convErr), reportType))
			} else {
				businessKey := nonipac.DocumentNumber
				if _, found := processedKeys[businessKey]; found {
					removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, "Duplicate within current report (NonIpac)", reportType))
				} else {
					processedNonIpacs = append(processedNonIpacs, nonipac)
					processedKeys[businessKey] = true
				}
			}
		case "VENDOR_CODE":
			agencyBureau, convErr := convertRecordToAgencyBureau(record, headerMap)
			if convErr != nil {
				removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, fmt.Sprintf("Data conversion/validation error: %v", convErr), reportType))
			} else {
				businessKey := agencyBureau.VendorCode
				if _, found := processedKeys[businessKey]; found {
					removedRows = append(removedRows, createRemovedRowEntry(uploadID, record, "Duplicate within current report (AgencyBureau)", reportType))
				} else {
					processedAgencyBureaus = append(processedAgencyBureaus, agencyBureau)
					processedKeys[businessKey] = true
				}
			}
		}
	}

	rowsUpserted, err := p.executeMergeTransaction(ctx, uploadID, reportType, removedRows, processedChargebacks, processedNonIpacs, processedAgencyBureaus)
	if err != nil {
		procLogger.ErrorContext(ctx, "Failed to execute database merge transaction", "error", err)
		return &ProcessingResult{Status: "FAILED_GENERIC", Error: err}
	}

	status := "COMPLETE"
	if len(removedRows) > 0 {
		status = "COMPLETE_WITH_ISSUES"
	}

	return &ProcessingResult{
		Status:       status,
		RowsRemoved:  len(removedRows),
		RowsUpserted: rowsUpserted,
	}
}

func (p *Processor) executeMergeTransaction(ctx context.Context, uploadID string, reportType string, removedRows []model.RemovedRow, chargebacks []model.Chargeback, nonipacs []model.NonIpac, agencyBureaus []model.AgencyBureau) (int64, error) {
	tx, err := p.db.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin pgx transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", p.systemUserID)
	_, err = tx.Exec(ctx, setQuery)
	if err != nil {
		return 0, fmt.Errorf("failed to set user for transaction: %w", err)
	}

	q := db.New(tx)
	var rowsAffected int64

	switch reportType {
	case "BC1300", "BC1048":
		if len(chargebacks) > 0 {
			_, err := tx.Exec(ctx, "CREATE TEMPORARY TABLE temp_chargeback_staging (LIKE chargeback INCLUDING DEFAULTS) ON COMMIT DROP;")
			if err != nil {
				return 0, fmt.Errorf("failed to create temp chargeback table: %w", err)
			}
			if err := q.DeactivateChargebacksBySource(ctx, db.ChargebackReportingSource(reportType)); err != nil {
				return 0, err
			}
			columnNames := []string{"reporting_source", "fund", "business_line", "region", "location_system", "program", "al_num", "source_num", "agreement_num", "title", "alc", "customer_tas", "task_subtask", "class_id", "customer_name", "org_code", "document_date", "accomp_date", "assigned_rebill_drn", "chargeback_amount", "statement", "bd_doc_num", "vendor", "articles_services", "current_status", "reason_code", "action", "is_active"}
			_, err = tx.CopyFrom(ctx, pgx.Identifier{"temp_chargeback_staging"}, columnNames, newChargebackCopySource(chargebacks))
			if err != nil {
				return 0, fmt.Errorf("failed to stage chargebacks: %w", err)
			}
			rowsAffected, err = q.UpsertChargebacks(ctx)
			if err != nil {
				return 0, fmt.Errorf("failed to upsert chargebacks: %w", err)
			}
		}
	case "OUTSTANDING_BILLS":
		if len(nonipacs) > 0 {
			_, err := tx.Exec(ctx, `CREATE TEMPORARY TABLE temp_nonipac_staging (LIKE "nonipac" INCLUDING DEFAULTS) ON COMMIT DROP;`)
			if err != nil {
				return 0, fmt.Errorf("failed to create temp non-ipac table: %w", err)
			}
			if err := q.DeactivateNonIpacsBySource(ctx, db.NonipacReportingSource(reportType)); err != nil {
				return 0, err
			}
			columnNames := []string{"reporting_source", "business_line", "billed_total_amount", "principle_amount", "interest_amount", "penalty_amount", "administration_charges_amount", "debit_outstanding_amount", "credit_total_amount", "credit_outstanding_amount", "title", "document_date", "address_code", "vendor", "debt_appeal_forbearance", "statement", "document_number", "vendor_code", "collection_due_date", "open_date", "is_active"}
			_, err = tx.CopyFrom(ctx, pgx.Identifier{"temp_nonipac_staging"}, columnNames, newNonIpacCopySource(nonipacs))
			if err != nil {
				return 0, fmt.Errorf("failed to stage non-ipacs: %w", err)
			}
			rowsAffected, err = q.UpsertNonIpacs(ctx)
			if err != nil {
				return 0, fmt.Errorf("failed to upsert non-ipacs: %w", err)
			}
		}
	case "VENDOR_CODE":
		if len(agencyBureaus) > 0 {
			_, err := tx.Exec(ctx, "CREATE TEMPORARY TABLE temp_agency_bureau_staging (LIKE agency_bureau INCLUDING DEFAULTS) ON COMMIT DROP;")
			if err != nil {
				return 0, fmt.Errorf("failed to create temp agency bureau table: %w", err)
			}
			columnNames := []string{"agency", "bureau_code", "vendor_code"}
			_, err = tx.CopyFrom(ctx, pgx.Identifier{"temp_agency_bureau_staging"}, columnNames, newAgencyBureauCopySource(agencyBureaus))
			if err != nil {
				return 0, fmt.Errorf("failed to stage agency bureaus: %w", err)
			}
			rowsAffected, err = q.UpsertAgencyBureaus(ctx)
			if err != nil {
				return 0, fmt.Errorf("failed to upsert agency bureaus: %w", err)
			}
		}
	}

	if len(removedRows) > 0 {
		p.logger.InfoContext(ctx, "Logging removed rows to database", "count", len(removedRows))
		columnNames := []string{"id", "upload_id", "timestamp", "report_type", "original_row_data", "reason_for_removal"}
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"removed_rows_log"}, columnNames, newRemovedRowCopySource(uploadID, removedRows))
		if err != nil {
			p.logger.ErrorContext(ctx, "Critical: failed to log removed rows after successful data merge", "error", err)
		}
	}

	return rowsAffected, tx.Commit(ctx)
}

type chargebackCopySource struct {
	rows []model.Chargeback
	idx  int
}

func newChargebackCopySource(rows []model.Chargeback) *chargebackCopySource {
	return &chargebackCopySource{rows: rows, idx: -1}
}
func (s *chargebackCopySource) Next() bool {
	s.idx++
	return s.idx < len(s.rows)
}
func (s *chargebackCopySource) Values() ([]any, error) {
	row := s.rows[s.idx]
	return []any{row.ReportingSource, row.Fund, row.BusinessLine, row.Region, row.LocationSystem, row.Program, row.ALNum, row.SourceNum, row.AgreementNum, row.Title, row.ALC, row.CustomerTAS, row.TaskSubtask, row.ClassID, row.CustomerName, row.OrgCode, row.DocumentDate, row.AccompDate, row.AssignedRebillDRN, row.ChargebackAmount, row.Statement, row.BDDocNum, row.Vendor, row.ArticlesServices, row.CurrentStatus, row.ReasonCode, row.Action, row.IsActive}, nil
}
func (s *chargebackCopySource) Err() error { return nil }

type nonipacCopySource struct {
	rows []model.NonIpac
	idx  int
}

func newNonIpacCopySource(rows []model.NonIpac) *nonipacCopySource {
	return &nonipacCopySource{rows: rows, idx: -1}
}
func (s *nonipacCopySource) Next() bool {
	s.idx++
	return s.idx < len(s.rows)
}
func (s *nonipacCopySource) Values() ([]any, error) {
	row := s.rows[s.idx]
	return []any{row.ReportingSource, row.BusinessLine, row.BilledTotalAmount, row.PrincipleAmount, row.InterestAmount, row.PenaltyAmount, row.AdministrationChargesAmount, row.DebitOutstandingAmount, row.CreditTotalAmount, row.CreditOutstandingAmount, row.Title, row.DocumentDate, row.AddressCode, row.Vendor, row.DebtAppealForbearance, row.Statement, row.DocumentNumber, row.VendorCode, row.CollectionDueDate, row.OpenDate, row.IsActive}, nil
}
func (s *nonipacCopySource) Err() error { return nil }

type agencyBureauCopySource struct {
	rows []model.AgencyBureau
	idx  int
}

func newAgencyBureauCopySource(rows []model.AgencyBureau) *agencyBureauCopySource {
	return &agencyBureauCopySource{rows: rows, idx: -1}
}
func (s *agencyBureauCopySource) Next() bool {
	s.idx++
	return s.idx < len(s.rows)
}
func (s *agencyBureauCopySource) Values() ([]any, error) {
	row := s.rows[s.idx]
	return []any{row.Agency, row.BureauCode, row.VendorCode}, nil
}
func (s *agencyBureauCopySource) Err() error { return nil }

type removedRowCopySource struct {
	uploadID uuid.UUID
	rows     []model.RemovedRow
	idx      int
}

func newRemovedRowCopySource(uploadIDStr string, rows []model.RemovedRow) *removedRowCopySource {
	uid, _ := uuid.Parse(uploadIDStr)
	return &removedRowCopySource{uploadID: uid, rows: rows, idx: -1}
}
func (s *removedRowCopySource) Next() bool {
	s.idx++
	return s.idx < len(s.rows)
}
func (s *removedRowCopySource) Values() ([]any, error) {
	row := s.rows[s.idx]
	return []any{row.ID, s.uploadID, row.Timestamp, row.ReportType, row.OriginalRowData, row.ReasonForRemoval}, nil
}
func (s *removedRowCopySource) Err() error { return nil }

func areHeadersValid(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if strings.TrimSpace(actual[i]) != strings.TrimSpace(expected[i]) {
			return false
		}
	}
	return true
}

func getString(record []string, headerMap map[string]int, headerName string) (string, bool) {
	if idx, ok := headerMap[headerName]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx]), true
	}
	return "", false
}

func getStringPtr(record []string, headerMap map[string]int, headerName string) *string {
	s, ok := getString(record, headerMap, headerName)
	if ok && s != "" {
		return &s
	}
	return nil
}

func parseInt16(record []string, headerMap map[string]int, headerName string) (int16, error) {
	s, ok := getString(record, headerMap, headerName)
	if !ok || s == "" {
		return 0, fmt.Errorf("missing or empty '%s'", headerName)
	}
	val, err := strconv.ParseInt(strings.ReplaceAll(s, ",", ""), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid '%s' format: %w", headerName, err)
	}
	return int16(val), nil
}

func parseDecimal(record []string, headerMap map[string]int, headerName string) (decimal.Decimal, error) {
	s, ok := getString(record, headerMap, headerName)
	if !ok || s == "" {
		return decimal.Zero, nil
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "$", "")
	val, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("invalid '%s' decimal format: %w", headerName, err)
	}
	return val, nil
}

func parseDate(record []string, headerMap map[string]int, headerName string) (time.Time, error) {
	s, ok := getString(record, headerMap, headerName)
	if !ok || s == "" {
		return time.Time{}, fmt.Errorf("missing or empty '%s'", headerName)
	}
	layouts := []string{"1/2/2006", "2006-01-02", "01-Jan-06", time.RFC3339}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("un-parsable '%s' date format: %s", headerName, s)
}

func parseBool(record []string, headerMap map[string]int, headerName string) (bool, error) {
	s, ok := getString(record, headerMap, headerName)
	if !ok || s == "" {
		return false, nil
	}
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "TRUE" || s == "1" || s == "Y" {
		return true, nil
	}
	if s == "FALSE" || s == "0" || s == "N" {
		return false, nil
	}
	return false, fmt.Errorf("invalid '%s' boolean format: %s", headerName, s)
}

func createRemovedRowEntry(uploadID string, originalData []string, reason string, reportType string) model.RemovedRow {
	uid, _ := uuid.Parse(uploadID)

	jsonData, err := json.Marshal(originalData)
	var rowData string
	if err != nil {
		rowData = `["error marshalling original row"]`
	} else {
		rowData = string(jsonData)
	}

	return model.RemovedRow{
		ID:               uuid.New(),
		UploadID:         uid,
		Timestamp:        time.Now(),
		ReportType:       reportType,
		OriginalRowData:  rowData,
		ReasonForRemoval: reason,
	}
}
func convertRecordToChargeback(record []string, headerMap map[string]int, reportType string) (model.Chargeback, error) {
	var parseErrors []error
	chargeback := model.Chargeback{
		ReportingSource: model.ChargebackReportingSource(reportType),
		IsActive:        true,
	}

	if val, ok := getString(record, headerMap, "Fund"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Fund'"))
	} else if !validChargebackFunds[val] {
		parseErrors = append(parseErrors, fmt.Errorf("invalid 'Fund' value: %s", val))
	} else {
		chargeback.Fund = model.ChargebackFund(val)
	}
	if val, ok := getString(record, headerMap, "Business Line"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Business Line'"))
	} else if !validBusinessLines[val] {
		parseErrors = append(parseErrors, fmt.Errorf("invalid 'Business Line' value: %s", val))
	} else {
		chargeback.BusinessLine = model.ChargebackBusinessLine(val)
	}
	if val, err := parseInt16(record, headerMap, "Region"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		chargeback.Region = val
	}
	chargeback.LocationSystem = getStringPtr(record, headerMap, "Location/System")
	if val, ok := getString(record, headerMap, "Program"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Program'"))
	} else {
		chargeback.Program = val
	}
	if val, err := parseInt16(record, headerMap, "AL Num"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		chargeback.ALNum = val
	}
	if val, ok := getString(record, headerMap, "Source Num"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Source Num'"))
	} else {
		chargeback.SourceNum = val
	}
	chargeback.AgreementNum = getStringPtr(record, headerMap, "Agreement Num")
	chargeback.Title = getStringPtr(record, headerMap, "Title")

	if val, ok := getString(record, headerMap, "ALC"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'ALC'"))
	} else if !alcFormatRegex.MatchString(val) {
		parseErrors = append(parseErrors, fmt.Errorf("value for 'ALC' does not match required format (4-8 digits): '%s'", val))
	} else {
		chargeback.ALC = val
	}
	if val, ok := getString(record, headerMap, "Org Code"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Org Code'"))
	} else if len(val) > 8 {
		parseErrors = append(parseErrors, fmt.Errorf("value for 'Org Code' exceeds 8 characters: '%s'", val))
	} else {
		chargeback.OrgCode = val
	}
	if valPtr := getStringPtr(record, headerMap, "Assigned Rebill DRN"); valPtr != nil && len(*valPtr) > 8 {
		parseErrors = append(parseErrors, fmt.Errorf("value for 'Assigned Rebill DRN' exceeds 8 characters: '%s'", *valPtr))
	} else {
		chargeback.AssignedRebillDRN = valPtr
	}
	if val, ok := getString(record, headerMap, "Statement"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Statement'"))
	} else if len(val) > 8 {
		parseErrors = append(parseErrors, fmt.Errorf("value for 'Statement' exceeds 8 characters: '%s'", val))
	} else {
		chargeback.Statement = val
	}
	if val, ok := getString(record, headerMap, "Vendor"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Vendor'"))
	} else if len(val) > 8 {
		parseErrors = append(parseErrors, fmt.Errorf("value for 'Vendor' exceeds 8 characters: '%s'", val))
	} else {
		chargeback.Vendor = val
	}

	if val, ok := getString(record, headerMap, "Customer TAS"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Customer TAS'"))
	} else {
		chargeback.CustomerTAS = val
	}
	if val, ok := getString(record, headerMap, "Task/Subtask"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Task/Subtask'"))
	} else {
		chargeback.TaskSubtask = val
	}
	chargeback.ClassID = getStringPtr(record, headerMap, "Class ID")
	if val, ok := getString(record, headerMap, "Vendor Name"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Customer Name' (from Vendor Name column)"))
	} else {
		chargeback.CustomerName = val
	}
	if val, err := parseDate(record, headerMap, "Doc Date"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		chargeback.DocumentDate = val
	}
	if val, err := parseDate(record, headerMap, "Accomp Date"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		chargeback.AccompDate = val
	}
	if val, err := parseDecimal(record, headerMap, "Chargeback Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		chargeback.ChargebackAmount = val
	}
	if val, ok := getString(record, headerMap, "BD Doc Num"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'BD Doc Num'"))
	} else {
		chargeback.BDDocNum = val
	}
	chargeback.ArticlesServices = getStringPtr(record, headerMap, "Articles or Services")

	if len(parseErrors) > 0 {
		return model.Chargeback{}, errors.Join(parseErrors...)
	}
	return chargeback, nil
}
func convertRecordToNonIpac(record []string, headerMap map[string]int, reportType string) (model.NonIpac, error) {
	nonipac := model.NonIpac{
		ReportingSource: model.NonIpacReportingSource(reportType),
		IsActive:        true,
	}
	var parseErrors []error

	if val, ok := getString(record, headerMap, "Business Line"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Business Line'"))
	} else if !validBusinessLines[val] {
		parseErrors = append(parseErrors, fmt.Errorf("invalid 'Business Line' value: %s", val))
	} else {
		nonipac.BusinessLine = model.ChargebackBusinessLine(val)
	}
	if val, err := parseDecimal(record, headerMap, "Chargeback Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.BilledTotalAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Principal_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.PrincipleAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Interest_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.InterestAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Penalty_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.PenaltyAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Administration_Charges_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.AdministrationChargesAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Outstanding_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.DebitOutstandingAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Credit_Total_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.CreditTotalAmount = val
	}
	if val, err := parseDecimal(record, headerMap, "Credit_Outstanding_Amount"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.CreditOutstandingAmount = val
	}
	nonipac.Title = getStringPtr(record, headerMap, "Title")
	if val, err := parseDate(record, headerMap, "Doc Date"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.DocumentDate = val
		nonipac.OpenDate = val
	}
	if val, ok := getString(record, headerMap, "Address_Code"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Address_Code'"))
	} else {
		nonipac.AddressCode = val
	}
	if val, ok := getString(record, headerMap, "Vendor Name"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Vendor Name'"))
	} else {
		nonipac.Vendor = val
	}
	if val, err := parseBool(record, headerMap, "Debt_Appeal_Forebearance"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.DebtAppealForbearance = val
	}
	if val, ok := getString(record, headerMap, "Statement"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Statement'"))
	} else {
		nonipac.Statement = val
	}
	if val, ok := getString(record, headerMap, "BD Doc Num"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Document Number' (from BD Doc Num)"))
	} else {
		nonipac.DocumentNumber = val
	}
	if val, ok := getString(record, headerMap, "Vendor"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Vendor'"))
	} else {
		nonipac.VendorCode = val
	}
	if val, err := parseDate(record, headerMap, "Collection_Due_Date"); err != nil {
		parseErrors = append(parseErrors, err)
	} else {
		nonipac.CollectionDueDate = val
	}

	if len(parseErrors) > 0 {
		return model.NonIpac{}, errors.Join(parseErrors...)
	}
	return nonipac, nil
}

func convertRecordToAgencyBureau(record []string, headerMap map[string]int) (model.AgencyBureau, error) {
	agencyBureau := model.AgencyBureau{}
	var parseErrors []error

	if val, ok := getString(record, headerMap, "Vendor Agency Code"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Vendor Agency Code'"))
	} else {
		agencyBureau.Agency = val
	}
	if val, ok := getString(record, headerMap, "Bureau Code"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Bureau Code'"))
	} else {
		agencyBureau.BureauCode = val
	}
	if val, ok := getString(record, headerMap, "Vendor Code"); !ok || val == "" {
		parseErrors = append(parseErrors, errors.New("missing 'Vendor Code'"))
	} else {
		agencyBureau.VendorCode = val
	}

	if len(parseErrors) > 0 {
		return model.AgencyBureau{}, errors.Join(parseErrors...)
	}
	return agencyBureau, nil
}
