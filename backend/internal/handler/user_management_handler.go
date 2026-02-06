package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================= LIST USERS =================
func (h *Handlers) ListUsers(c *gin.Context) {
	users, err := h.UserManagementUC.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// ================= CREATE USER =================
func (h *Handlers) CreateUser(c *gin.Context) {
	var req struct {
		NPP      string `json:"npp"`
		Name     string `json:"name"`
		Jabatan  string `json:"jabatan"` // ⬅️ ganti makna email
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "format request tidak valid"})
		return
	}

	if req.NPP == "" {
		c.JSON(400, gin.H{"error": "NPP wajib diisi"})
		return
	}

	adminID := c.GetString("user_id")

	if err := h.UserManagementUC.CreateUser(
		req.NPP,
		req.Name,
		req.Jabatan,
		req.Password,
		req.Role,
		adminID,
	); err != nil {

		log.Println("CREATE USER ERROR >>>", err.Error()) // ⬅️ TAMBAH DI SINI

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "user berhasil dibuat"})
}

// ================= ACTIVATE =================
func (h *Handlers) ActivateUser(c *gin.Context) {
	userID := c.Param("id")
	adminID := c.GetString("user_id")

	if err := h.UserManagementUC.ActivateUser(userID, adminID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user diaktifkan"})
}

// ================= DEACTIVATE =================
func (h *Handlers) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	adminID := c.GetString("user_id")

	if err := h.UserManagementUC.DeactivateUser(userID, adminID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user dinonaktifkan"})
}

func (h *Handlers) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.UserManagementUC.GetProfile(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id":         user.ID,
		"npp":        user.NPP,
		"name":       user.Name,
		"jabatan":    user.Jabatan, // ⬅️ makna baru
		"role":       user.Role,
		"avatar_url": user.AvatarURL,
	})
}

func (h *Handlers) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Name     *string `json:"name"`
		Jabatan  *string `json:"jabatan"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "request tidak valid"})
		return
	}

	if err := h.UserManagementUC.UpdateUser(
		userID,
		req.Name,
		req.Jabatan,
		req.Password,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user berhasil diupdate"})
}
