package web

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/constructors"
	"github.com/serhii-marchuk/blog/internal/models"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"time"
)

type WebHandler struct {
	NavCfg      *configs.NavConfig
	Pages       []configs.NavItem
	RedisClient *redis.Client
	Db          *gorm.DB
}

func NewWebHandler(
	cfg *configs.WebCfg,
	rc *constructors.RedisClient,
	db *constructors.Db,
) *WebHandler {
	return &WebHandler{
		NavCfg:      cfg.NavCfg,
		RedisClient: rc.RC,
		Db:          db.DB,
		Pages:       cfg.NavCfg.NavBar,
	}
}

func (h *WebHandler) Setup(e *echo.Echo) {
	e.GET("/", h.Page)
	e.GET("/:page", h.Page)
	e.POST("/contact/form", h.HandleContactForm)
	e.GET("/assets/:assetsFilePath", h.Assets)
}

func (h *WebHandler) Page(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		page = "home"
	}

	return h.Render(c, http.StatusOK, page)
}

func (h *WebHandler) Render(
	c echo.Context,
	code int,
	pageName string,
) error {
	var html []byte
	var err error

	html, _ = h.RedisClient.Get(context.Background(), pageName).Bytes()
	if html == nil {
		html, err = h.getPage(pageName)
		if err != nil {
			return err
		}
	}

	return c.HTMLBlob(code, html)
}

func (h *WebHandler) getPage(pageName string) ([]byte, error) {
	for k, item := range h.Pages {
		if item.Name == pageName {
			h.Pages[k].Active = true
		} else {
			h.Pages[k].Active = false
		}
	}
	tmpl := template.Must(template.ParseFiles(
		h.NavCfg.GetContentFileByPageName(pageName),
		h.NavCfg.BaseTemplatePath),
	)

	buf := new(bytes.Buffer)
	tmpErr := tmpl.ExecuteTemplate(buf, "base", map[string]interface{}{"pages": &h.Pages})
	if tmpErr != nil {
		return nil, tmpErr
	}

	rdsErr := h.RedisClient.Set(context.Background(), pageName, buf.Bytes(), 24*time.Hour).Err()
	if rdsErr != nil {
		return nil, rdsErr
	}

	return buf.Bytes(), nil
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

	return h.Render(c, http.StatusOK, "contact")
}
