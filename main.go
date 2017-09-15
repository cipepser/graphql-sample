package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Repos struct {
	Data struct {
		Viewer struct {
			Repositories struct {
				Edges []struct {
					Node struct {
						Name string `json:"name"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"repositories"`
		} `json:"viewer"`
	} `json:"data"`
}

func main() {
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			hello
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}

	// c, err := sdk.NewClient("https://api.github.com/graphql", "user", "pass", nil)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	//
	// body := `{"query": "query { viewer { name }}"}`
	//
	// req, err := c.NewRequest(ctx, "POST", "", strings.NewReader(body))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// req.Header.Set("Authorization", "bearer "+os.Getenv("GITHUB_ACCESSTOKEN"))
	//
	// resp, err := c.HTTPClient.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// b, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(b))
	//
	// r := Repos{}
	// err = sdk.DecodeBody(resp, &r)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(r)

}
