package handler

import (
	"ahu-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	handlers *Handlers,
	scheduleHandler *ScheduleHandler,
	buildingHandler *BuildingHandler,
	areaHandler *AreaHandler,
	schedulePlanHandler *SchedulePlanHandler,
	inspectionHandler *InspectionHandler,
	inspectionFormHandler *InspectionFormHandler,
) {

	api := r.Group("/api")

	// ================= PUBLIC =================
	api.POST("/auth/login", handlers.Login)

	// ================= AUTHENTICATED =================
	secured := api.Group("")
	secured.Use(middleware.AuthMiddleware())
	{
		// =================================================
		// ================= INSPECTOR =====================
		// =================================================
		inspector := secured.Group("/inspection")
		inspector.Use(middleware.RequireRole("inspector"))
		{
			inspector.GET("", inspectionHandler.List)
			inspector.GET("/dashboard", inspectionHandler.Dashboard)
			inspector.GET("/:inspection_id/detail", inspectionHandler.Detail)
			inspector.POST("/scan-nfc", inspectionHandler.ScanNFC)
			inspector.POST("/:inspection_id/form/submit", inspectionFormHandler.Submit)
			inspector.GET("/:inspection_id/form", inspectionFormHandler.GetForm)
		}

		// =================================================
		// ================= ASMEN =========================
		// =================================================
		asmen := secured.Group("/asmen")
		asmen.Use(middleware.RequireRole("asmen"))
		{
			asmen.POST("/inspection/:id/approve", handlers.ApproveInspection)
			asmen.POST("/inspection/:id/reject", handlers.RejectInspection)
		}

		secured.GET("/profile", handlers.GetMyProfile)
		secured.PATCH("/profile", handlers.UpdateMyProfile)
		secured.PATCH("/profile/password", handlers.ChangeMyPassword)
		secured.POST("/profile/avatar", handlers.UploadAvatar)
		//secured.GET("/ahus/by-nfc/:nfc_uid", handlers.GetAHUByNFC)

		// =================================================
		// ================= ADMIN =========================
		// =================================================
		admin := secured.Group("/admin")
		admin.Use(middleware.RequireAdminLike())
		{
			// ================= SCHEDULE =================
			admin.POST("/schedule-plan", schedulePlanHandler.Create)
			admin.GET("/schedule-plan", schedulePlanHandler.List)
			admin.PATCH("/schedule-plan/:id", handlers.UpdateSchedulePlan)
			admin.DELETE("/schedule-plan/:id", handlers.DeleteSchedulePlan)

			admin.POST("/schedules/generate", handlers.GenerateSchedule)

			admin.GET("/schedules", scheduleHandler.List)

			admin.GET("/inspectors", inspectionHandler.ListDropdown)

			admin.PATCH("/schedules/:id/assign-inspector", scheduleHandler.AssignInspector)

			// ================= USER =================
			admin.POST("/users", handlers.CreateUser)
			admin.PATCH("/users/:id", handlers.UpdateUser)
			admin.PATCH("/users/:id/activate", handlers.ActivateUser)
			admin.PATCH("/users/:id/deactivate", handlers.DeactivateUser)
			admin.GET("/users", handlers.ListUsers)

			// ================= MASTER DATA =================
			admin.POST("/buildings", buildingHandler.Create)
			admin.GET("/buildings", buildingHandler.List)
			admin.PATCH("/buildings/:id", buildingHandler.Update)
			admin.PATCH("/buildings/:id/deactivate", buildingHandler.Deactivate)

			admin.POST("/areas", areaHandler.Create)
			admin.GET("/areas", areaHandler.List)
			admin.PATCH("/areas/:id", areaHandler.Update)
			admin.PATCH("/areas/:id/deactivate", areaHandler.Deactivate)

			// 🔥🔥🔥 AHU ENDPOINT (INI YANG KITA TAMBAH)
			admin.POST("/ahus", handlers.CreateAHU)
			admin.GET("/ahus", handlers.ListAHUs)
			admin.GET("/ahus/:id", handlers.GetAHUDetail)
			admin.PATCH("/ahus/:id", handlers.UpdateAHU)
			admin.PATCH("/ahus/:id/deactivate", handlers.DeactivateAHU)

			// ============== FORMS ====================
			admin.POST("/form-templates", handlers.CreateFormTemplate)
			admin.GET("/form-templates/:id", handlers.GetFormTemplateDetail)
			admin.GET("/form-templates", handlers.ListFormTemplates)
			admin.PATCH("/form-templates/:id/active", handlers.SetFormTemplateActive)
			admin.POST("/form-templates/:id/version", handlers.CreateNewFormTemplateVersion)

			admin.GET(
				"/form-templates/:id/versions",
				handlers.ListFormTemplateVersions,
			)
			admin.GET(
				"/form-templates/compare/:fromId/:toId",
				handlers.CompareFormTemplate,
			)

			// ================= AUDIT =================
			admin.GET("/audit-trails", handlers.ListAuditTrails)
			admin.GET("/audit-trails/:entity/:id", handlers.ListAuditByEntity)
		}
	}
}
