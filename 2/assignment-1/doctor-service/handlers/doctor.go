package handlers

import "github.com/gin-gonic/gin"

func DoctorGETByIDHandler(c *gin.Context) {
	id := c.Param("id")
}

func DoctorGETListHandler(c *gin.Context) {
}

type DoctorPOSTBind struct {
	FullName       string `json:"full_name" binding:"required"`
	Specialization string `json:"specialization" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

func DoctorPOSTHandler(c *gin.Context) {
	var doctor_bind DoctorPOSTBind

	if err := c.BindJSON(&doctor_bind); err != nil {
		return
	}
}
