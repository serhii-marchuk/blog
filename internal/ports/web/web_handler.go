package web

import (
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/constructors"
	"github.com/serhii-marchuk/blog/internal/models"
	"gorm.io/gorm"
	"net/http"
)

type WebHandler struct {
	RedisClient *redis.Client
	Db          *gorm.DB
	Renderer    *constructors.Renderer
}

func NewWebHandler(
	rc *constructors.RedisClient,
	db *constructors.Db,
	r *constructors.Renderer,
) *WebHandler {
	return &WebHandler{
		RedisClient: rc.RC,
		Db:          db.DB,
		Renderer:    r,
	}
}

func (h *WebHandler) Setup(e *echo.Echo) {
	e.GET("/", h.BasePage)
	e.GET("/:page", h.BasePage)
	e.POST("/contact/form", h.HandleContactForm)
	e.GET("/assets/:assetsFilePath", h.Assets)
}

func (h *WebHandler) BasePage(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		page = "home"
	}

	return h.render(c, page)
}

func (h *WebHandler) Assets(c echo.Context) error {
	return c.File(fmt.Sprintf("web/assets/%s", c.Param("assetsFilePath")))
}

func (h *WebHandler) HandleContactForm(c echo.Context) error {
	h.Db.Save(&models.ContactRequest{
		ID:          uuid.Generate(),
		FirstName:   c.FormValue("first_name"),
		LastName:    c.FormValue("last_name"),
		Email:       c.FormValue("email"),
		Description: c.FormValue("description"),
		Status:      "new",
	})

	return h.render(c, "contact")
}

func (h *WebHandler) render(c echo.Context, pageName string) error {
	html, err := h.Renderer.Render(h.RedisClient, pageName)
	if err != nil {
		return err
	}

	return c.HTMLBlob(http.StatusOK, html)
}
