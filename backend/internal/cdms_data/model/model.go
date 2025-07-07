// internal/cdms_data/model/model.go
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ChargebackReportingSource string
type ChargebackStatus string
type ChargebackReasonCode string
type ChargebackAction string
type ChargebackBusinessLine string
type ChargebackFund string

type NonIpacReportingSource string
type NonIpacStatus string

type Chargeback struct {
	ID                  int64                     `json:"id"`
	ReportingSource     ChargebackReportingSource `json:"reporting_source"`
	Fund                ChargebackFund            `json:"fund"`
	BusinessLine        ChargebackBusinessLine    `json:"business_line"`
	Region              int16                     `json:"region"`
	LocationSystem      *string                   `json:"location_system"`
	Program             string                    `json:"program"`
	ALNum               int16                     `json:"al_num"`
	SourceNum           string                    `json:"source_num"`
	AgreementNum        *string                   `json:"agreement_num"`
	Title               *string                   `json:"title"`
	ALC                 string                    `json:"alc"`
	CustomerTAS         string                    `json:"customer_tas"`
	TaskSubtask         string                    `json:"task_subtask"`
	ClassID             *string                   `json:"class_id"`
	CustomerName        string                    `json:"customer_name"`
	OrgCode             string                    `json:"org_code"`
	DocumentDate        time.Time                 `json:"document_date"`
	AccompDate          time.Time                 `json:"accomp_date"`
	AssignedRebillDRN   *string                   `json:"assigned_rebill_drn"`
	ChargebackAmount    decimal.Decimal           `json:"chargeback_amount"`
	Statement           string                    `json:"statement"`
	BDDocNum            string                    `json:"bd_doc_num"`
	Vendor              string                    `json:"vendor"`
	ArticlesServices    *string                   `json:"articles_services"`
	CurrentStatus       ChargebackStatus          `json:"current_status"`
	IssueInResearchDate *time.Time                `json:"issue_in_research_date"`
	ReasonCode          *ChargebackReasonCode     `json:"reason_code"`
	Action              *ChargebackAction         `json:"action"`
	DaysOld             time.Time                 `json:"days_old"`
	AbsAmount           decimal.Decimal           `json:"abs_amount"`
	AgingCategory       string                    `json:"aging_category"`
	DaysOpenToPFS       int16                     `json:"days_open_to_pfs"`
	DaysPFSToComplete   int16                     `json:"days_pfs_to_complete"`
	DaysComplete        int16                     `json:"days_complete"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
	IsActive            bool                      `json:"is_active"`
}

type NonIpac struct {
	ID                          int64                  `json:"id"`
	ReportingSource             NonIpacReportingSource `json:"reporting_source"`
	BusinessLine                ChargebackBusinessLine `json:"business_line"`
	BilledTotalAmount           decimal.Decimal        `json:"billed_total_amount"`
	PrincipleAmount             decimal.Decimal        `json:"principle_amount"`
	InterestAmount              decimal.Decimal        `json:"interest_amount"`
	PenaltyAmount               decimal.Decimal        `json:"penalty_amount"`
	AdministrationChargesAmount decimal.Decimal        `json:"administration_charges_amount"`
	DebitOutstandingAmount      decimal.Decimal        `json:"debit_outstanding_amount"`
	CreditTotalAmount           decimal.Decimal        `json:"credit_total_amount"`
	CreditOutstandingAmount     decimal.Decimal        `json:"credit_outstanding_amount"`
	Title                       *string                `json:"title"`
	DocumentDate                time.Time              `json:"document_date"`
	AddressCode                 string                 `json:"address_code"`
	Vendor                      string                 `json:"vendor"`
	DebtAppealForbearance       bool                   `json:"debt_appeal_forbearance"`
	Statement                   string                 `json:"statement"`
	DocumentNumber              string                 `json:"document_number"`
	VendorCode                  string                 `json:"vendor_code"`
	CollectionDueDate           time.Time              `json:"collection_due_date"`
	CurrentStatus               *NonIpacStatus         `json:"current_status"`
	PFSPoc                      *int64                 `json:"pfs_poc"`
	GSAPoc                      *int64                 `json:"gsa_poc"`
	CustomerPoc                 *int64                 `json:"customer_poc"`
	PFSContacts                 int16                  `json:"pfs_contacts"`
	OpenDate                    time.Time              `json:"open_date"`
	ReconciledDate              *time.Time             `json:"reconciled_date"`
	DaysOld                     time.Time              `json:"days_old"`
	AgingCategory               string                 `json:"aging_category"`
	AbsAmount                   decimal.Decimal        `json:"abs_amount"`
	CreatedAt                   time.Time              `json:"created_at"`
	UpdatedAt                   time.Time              `json:"updated_at"`
	IsActive                    bool                   `json:"is_active"`
}

type AgencyBureau struct {
	Agency     string    `json:"agency"`
	BureauCode string    `json:"bureau_code"`
	VendorCode string    `json:"vendor_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RemovedRow struct {
	ID               uuid.UUID `json:"id"`
	UploadID         uuid.UUID `json:"upload_id"`
	Timestamp        time.Time `json:"timestamp"`
	ReportType       string    `json:"report_type"`
	OriginalRowData  string    `json:"original_row_data"`
	ReasonForRemoval string    `json:"reason_for_removal"`
}

type Upload struct {
	ID                uuid.UUID  `json:"id"`
	StorageKey        string     `json:"storage_key"`
	Filename          string     `json:"filename"`
	ReportType        string     `json:"report_type"`
	Status            string     `json:"status"`
	UploadedAt        time.Time  `json:"uploaded_at"`
	ProcessedAt       *time.Time `json:"processed_at"`
	ErrorDetails      *string    `json:"error_details"`
	ProcessedByUserID *uuid.UUID `json:"processed_by_user_id"`
}
