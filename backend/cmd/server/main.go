package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ahu-backend/internal/database"
	"ahu-backend/internal/handler"
	"ahu-backend/internal/repository/postgres"
	"ahu-backend/internal/usecase"
)

func main() {
    // 1. Force aplikasi Go pakai WIB secara internal
    loc, _ := time.LoadLocation("Asia/Jakarta")
    time.Local = loc

    // 2. Ambil DSN dari ENV
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        // Tambahkan &timezone=Asia/Jakarta di akhir string
        dsn = "postgres://postgres:postgres@localhost:5432/ahu_db?sslmode=disable&timezone=Asia/Jakarta"
    }

    db := database.NewPostgresPool(dsn)

	// ================= REPOSITORY =================
	schedulePlanRepo := postgres.NewSchedulePlanPostgresRepository(db)
	buildingRepo := postgres.NewBuildingPostgresRepository(db)
	areaRepo := postgres.NewAreaPostgresRepository(db)
	ahuRepo := postgres.NewAHUPostgresRepository(db)
	scheduleRepo := postgres.NewSchedulePostgresRepository(db)
	inspectionRepo := postgres.NewInspectionPostgresRepository(db)
	auditRepo := postgres.NewAuditTrailPostgresRepository(db)
	scheduleApprovalRepo := postgres.NewScheduleApprovalPostgres(db)
	userRepo := postgres.NewUserPostgresRepository(db)
	formRepo := postgres.NewFormPostgresRepository(db)
	inspectionResultRepo := postgres.NewInspectionResultPostgresRepository(db)

	// ================= USECASE CORE =================
	auditUC := usecase.NewAuditTrailUsecase(auditRepo)
	authUC := usecase.NewAuthUsecase(userRepo, auditUC)

	schedulePlanUC := usecase.NewSchedulePlanUsecase(schedulePlanRepo, auditUC)

	buildingUC := usecase.NewBuildingUsecase(buildingRepo, auditUC)
	areaUC := usecase.NewAreaUsecase(areaRepo, auditUC)

	ahuUC := usecase.NewAHUUsecase(ahuRepo, auditUC)

	userManagementUC := usecase.NewUserManagementUsecase(
		userRepo,
		auditRepo,
	)

	// ================= SCHEDULE =================
	scheduleQueryUC := usecase.NewScheduleQueryUsecase(scheduleRepo)

	assignScheduleUC := usecase.NewScheduleAssignUsecase(
		scheduleRepo,
		auditRepo,
	)

	generateScheduleUC := usecase.NewGenerateScheduleUsecase(
		schedulePlanRepo,
		scheduleRepo,
	)

	createFilterScheduleUC := usecase.NewCreateFilterScheduleUsecase(
		scheduleRepo,
		schedulePlanRepo,
	)

	bypassUC := usecase.NewScheduleBypassNFCUsecase(
		scheduleRepo,
		auditRepo,
	)

	// ================= PDF =================
	pdfService := usecase.NewSchedulePDFService(scheduleRepo)
	inspectionPDFService := usecase.NewInspectionPDFService(inspectionRepo)

	// ================= FORM TEMPLATE =================
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

	// ================= INSPECTION =================
	inspectionTaskUC := usecase.NewInspectionTaskUsecase(scheduleRepo)

	inspectionQueryUC := usecase.NewInspectionQueryUsecase(
		inspectionRepo,
		userRepo,
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

	signUC := usecase.NewSignInspectionUsecase(
		inspectionRepo,
		inspectionPDFService,
	)

	approveInspectionUsecase := usecase.NewApproveInspectionUsecase(
		inspectionRepo,
	)

	inspectionApprovalUC := usecase.NewInspectionApprovalUsecase(
		inspectionRepo,
		scheduleRepo,
		auditRepo,
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
		scheduleRepo,
		inspectionPDFService,
	)

	// ================= DASHBOARD =================
	dashboardUC := usecase.NewDashboardUsecase(
		scheduleRepo,
		inspectionRepo,
	)

	// ================= APPROVAL =================
	scheduleApprovalUC := usecase.NewScheduleApprovalUsecase(
		scheduleApprovalRepo,
		pdfService,
	)

	// ================= HANDLER =================
	schedulePlanHandler := handler.NewSchedulePlanHandler(schedulePlanUC)

	buildingHandler := handler.NewBuildingHandler(buildingUC)

	areaHandler := handler.NewAreaHandler(areaUC)

	inspectionHandler := handler.NewInspectionHandler(
		inspectionUC,
		inspectionQueryUC,
		scanNFCUsecase,
		inspectionTaskUC,
		signUC,
		approveInspectionUsecase,
		inspectionPDFService,
	)

	scheduleHandler := handler.NewScheduleHandler(
		bypassUC,
		assignScheduleUC,
		scheduleQueryUC,
		scheduleApprovalUC,
	)

	inspectionFormHandler := handler.NewInspectionFormHandler(
		getFormUC,
		submitInspectionFormUC,
	)

	handlers := handler.NewHandlers(
		authUC,
		userManagementUC,
		schedulePlanUC,
		createFilterScheduleUC,
		generateScheduleUC,
		inspectionUC,
		inspectionApprovalUC,
		scheduleApprovalUC,
		dashboardUC,
		auditUC,
		ahuUC,
		createFormTemplateUC,
		getFormTemplateDetailUC,
		listFormTemplateUC,
		setFormTemplateActiveUC,
		createNewFormTemplateVersionUC,
		listFormTemplateVersionsUC,
		compareFormTemplateUC,
		formRepo,
	)

	// ================= ROUTER =================
	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20

	r.Static("/uploads", "./uploads")
	r.Static("/api/files", "./files")
	r.Static("/files", "./files")

	r.LoadHTMLGlob("templates/*")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://10.9.118.16:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
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

	// ================= RUN =================
	r.Run("0.0.0.0:8080")
}