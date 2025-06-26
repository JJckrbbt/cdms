package api

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/db"
)

// derefString safely dereferences a string pointer and returns its value.
// If the pointer is nil, it returns an empty string.
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// derefStringWithDefault safely dereferences a string pointer.
// If the pointer is nil, it returns the provided default value.
func derefStringWithDefault(s *string, defaultValue string) string {
	if s != nil {
		return *s
	}
	return defaultValue
}

// --- Helper functions to build the update params ---

func buildUserUpdateParams(req *UserUpdateChargebackRequest, existing *db.Chargeback) db.UserUpdateChargebackParams {
	params := db.UserUpdateChargebackParams{
		ID:                     existing.ID,
		CurrentStatus:          existing.CurrentStatus,
		ReasonCode:             existing.ReasonCode,
		Action:                 existing.Action,
		IssueInResearchDate:    existing.IssueInResearchDate,
		AlcToRebill:            existing.AlcToRebill,
		TasToRebill:            existing.TasToRebill,
		LineOfAccountingRebill: existing.LineOfAccountingRebill,
		SpecialInstruction:     existing.SpecialInstruction,
		PassedToPsf:            existing.PassedToPsf,
	}
	// Merge updates from request
	if req.CurrentStatus != nil {
		params.CurrentStatus = db.ChargebackStatus(*req.CurrentStatus)
	}
	if req.ReasonCode != nil {
		params.ReasonCode = db.NullChargebackReasonCode{ChargebackReasonCode: db.ChargebackReasonCode(*req.ReasonCode), Valid: true}
	}
	if req.Action != nil {
		params.Action = db.NullChargebackAction{ChargebackAction: db.ChargebackAction(*req.Action), Valid: true}
	}
	if req.IssueInResearchDate != nil {
		params.IssueInResearchDate = parseDateToPG(*req.IssueInResearchDate)
	}
	if req.ALCToRebill != nil {
		params.AlcToRebill = pgtype.Text{String: *req.ALCToRebill, Valid: true}
	}
	if req.TASToRebill != nil {
		params.TasToRebill = pgtype.Text{String: *req.TASToRebill, Valid: true}
	}
	if req.LineOfAccountingRebill != nil {
		params.LineOfAccountingRebill = pgtype.Text{String: *req.LineOfAccountingRebill, Valid: true}
	}
	if req.SpecialInstruction != nil {
		params.SpecialInstruction = pgtype.Text{String: *req.SpecialInstruction, Valid: true}
	}
	if req.PassedToPSF != nil {
		params.PassedToPsf = parseDateToPG(*req.PassedToPSF)
	}

	return params
}

func buildPFSUpdateParams(req *PFSUpdateChargebackRequest, existing *db.Chargeback) db.PFSUpdateChargebackParams {
	params := db.PFSUpdateChargebackParams{
		ID:                 existing.ID,
		CurrentStatus:      existing.CurrentStatus,
		PassedToPsf:        existing.PassedToPsf,
		NewIpacDocumentRef: existing.NewIpacDocumentRef,
		PfsCompletionDate:  existing.PfsCompletionDate,
	}
	// Merge updates
	if req.CurrentStatus != nil {
		params.CurrentStatus = db.ChargebackStatus(*req.CurrentStatus)
	}
	if req.PassedToPSF != nil {
		params.PassedToPsf = parseDateToPG(*req.PassedToPSF)
	}
	if req.NewIPACDocumentRef != nil {
		params.NewIpacDocumentRef = pgtype.Text{String: *req.NewIPACDocumentRef, Valid: true}
	}
	if req.PFSCompletionDate != nil {
		params.PfsCompletionDate = parseDateToPG(*req.PFSCompletionDate)
	}

	return params
}

func buildAdminUpdateParams(req *AdminUpdateChargebackRequest, existing *db.Chargeback) db.AdminUpdateChargebackParams {
	params := db.AdminUpdateChargebackParams{
		ID:                     existing.ID,
		CurrentStatus:          existing.CurrentStatus,
		ReasonCode:             existing.ReasonCode,
		Action:                 existing.Action,
		IssueInResearchDate:    existing.IssueInResearchDate,
		AlcToRebill:            existing.AlcToRebill,
		TasToRebill:            existing.TasToRebill,
		LineOfAccountingRebill: existing.LineOfAccountingRebill,
		SpecialInstruction:     existing.SpecialInstruction,
		PassedToPsf:            existing.PassedToPsf,
		PfsCompletionDate:      existing.PfsCompletionDate,
	}
	// Merge updates
	if req.CurrentStatus != nil {
		params.CurrentStatus = db.ChargebackStatus(*req.CurrentStatus)
	}
	if req.ReasonCode != nil {
		params.ReasonCode = db.NullChargebackReasonCode{ChargebackReasonCode: db.ChargebackReasonCode(*req.ReasonCode), Valid: true}
	}
	if req.Action != nil {
		params.Action = db.NullChargebackAction{ChargebackAction: db.ChargebackAction(*req.Action), Valid: true}
	}
	if req.IssueInResearchDate != nil {
		params.IssueInResearchDate = parseDateToPG(*req.IssueInResearchDate)
	}
	if req.ALCToRebill != nil {
		params.AlcToRebill = pgtype.Text{String: *req.ALCToRebill, Valid: true}
	}
	if req.TASToRebill != nil {
		params.TasToRebill = pgtype.Text{String: *req.TASToRebill, Valid: true}
	}
	if req.LineOfAccountingRebill != nil {
		params.LineOfAccountingRebill = pgtype.Text{String: *req.LineOfAccountingRebill, Valid: true}
	}
	if req.SpecialInstruction != nil {
		params.SpecialInstruction = pgtype.Text{String: *req.SpecialInstruction, Valid: true}
	}
	if req.PassedToPSF != nil {
		params.PassedToPsf = parseDateToPG(*req.PassedToPSF)
	}
	if req.PFSCompletionDate != nil {
		params.PfsCompletionDate = parseDateToPG(*req.PFSCompletionDate)
	}

	return params
}

// Helper for converting date strings to pgtype.Date
func parseDateToPG(dateStr string) pgtype.Date {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}
