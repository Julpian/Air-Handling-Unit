package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	NPP      string `json:"npp"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req loginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "format request tidak valid"})
		return
	}

	if req.NPP == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "NPP dan password wajib diisi"})
		return
	}

	token, role, err := h.AuthUC.Login(req.NPP, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "NPP atau password salah"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"role":  role,
	})
}

func (h *Handlers) Logout(c *gin.Context) {
    // Ambil data dari context (hasil extract token di middleware)
    valID, _ := c.Get("user_id")
    valName, _ := c.Get("user_name")

    // Konversi ke string dengan aman
    userID, _ := valID.(string)
    userName, _ := valName.(string)

    // Sekarang argumennya sudah pas: 2 string
    h.AuthUC.Logout(userID, userName)

    c.JSON(http.StatusOK, gin.H{
        "message": "Logout berhasil",
    })
}