package handlers

import (
	"net/http"
	"strconv"

	"github.com/ZAUakaAlexey/backend_go/internal/database"
	"github.com/ZAUakaAlexey/backend_go/internal/models"
	"github.com/ZAUakaAlexey/backend_go/internal/responses"
	"github.com/ZAUakaAlexey/backend_go/internal/validators"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(context *gin.Context) {
	userId, exists := context.Get("user_id")
	if !exists {
		errors := map[string][]string{
			"auth": {"User not authenticated"},
		}
		responses.ErrorResponse(context, http.StatusUnauthorized, "Authentication required", errors)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		errors := map[string][]string{
			"user": {"User not found"},
		}
		responses.ErrorResponse(context, http.StatusNotFound, "User not found", errors)
		return
	}

	responses.SuccessResponse(context, http.StatusOK, user, "User retrieved successfully")
}

// GetUsers возвращает список пользователей с пагинацией
func GetUsers(c *gin.Context) {
	// Получаем параметры пагинации из query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Подсчет общего количества записей
	var total int64
	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to count users", nil)
		return
	}

	// Если записей нет
	if total == 0 {
		meta := responses.PaginationMeta{
			CurrentPage: page,
			From:        0,
			LastPage:    0,
			PerPage:     perPage,
			To:          0,
			Total:       0,
		}
		responses.PaginatedSuccessResponse(
			c,
			http.StatusOK,
			[]models.User{},
			meta,
			"No users found",
		)
		return
	}

	// Вычисляем offset
	offset := (page - 1) * perPage

	// Получаем пользователей с пагинацией
	var users []models.User
	if err := database.DB.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", nil)
		return
	}

	// Вычисляем метаданные пагинации
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	from := offset + 1
	to := offset + len(users)

	meta := responses.PaginationMeta{
		CurrentPage: page,
		From:        from,
		LastPage:    lastPage,
		PerPage:     perPage,
		To:          to,
		Total:       int(total),
	}

	responses.PaginatedSuccessResponse(
		c,
		http.StatusOK,
		users,
		meta,
		"Users retrieved successfully",
	)
}

// SearchUsers поиск пользователей с пагинацией
func SearchUsers(c *gin.Context) {
	// Получаем параметры
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Базовый запрос
	query := database.DB.Model(&models.User{})

	// Применяем фильтр поиска
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Подсчет общего количества записей
	var total int64
	if err := query.Count(&total).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to count users", nil)
		return
	}

	// Если записей нет
	if total == 0 {
		meta := responses.PaginationMeta{
			CurrentPage: page,
			From:        0,
			LastPage:    0,
			PerPage:     perPage,
			To:          0,
			Total:       0,
		}
		responses.PaginatedSuccessResponse(
			c,
			http.StatusOK,
			[]models.User{},
			meta,
			"No users found",
		)
		return
	}

	// Вычисляем offset
	offset := (page - 1) * perPage

	// Получаем пользователей с пагинацией
	var users []models.User
	if err := query.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", nil)
		return
	}

	// Вычисляем метаданные пагинации
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	from := offset + 1
	to := offset + len(users)

	meta := responses.PaginationMeta{
		CurrentPage: page,
		From:        from,
		LastPage:    lastPage,
		PerPage:     perPage,
		To:          to,
		Total:       int(total),
	}

	responses.PaginatedSuccessResponse(
		c,
		http.StatusOK,
		users,
		meta,
		"Search completed successfully",
	)
}

// GetUser получить одного пользователя по ID
func GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		errors := map[string][]string{
			"user_id": {"User not found"},
		}
		responses.ErrorResponse(c, http.StatusNotFound, "User not found", errors)
		return
	}

	responses.SuccessResponse(c, http.StatusOK, user, "User retrieved successfully")
}

// UpdateUser обновление пользователя
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	// Проверяем, что пользователь обновляет свой профиль
	currentUserID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Простая проверка прав (пользователь может обновлять только свой профиль)
	if strconv.Itoa(int(currentUserID.(uint))) != userID {
		errors := map[string][]string{
			"permission": {"You can only update your own profile"},
		}
		responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", errors)
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := validators.FormatValidationErrors(err)
		responses.ValidationErrorResponse(c, "Validation failed", validationErrors)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		errors := map[string][]string{
			"user_id": {"User not found"},
		}
		responses.ErrorResponse(c, http.StatusNotFound, "User not found", errors)
		return
	}

	// Обновляем только переданные поля
	if input.Name != "" {
		user.Name = input.Name
	}

	if err := database.DB.Save(&user).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", nil)
		return
	}

	responses.SuccessResponse(c, http.StatusOK, user, "User updated successfully")
}

// DeleteUser удаление пользователя (soft delete)
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// Проверяем, что пользователь удаляет свой профиль
	currentUserID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Простая проверка прав
	if strconv.Itoa(int(currentUserID.(uint))) != userID {
		errors := map[string][]string{
			"permission": {"You can only delete your own profile"},
		}
		responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", errors)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		errors := map[string][]string{
			"user_id": {"User not found"},
		}
		responses.ErrorResponse(c, http.StatusNotFound, "User not found", errors)
		return
	}

	// Soft delete (GORM automatically uses DeletedAt)
	if err := database.DB.Delete(&user).Error; err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", nil)
		return
	}

	responses.SuccessResponse(c, http.StatusOK, nil, "User deleted successfully")
}
