package elasticsearch

func ProductMapping() string {

	return `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {

      "id": {
        "type": "long"
      },

      "name": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
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
        "type": "float"
      },

      "sale_price": {
        "type": "float"
      },

      "status": {
        "type": "keyword"
      },

      "created_at": {
        "type": "date"
      }
    }
  }
}
`
}
