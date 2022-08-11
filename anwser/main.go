package main

import "fmt"

func main() {
	m := map[string]float32{
		"code_update":                0.889,
		"collections_create":         1,
		"collections_delete":         0.378,
		"collections_update":         125,
		"country_blockers":           0.111,
		"discount-flashsale_refresh": 1.84,
		"filters_update":             0.0667,
		"global-config_create":       0.267,
		"ip_blockers":                0.0889,
		"menu_delete":                0.289,
		"menu_update":                0.178,
		"page_create":                0.398,
		"page_delete":                0.244,
		"page_update":                0.333,
		"product_create":             48.1,
		"product_delete":             40,
		"product_update":             258,
		"redirect_update":            0.133,
		"script-tag_create":          0.133,
		"script-tag_delete":          0.0917,
		"script-tag_update":          1.69,
		"seo_update":                 0.0222,
		"store2_update":              0.267,
		"store-subscribe":            0.161,
		"store_update":               0.245,
		"app-install":                0.0444,
		"app-uninstall":              0.0343333,
		"themes_publish":             0.333,
		"themes_update":              5.45,
		"website":                    113,
	}
	r := Percent(m)
	for k, v := range r {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func Percent(m map[string]float32) map[string]string {
	var (
		sum float32 = 0
		res         = make(map[string]string)
	)

	for _, v := range m {
		sum += v
	}
	for k, v := range m {
		res[k] = fmt.Sprintf("%.2f%%", v/sum*100)
	}
	res["sum"] = fmt.Sprintf("%.2f", sum)
	return res
}
