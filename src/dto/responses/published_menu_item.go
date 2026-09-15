package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedMenuItem struct {
	ID            uint                `json:"id"`
	Position      uint                `json:"position"`
	Page          string              `json:"page"`
	Name          string              `json:"name"`
	URLName       string              `json:"urlName"`
	Hashtag       *string             `json:"hashtag"`
	Icon          *string             `json:"icon"`
	IsHome        bool                `json:"isHome"`
	IsCustom      bool                `json:"isCustom"`
	Url           *string             `json:"url"`
	UrlEnabled    bool                `json:"urlEnabled"`
	NewTabEnabled bool                `json:"newTabEnabled"`
	Routes        map[string][]string `json:"routes"`
	Items         []PublishedMenuItem `json:"items"`

	routeMenuItemIDs []uint
}

// SetMenuItem sets the MenuItem response from the models.MenuItem model.
func (pmi *PublishedMenuItem) SetMenuItem(menuItem *models.MenuItem, position uint, pages *[]models.Page, routeMenuItemIDs []uint) {
	var page *models.Page
	if len(menuItem.Pages) > 0 {
		page = &menuItem.Pages[0]
	}

	pmi.ID = menuItem.ID
	pmi.Position = position
	pmi.Page = menuItem.Name
	pmi.Routes = getMenuItemRoutes(routeMenuItemIDs, pages)
	pmi.routeMenuItemIDs = routeMenuItemIDs

	if page != nil {
		pmi.Name = page.Name
		pmi.URLName = utils.URLEncode(page.Name)
		pmi.UrlEnabled = page.UrlEnabled
		pmi.NewTabEnabled = page.NewTabEnabled
		pmi.Hashtag = utils.PtrFromNullString(page.Hashtag)
		pmi.Url = utils.PtrFromNullString(page.Url)
	}

	pmi.Icon = utils.PtrFromNullString(menuItem.Icon)
	pmi.IsHome = menuItem.IsHome
	pmi.IsCustom = menuItem.IsCustom

	pmi.Items = make([]PublishedMenuItem, 0)
}

// getMenuItemRoutes returns complete route segments per locale for the ordered menu item path.
func getMenuItemRoutes(menuItemIDs []uint, pages *[]models.Page) map[string][]string {
	routes := make(map[string][]string)
	if pages == nil || len(menuItemIDs) == 0 {
		return routes
	}

	for i, menuItemID := range menuItemIDs {
		pageNamesByLocale := make(map[string]string)
		for j := range *pages {
			page := &(*pages)[j]
			if page.MenuItemID == menuItemID {
				pageNamesByLocale[page.Locale] = utils.URLEncode(page.Name)
			}
		}

		if i == 0 {
			for locale, name := range pageNamesByLocale {
				routes[locale] = []string{name}
			}
			continue
		}

		for locale := range routes {
			if name, ok := pageNamesByLocale[locale]; ok {
				routes[locale] = append(routes[locale], name)
				continue
			}

			delete(routes, locale)
		}
	}

	return routes
}
