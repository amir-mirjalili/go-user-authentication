package handlers

import (
	"net/http"

	"github.com/amir-mirjalili/go-user-authentication/internal/params"
	"github.com/amir-mirjalili/go-user-authentication/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SendOTP godoc
// @Summary Send OTP to phone number
// @Description Send a 6-digit OTP code to the specified phone number
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{phone_number=string} true "Phone number"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Router /auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req params.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.authService.SendOTP(req.PhoneNumber)
	if err != nil {
		if err.Error() == "rate limit exceeded: maximum 3 OTP requests per 10 minutes" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP godoc
// @Summary Verify OTP and authenticate user
// @Description Verify the OTP code and return authentication token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{phone_number=string,code=string} true "Phone and OTP"
// @Success 200 {object} object{token=string,user=object}
// @Failure 401 {object} object{error=string}
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req params.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.authService.VerifyOTPAndLogin(req.PhoneNumber, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
