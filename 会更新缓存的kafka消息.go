package main

var RedirectUpdateTopic = "redirect_update"

var ProductsUpdateTopic = "products_update"

var StoreSubscribeTopic = map[string]bool{
	"store_subscribe": true,
	"store2_update":   true,
}

var WebhookTopic = map[string]bool{
	"products_create":    true,
	"products_update":    true,
	"products_unpublish": true,
	"products_publish":   true,
	"products_delete":    true,

	"collections_update": true,
	"collections_delete": true,
	"collections_create": true,

	"discount-flashsale_refresh": true,
	"page_update":                true,
	"page_create":                true,
	"page_delete":                true,

	"code_update": true,

	"seo_update": true,

	"script-tag_create": true,
	"script-tag_update": true,
	"script-tag_delete": true,

	"global-config_create": true,
	"global-config_delete": true,

	"store_update": true,

	"themes_publish": true,
	"themes_update":  true,
	"filters_update": true,

	"third_part_search_enable":  true,
	"third_part_search_disable": true,
	"theme_app_install":         true,
	"theme_app_uninstall":       true,

	"menu_update": true,
	"menu_delete": true,

	"redirect_update": true,
}

var MenuUpdateTopics = map[string]bool{
	"products_update":    true,
	"collections_update": true,
	"page_update":        true,
}

var CanalTopic = map[string]bool{
	"website": true,
}

var CfDataUpdateTopics = map[string]bool{
	"ip_blockers":      true,
	"country_blockers": true,
}
