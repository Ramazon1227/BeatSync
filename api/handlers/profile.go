package handlers

import (
	"context"
	"fmt"
	// "fmt"
	// "hash"

	httpapi "github.com/Ramazon1227/BeatSync/api/http"
	"github.com/Ramazon1227/BeatSync/models"
	"github.com/Ramazon1227/BeatSync/pkg/jwt.go"
	"github.com/Ramazon1227/BeatSync/pkg/utils"
	"github.com/Ramazon1227/BeatSync/storage"
	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @ID register-user
// @Router /v1/auth/register [POST]
// @Summary Register User
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegisterModel true "user registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) RegisterUser(c *gin.Context) {
    var user models.UserRegisterModel

    err := c.ShouldBindJSON(&user)
    if err != nil {
        h.handleResponse(c, httpapi.BadRequest, err)
        return
    }

	if user.Password != user.ConfirmPassword {
		h.handleResponse(c, httpapi.BadRequest, "passwords do not match")
		return
	}

    resp, err := h.storage.User().Create(context.Background(), &user)
    if err != nil {
        h.handleResponse(c, httpapi.InternalServerError, err)
        return
    }

    status := httpapi.Created
    status.Description = "user registered successfully"
    h.handleResponse(c, status, resp)
}




// Login godoc
// @ID login
// @Router /v1/auth/login [POST]
// @Summary Login User
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} httpapi.Response
// @Failure 401 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err.Error())
		return
	}

	user, err := h.storage.User().GetByEmail(context.Background(), req.Email)
	if err != nil {
		if err == storage.ErrorNotFound {
			h.handleResponse(c, httpapi.Unauthorized, "invalid credentials")
			return
		}
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}
    fmt.Println("User found:", user)
	// Check if the password is correct
	if !utils.CheckPassword(user.Password, req.Password) {
		h.handleResponse(c, httpapi.Unauthorized, "invalid credentials")
		return
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, &models.LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: jwt.GetTokenExpiryTime(),
	})
}

// Logout godoc
// @ID logout
// @Router /v1/auth/logout [POST]
// @Summary Logout User
// @Description Invalidate user's JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param logout body models.LogoutRequest true "logout request"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 401 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) Logout(c *gin.Context) {
	var req models.LogoutRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err.Error())
		return
	}

	// Add token to blacklist or invalidate it
	err = jwt.InvalidateToken(req.Token)
	if err != nil {
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, "successfully logged out")
}

// GetProfile godoc
// @ID get-profile
// @Router /v1/profile [GET]
// @Summary Get User Profile
// @Description Get authenticated user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 401 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// func (h *Handler) GetProfile(c *gin.Context) {
//     userId, exists := c.Get("user_id")
//     if !exists {
//         h.handleResponse(c, httpapi.Unauthorized, "user not authenticated")
//         return
//     }

//     user, err := h.storage.User().GetById(context.Background(), &models.PrimaryKey{Id: userId.(string)})
//     if err != nil {
//         h.handleResponse(c, httpapi.InternalServerError, err)
//         return
//     }

//     // Clear sensitive information
//     user.Password = ""

//     c.JSON(httpapi.OK.Code, user)
// }

// UpdateProfile godoc
// @ID update-profile
// @Router /v1/profile [PUT]
// @Summary Update User Profile
// @Description Update authenticated user's profile
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param profile body models.UpdateProfileRequest true "profile data"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 401 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// func (h *Handler) UpdateProfile(c *gin.Context) {
// 	userId, exists := c.Get("user_id")
// 	if !exists {
// 		h.handleResponse(c, httpapi.Unauthorized, "unauthorized")
// 		return
// 	}

// 	var req models.UpdateProfileRequest
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		h.handleResponse(c, httpapi.BadRequest, err.Error())
// 		return
// 	}

// 	err = h.storage.User().UpdateUserProfile(context.Background(), userId.(string), &req)
// 	if err != nil {
// 		if err == storage.ErrorNotFound {
// 			h.handleResponse(c, httpapi.NoContent, err)
// 			return
// 		}
// 		h.handleResponse(c, httpapi.InternalServerError, err)
// 		return
// 	}

// 	h.handleResponse(c, httpapi.OK, "profile updated successfully")
// }

// UpdatePassword godoc
// @ID update-password
// @Router /v1/profile/password [PUT]
// @Summary Update User Password
// @Description Update authenticated user's password
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body models.UpdatePasswordRequest true "password data"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 401 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// func (h *Handler) UpdatePassword(c *gin.Context) {
// 	userId, exists := c.Get("user_id")
// 	if !exists {
// 		h.handleResponse(c, httpapi.Unauthorized, "unauthorized")
// 		return
// 	}

// 	var req models.UpdatePasswordRequest
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		h.handleResponse(c, httpapi.BadRequest, err.Error())
// 		return
// 	}

// 	err = h.storage.User().UpdatePassword(context.Background(), userId.(string), req.CurrentPassword, req.NewPassword)
// 	if err != nil {
// 		if err == storage.ErrorNotFound {
// 			h.handleResponse(c, httpapi.NoContent, err)
// 			return
// 		}
// 		h.handleResponse(c, httpapi.InternalServerError, err)
// 		return
// 	}

// 	h.handleResponse(c, httpapi.OK, "password updated successfully")
// }


// DeleteUser godoc
// @ID delete-user
// @Router /v1/profile/{user_id} [DELETE]
// @Summary Delete User
// @Description Delete a user by ID
// @Tags profile
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 204 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
// func (h *Handler) DeleteUser(c *gin.Context) {
// 	id := c.Param("user_id")
// 	if id == "" {
// 		h.handleResponse(c, httpapi.BadRequest, "user id required")
// 		return
// 	}

// 	err := h.storage.User().Delete(context.Background(), &models.PrimaryKey{Id: id})
// 	if err != nil {
// 		if err == storage.ErrorNotFound {
// 			h.handleResponse(c, httpapi.NoContent, err)
// 			return
// 		}
// 		h.handleResponse(c, httpapi.InternalServerError, err)
// 		return
// 	}

// 	h.handleResponse(c, httpapi.OK, "successfully deleted")
// }