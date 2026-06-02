package elasticsearch

import (
	"bytes"
	"io"
)

func CreateProductIndex(
	tenantCode string,
	countryCode string,
) error {

	indexName :=
		GetProductIndex(
			tenantCode,
			countryCode,
		)

	exists, err :=
		Client.Indices.Exists(
			[]string{
				indexName,
			},
		)

	if err != nil {
		return err
	}

	defer exists.Body.Close()

	if exists.StatusCode == 200 {
		return nil
	}

	response, err :=
		Client.Indices.Create(
			indexName,
			Client.Indices.Create.WithBody(
				bytes.NewReader(
					[]byte(
						ProductMapping(),
					),
				),
			),
		)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	_, err =
		io.ReadAll(response.Body)

	return err
}
