package main

import (
	"fmt"
	"time"
)

var k = []string{"single_line_text_field", "number_decimal", "number_integer", "date", "date_time",
	"dimension", "money", "rating", "volume", "weight", "list.single_line_text_field",
	"multi_line_text_field", "page_reference", "product_reference", "collection_reference", "file_reference", "url", "color"}
var pattern = "(%d, %d, '%s', '%s', '%s', '%s', '%s', '',false, 0),"
var ownerResouces = []string{"product", "collection", "page", "shop"}

func main() {
	//owner_resource, namespace, meta_key, meta_type,
	var keys [][]any
	for i := 0; i < len(k); i++ {
		id := time.Now().UnixMicro()
		time.Sleep(1 * time.Microsecond)
		storeId := 205651
		ownerResource := "collection"
		namespace := "custom"
		metaName := fmt.Sprintf("%d_%s_%d", storeId, ownerResource, i)
		metaKey := fmt.Sprintf("%d_%s_%d", storeId, ownerResource, i)
		metaType := k[i]
		keys = append(keys, []any{id, storeId, metaName, ownerResource, namespace, metaKey, metaType})
	}
	p := &MetaFieldsDefinitions{
		Keys:    keys,
		Pattern: pattern,
	}
	p.Output(nil)
	//keys := []string{
	//	"date", "date_time",
	//	"dimension", "money", "rating", "volume", "weight",
	//	"multi_line_text_field",
	//}
	//constants := PhpConstantKeys{
	//	Keys:    keys,
	//	Pattern: "METAFIELD_%s",
	//}
	//constants.Output(func1)
}

func func1(m map[string]string) {
	for k, _ := range m {
		fmt.Printf("self::%s,\n", k)
	}
}
