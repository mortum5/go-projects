package dgraph

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"github.com/mortum5/go-projects/dgraph-go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CancelFunc func()

type Url struct {
	model.Url
	Dtype []string `json:"dgraph.type,omitempty"`
}

type DGraph struct {
	*dgo.Dgraph
}

func New() (DGraph, CancelFunc) {
	conn, err := grpc.Dial("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cant connect to graph")
	}
	dc := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	return DGraph{dc}, func() {
		if err := conn.Close(); err != nil {
			log.Printf("dgraph close connect err: %v", err)
		}
	}
}

func (dg DGraph) Setup() {
	err := dg.Alter(context.Background(), &api.Operation{
		Schema: `
			url: string .
			slug: string @index(term) .
			expire: datetime .

			type Url {
				url
				slug
				expire
			}
		`,
	})
	if err != nil {
		log.Printf("init error: %v\n", err)
	}
}

func (dg DGraph) DropAll() {
	err := dg.Alter(context.Background(), &api.Operation{DropOp: api.Operation_ALL})
	if err != nil {
		log.Printf("drop error: %v\n", err)
	}
}

func (dg DGraph) Get(slug string) (model.Url, bool) {
	q := `query get($slug: string) {
		q(func: allofterms(slug, $slug))  {
			uid
			expand(_all_)
		}
	}`

	args := make(map[string]string)
	args["$slug"] = slug

	resp, err := dg.NewReadOnlyTxn().QueryWithVars(context.Background(), q, args)
	if err != nil {
		log.Printf("dgraph: request error: %v\n", err)
	}

	type Root struct {
		Urls []model.Url `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.GetJson(), &r)
	if err != nil {
		log.Printf("dgraph: unmarshall error: %v\n", err)
	}

	if len(r.Urls) > 0 {
		return r.Urls[0], true
	}

	return model.Url{}, false
}

func (dg DGraph) Set(url model.Url) {
	dUrl := Url{
		Url:   url,
		Dtype: []string{"Url"},
	}

	json, _ := json.Marshal(dUrl)

	mu := &api.Mutation{
		SetJson:   json,
		CommitNow: true,
	}

	_, err := dg.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		log.Printf("mutate error: %v\n", err)
	}

}
