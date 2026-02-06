package handler

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateProfileRequest struct {
	Name      string  `json:"name"`
	BirthDate *string `json:"birth_date"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ================= GET MY PROFILE =================
func (h *Handlers) GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.UserManagementUC.GetMyProfile(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 🔥 SAFE DEREFERENCE
	var npp, jabatan, avatarURL string
	if user.NPP != nil {
		npp = *user.NPP
	}
	if user.Jabatan != nil {
		jabatan = *user.Jabatan
	}
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	c.JSON(200, gin.H{
		"id":         user.ID,
		"npp":        npp,
		"name":       user.Name,
		"jabatan":    jabatan,
		"role":       user.Role,
		"avatar_url": avatarURL,
	})
}

// ================= UPDATE PROFILE =================
func (h *Handlers) UpdateMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Name      string `json:"name"`
		BirthDate string `json:"birth_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	var birth *time.Time
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid birth date"})
			return
		}
		birth = &t
	}

	if err := h.UserManagementUC.UpdateMyProfile(userID, req.Name, birth); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "profile updated"})
}

// ================= CHANGE PASSWORD =================
func (h *Handlers) ChangeMyPassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if err := h.UserManagementUC.ChangePassword(
		userID,
		req.OldPassword,
		req.NewPassword,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password updated"})
}

// ================= UPLOAD AVATAR =================
func (h *Handlers) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "file avatar wajib"})
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)
	path := "uploads/avatar/" + filename
	avatarURL := "/" + path

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(500, gin.H{"error": "gagal menyimpan file"})
		return
	}

	if err := h.UserManagementUC.UpdateAvatar(userID, avatarURL); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"avatar_url": avatarURL,
	})
}
