package server

import (
	"auth/iternal/services/notifications"
	"auth/iternal/token"
	"auth/pkg/database/models"
	"auth/pkg/logger"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (r *Router) createToken(c fiber.Ctx) error {
	guid := c.Query("id")

	if _, err := uuid.Parse(guid); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "неверный id")
	}

	ip := c.IP()

	accessToken, session, err := token.Create(r.cfg, guid, ip)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось создать токен")
	}

	refresh, hash, err := token.CreateRefresh()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось создать refresh-токен")
	}

	modelToken := models.RefreshToken{
		UserGUID:  guid,
		TokenHash: hash,
		Session:   session,
		IPAddress: ip,
		Used:      false,
	}

	if err := modelToken.Create(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось сохранить refresh-токен")
	}

	response := tokenPair{
		AccessToken:  accessToken,
		RefreshToken: refresh,
	}

	return c.JSON(response)
}

func (r *Router) refreshToken(c fiber.Ctx) error {
	var request tokenPair

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "неверное тело запроса")
	}

	accessToken, err := token.Parse(r.cfg, request.AccessToken)
	if err != nil || !accessToken.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "неверный токен доступа")
	}

	claims := accessToken.Claims.(jwt.MapClaims)
	session := claims["session"].(string)
	guid := claims["sub"].(string)
	ip := claims["ip"].(string)
	currentIP := c.IP()

	// Получаем токен из БД по JTI
	refreshToken, err := models.FindRefreshBySession(session)
	if err != nil || refreshToken.Used {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh-токен не найден или уже был использован ранее")
	}

	if !token.ValidateRefresh(request.RefreshToken, refreshToken.TokenHash) {
		return fiber.NewError(fiber.StatusUnauthorized, "неверный токен/refresh-токен")
	}

	if ip != currentIP {
		err := notifications.SendEmail(fmt.Sprintf("В Ваш аккаунт вошли с другого IP: %s", currentIP), "ВНИМАНИЕ! НОВОЕ УСТРОЙСТВО!", "test@email.ru")
		if err != nil {
			logger.Warn("Не удалось отправить уведомление на почту: %v", err)
		}
	}

	if err := refreshToken.Marked(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось обновить refresh-токен")
	}

	newAccessToken, newSession, err := token.Create(r.cfg, guid, currentIP)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось создать токен")
	}

	newRefreshToken, newHash, err := token.CreateRefresh()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось создать refresh-token")
	}

	newToken := models.RefreshToken{
		UserGUID:  guid,
		TokenHash: newHash,
		Session:   newSession,
		IPAddress: currentIP,
		Used:      false,
	}

	if err := newToken.Create(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "не удалось сохранить refresh-token")
	}

	response := tokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return c.JSON(response)
}
