package handler

import (
	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase"
)

type Handlers struct {
	AuthUC *usecase.AuthUsecase

	GenerateScheduleUC  *usecase.GenerateScheduleUsecase
	ScheduleAssignUC    *usecase.ScheduleAssignUsecase
	ScheduleBypassNFCUC *usecase.ScheduleBypassNFCUsecase

	InspectionUC         *usecase.InspectionUsecase
	InspectionApprovalUC *usecase.InspectionApprovalUsecase
	ScheduleApprovalUC   *usecase.ScheduleApprovalUsecase

	SchedulePlanUC   *usecase.SchedulePlanUsecase
	GenerateUC       *usecase.GenerateScheduleUsecase
	UserManagementUC *usecase.UserManagementUsecase
	AuditUC          *usecase.AuditTrailUsecase
	AHUUC            *usecase.AHUUsecase
	bypassUC         *usecase.ScheduleBypassNFCUsecase
	assignUC         *usecase.ScheduleAssignUsecase

	createFormTemplateUsecase           *usecase.CreateFormTemplateUsecase
	getFormTemplateDetailUsecase        *usecase.GetFormTemplateDetailUsecase
	listFormTemplateUsecase             *usecase.ListFormTemplateUsecase
	setFormTemplateActiveUsecase        *usecase.SetFormTemplateActiveUsecase
	createNewFormTemplateVersionUsecase *usecase.CreateNewFormTemplateVersionUsecase
	listFormTemplateVersionsUsecase     *usecase.ListFormTemplateVersionsUsecase
	compareFormTemplateUsecase          *usecase.CompareFormTemplateUsecase

	formRepo repository.FormRepository
}

func NewHandlers(
	authUC *usecase.AuthUsecase,
	userManagementUC *usecase.UserManagementUsecase,
	schedulePlanUC *usecase.SchedulePlanUsecase,
	generateUC *usecase.GenerateScheduleUsecase,
	inspectionUC *usecase.InspectionUsecase,
	inspectionApprovalUC *usecase.InspectionApprovalUsecase,
	scheduleApprovalUC *usecase.ScheduleApprovalUsecase,
	auditUC *usecase.AuditTrailUsecase,
	ahuUC *usecase.AHUUsecase,

	createFormTemplateUC *usecase.CreateFormTemplateUsecase,
	getFormTemplateDetailUC *usecase.GetFormTemplateDetailUsecase,
	listFormTemplateUC *usecase.ListFormTemplateUsecase, // ✅
	setFormTemplateActiveUC *usecase.SetFormTemplateActiveUsecase,
	createNewFormTemplateVersionUC *usecase.CreateNewFormTemplateVersionUsecase,
	listFormTemplateVersionsUC *usecase.ListFormTemplateVersionsUsecase,
	compareFormTemplateUC *usecase.CompareFormTemplateUsecase,

	formRepo repository.FormRepository,
) *Handlers {
	return &Handlers{
		AuthUC:               authUC,
		UserManagementUC:     userManagementUC,
		SchedulePlanUC:       schedulePlanUC,
		GenerateUC:           generateUC,
		AuditUC:              auditUC,
		AHUUC:                ahuUC,
		InspectionApprovalUC: inspectionApprovalUC,
		ScheduleApprovalUC:   scheduleApprovalUC,

		createFormTemplateUsecase:           createFormTemplateUC,
		getFormTemplateDetailUsecase:        getFormTemplateDetailUC,
		listFormTemplateUsecase:             listFormTemplateUC,
		setFormTemplateActiveUsecase:        setFormTemplateActiveUC,
		createNewFormTemplateVersionUsecase: createNewFormTemplateVersionUC,
		listFormTemplateVersionsUsecase:     listFormTemplateVersionsUC,
		compareFormTemplateUsecase:          compareFormTemplateUC,

		formRepo: formRepo,
	}
}
