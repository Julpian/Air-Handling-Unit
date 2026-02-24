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
	scheduleApprovalRepo := postgres.NewScheduleApprovalPostgres(db)
	pdfService := usecase.NewSchedulePDFService(scheduleRepo)
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
		ahuRepo,
		scheduleRepo,
	)

	submitInspectionFormUC := usecase.NewSubmitInspectionFormUsecase(
		inspectionResultRepo,
		inspectionRepo,
		formRepo,
		scheduleRepo, // 🔥 TAMBAH
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

	scanNFCUsecase := usecase.NewScanNFCUsecase(
		ahuRepo,
		scheduleRepo,
		inspectionRepo,
		formRepo,
	)

	inspectionHandler := handler.NewInspectionHandler(
		inspectionUC,
		inspectionQueryUC,
		scanNFCUsecase,
	)

	scheduleApprovalUC := usecase.NewScheduleApprovalUsecase(
		scheduleApprovalRepo,
		pdfService,
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
		scheduleApprovalUC,
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
		scheduleApprovalUC, // ✅ TAMBAH
	)

	inspectionFormHandler := handler.NewInspectionFormHandler(
		getFormUC,
		submitInspectionFormUC,
	)

	// ================= ROUTER =================
	r := gin.Default()

	// 🔥 STATIC & LIMIT
	r.MaxMultipartMemory = 8 << 20
	r.Static("/uploads", "./uploads")
	r.Static("/files", "./files")
	r.LoadHTMLGlob("templates/*")

	// 🔥 PASANG CORS DI SINI
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://10.9.118.16:3000", // HP / LAN
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
	}))

	// ================= REGISTER ROUTES =================
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

	// ================= RUN SERVER =================
	r.Run("0.0.0.0:8080")
}
