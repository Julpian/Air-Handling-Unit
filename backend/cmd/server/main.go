package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ahu-backend/internal/database"
	"ahu-backend/internal/handler"
	"ahu-backend/internal/repository/postgres"
	"ahu-backend/internal/usecase"
)

func main() {
	// ================= DB =================
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/ahu_db?sslmode=disable"
	}
	db := database.NewPostgresPool(dsn)

	// ================= REPOSITORY =================
	// ================= REPOSITORY =================
	schedulePlanRepo := postgres.NewSchedulePlanPostgresRepository(db)

	// ================= USECASE =================
	schedulePlanUC := usecase.NewSchedulePlanUsecase(schedulePlanRepo)

	// ================= HANDLER =================
	schedulePlanHandler := handler.NewSchedulePlanHandler(schedulePlanUC)
	buildingRepo := postgres.NewBuildingPostgresRepository(db)
	buildingUC := usecase.NewBuildingUsecase(buildingRepo)
	areaRepo := postgres.NewAreaPostgresRepository(db)
	areaUC := usecase.NewAreaUsecase(areaRepo)
	areaHandler := handler.NewAreaHandler(areaUC)
	buildingHandler := handler.NewBuildingHandler(buildingUC)
	ahuRepo := postgres.NewAHUPostgresRepository(db)
	scheduleRepo := postgres.NewSchedulePostgresRepository(db)
	inspectionRepo := postgres.NewInspectionPostgresRepository(db)
	auditRepo := postgres.NewAuditTrailPostgresRepository(db)
	userRepo := postgres.NewUserPostgresRepository(db)
	scheduleQueryUC := usecase.NewScheduleQueryUsecase(scheduleRepo)
	ahuUC := usecase.NewAHUUsecase(ahuRepo)
	formRepo := postgres.NewFormPostgresRepository(db)
	inspectionResultRepo := postgres.NewInspectionResultPostgresRepository(db)
	createFormTemplateUC := usecase.NewCreateFormTemplateUsecase(formRepo)
	getFormTemplateDetailUC := usecase.NewGetFormTemplateDetailUsecase(formRepo)
	listFormTemplateUC := usecase.NewListFormTemplateUsecase(formRepo)
	setFormTemplateActiveUC := usecase.NewSetFormTemplateActiveUsecase(formRepo)
	createNewFormTemplateVersionUC :=
		usecase.NewCreateNewFormTemplateVersionUsecase(formRepo)
	listFormTemplateVersionsUC :=
		usecase.NewListFormTemplateVersionsUsecase(formRepo)
	compareFormTemplateUC :=
		usecase.NewCompareFormTemplateUsecase(formRepo)

	inspectionQueryUC := usecase.NewInspectionQueryUsecase(
		inspectionRepo,
		userRepo,
	)

	assignScheduleUC := usecase.NewScheduleAssignUsecase(
		scheduleRepo,
		auditRepo,
	)

	// ================= USECASE =================

	generateScheduleUC := usecase.NewGenerateScheduleUsecase(
		schedulePlanRepo,
		scheduleRepo,
	)

	getFormUC := usecase.NewGetFormByInspectionUsecase(
		inspectionRepo,
		formRepo,
	)

	submitInspectionFormUC := usecase.NewSubmitInspectionFormUsecase(
		inspectionResultRepo,
		inspectionRepo,
		formRepo,
	)

	authUC := usecase.NewAuthUsecase(userRepo)
	userManagementUC := usecase.NewUserManagementUsecase(
		userRepo,
		auditRepo,
	)

	inspectionUC := usecase.NewInspectionUsecase(
		ahuRepo,
		scheduleRepo,
		inspectionRepo,
		formRepo,
		auditRepo,
	)

	inspectionHandler := handler.NewInspectionHandler(
		inspectionUC,
		inspectionQueryUC,
	)

	inspectionApprovalUC := usecase.NewInspectionApprovalUsecase(
		inspectionRepo,
		scheduleRepo,
		auditRepo,
	)

	auditUC := usecase.NewAuditTrailUsecase(auditRepo)

	// ================= HANDLER =================
	handlers := handler.NewHandlers(
		authUC,
		userManagementUC,
		schedulePlanUC,
		generateScheduleUC,
		inspectionUC,
		inspectionApprovalUC,
		auditUC,
		ahuUC,
		createFormTemplateUC,
		getFormTemplateDetailUC,
		listFormTemplateUC, // ✅ HARUS DI SINI
		setFormTemplateActiveUC,
		createNewFormTemplateVersionUC,
		listFormTemplateVersionsUC,
		compareFormTemplateUC,
		formRepo, // ✅ PALING TERAKHIR
	)

	// bypass NFC handler
	bypassUC := usecase.NewScheduleBypassNFCUsecase(
		scheduleRepo,
		auditRepo,
	)

	scheduleHandler := handler.NewScheduleHandler(
		bypassUC,
		assignScheduleUC,
		scheduleQueryUC,
	)

	inspectionFormHandler := handler.NewInspectionFormHandler(
		getFormUC,
		submitInspectionFormUC,
	)

	// ================= ROUTER =================
	r := gin.Default()

	// 🔥 TAMBAHKAN INI
	r.MaxMultipartMemory = 8 << 20 // 8 MB
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	handler.RegisterRoutes(
		r,
		handlers,
		scheduleHandler,
		buildingHandler,
		areaHandler,
		schedulePlanHandler,
		inspectionHandler,
		inspectionFormHandler,
	)

	r.Run(":8080")
}
