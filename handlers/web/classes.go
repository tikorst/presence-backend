package web

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/tikorst/presence-backend/config"
    "github.com/tikorst/presence-backend/models"
)

func Classes() gin.HandlerFunc {
    return func(c *gin.Context) {
        var classes []models.Kelas
        claims, _ := c.Get("claims")
        jwtClaims := claims.(jwt.MapClaims)
        username := jwtClaims["sub"].(string)

        user := models.User{}
        var latestSemester models.Semester
        if err := config.DB.
            Last(&latestSemester).Error; err != nil {
            c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
            return
        }
        if err := config.DB.
            Joins("JOIN dosen_pengampu ON dosen_pengampu.id_kelas = kelas.id_kelas").
            Preload("MataKuliah").
            Preload("DosenPengampu").
            Where("dosen_pengampu.nip = ? AND id_semester = ?", username, latestSemester.IDSemester).
            Find(&classes).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
            return
        }

        if err:= config.DB.
            Where("username = ?", username).
            First(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"classes": classes, "user" : user})
    }
}
