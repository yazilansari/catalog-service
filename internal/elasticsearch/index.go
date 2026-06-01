package elasticsearch

import (
	"bytes"
	"context"
)

func CreateProductIndex() error {

	mapping := `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },

      "name": {
        "type": "text"
      },

      "slug": {
        "type": "keyword"
      },

      "category": {
        "type": "keyword"
      },

      "subcategory": {
        "type": "keyword"
      },

      "brand": {
        "type": "keyword"
      },

      "price": {
        "type": "double"
      },

      "discount_price": {
        "type": "double"
      },

      "status": {
        "type": "keyword"
      },

      "created_at": {
        "type": "date"
      },

      "sales_count": {
        "type": "long"
      }
    }
  }
}`

	_, err :=
		Client.Indices.Create(
			"products",

			Client.Indices.Create.WithContext(
				context.Background(),
			),

			Client.Indices.Create.WithBody(
				bytes.NewReader(
					[]byte(mapping),
				),
			),
		)

	return err
}
