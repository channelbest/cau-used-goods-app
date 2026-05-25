package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/pkg/response"
)

type Handler struct {
	service *Service
}

type updateProfileRequest struct {
	Nickname  *string `json:"nickname"`
	AvatarURL *string `json:"avatarUrl"`
	Phone     *string `json:"phone"`
}

type submitStudentVerificationRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	RealName  string `json:"realName" binding:"required"`
	College   string `json:"college" binding:"required"`
}

type reviewStudentVerificationRequest struct {
	AuthStatus  string `json:"authStatus" binding:"required"`
	Description string `json:"description"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if req.AvatarURL != nil {
		avatarURL, err := saveRemoteAvatar(*req.AvatarURL)
		if err != nil {
			response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
			return
		}
		req.AvatarURL = &avatarURL
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, UpdateProfileInput{
		Nickname:  req.Nickname,
		AvatarURL: req.AvatarURL,
		Phone:     req.Phone,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "avatar file is required")
		return
	}
	if file.Size > 2*1024*1024 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "avatar file must be <= 2MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "avatar file must be jpg, jpeg, png or webp")
		return
	}

	name, err := randomFileName(ext)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "generate avatar filename failed")
		return
	}

	dir := filepath.Join("uploads", "avatar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "create avatar directory failed")
		return
	}

	dst := filepath.Join(dir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "save avatar file failed")
		return
	}

	avatarURL := "/" + filepath.ToSlash(dst)
	user, err := h.service.UpdateProfile(c.Request.Context(), userID, UpdateProfileInput{
		AvatarURL: &avatarURL,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, gin.H{
		"avatarUrl": avatarURL,
		"user":      user,
	})
}

func (h *Handler) SubmitStudentVerification(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req submitStudentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	verification, err := h.service.SubmitStudentVerification(c.Request.Context(), userID, SubmitStudentVerificationInput{
		StudentID: req.StudentID,
		RealName:  req.RealName,
		College:   req.College,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, verification)
}

func randomFileName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf) + ext, nil
}

func saveRemoteAvatar(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("avatarUrl is empty")
	}
	if strings.HasPrefix(rawURL, "/uploads/avatar/") {
		return rawURL, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", fmt.Errorf("avatarUrl must be http(s) URL or /uploads/avatar path")
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("download avatar failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download avatar failed with status %d", resp.StatusCode)
	}
	if resp.ContentLength > 2*1024*1024 {
		return "", fmt.Errorf("avatar file must be <= 2MB")
	}

	ext := avatarExt(resp.Header.Get("Content-Type"), filepath.Ext(parsed.Path))
	if ext == "" {
		return "", fmt.Errorf("avatar file must be jpg, jpeg, png or webp")
	}

	name, err := randomFileName(ext)
	if err != nil {
		return "", fmt.Errorf("generate avatar filename failed")
	}

	dir := filepath.Join("uploads", "avatar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create avatar directory failed")
	}

	dst := filepath.Join(dir, name)
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create avatar file failed")
	}
	defer out.Close()

	limited := io.LimitReader(resp.Body, 2*1024*1024+1)
	written, err := io.Copy(out, limited)
	if err != nil {
		return "", fmt.Errorf("save avatar file failed")
	}
	if written > 2*1024*1024 {
		_ = os.Remove(dst)
		return "", fmt.Errorf("avatar file must be <= 2MB")
	}

	return "/" + filepath.ToSlash(dst), nil
}

func avatarExt(contentType string, pathExt string) string {
	switch strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}

	switch strings.ToLower(pathExt) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return strings.ToLower(pathExt)
	default:
		return ""
	}
}

func (h *Handler) StudentVerification(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	verification, err := h.service.StudentVerification(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, verification)
}

func (h *Handler) ListStudentVerifications(c *gin.Context) {
	status := c.DefaultQuery("authStatus", "PENDING")
	items, err := h.service.ListStudentVerifications(c.Request.Context(), status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"items": items})
}

func (h *Handler) ReviewStudentVerification(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || userID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid user id")
		return
	}

	var req reviewStudentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	verification, err := h.service.ReviewStudentVerification(c.Request.Context(), adminID, ReviewStudentVerificationInput{
		UserID:      userID,
		AuthStatus:  req.AuthStatus,
		Description: req.Description,
	})
	if err != nil {
		response.Error(c, http.StatusConflict, response.CodeConflict, err.Error())
		return
	}
	response.Success(c, verification)
}
