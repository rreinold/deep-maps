package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"encoding/json"    
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
    
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
)

func initialize(dgraph *string) (*dgo.Dgraph, *gin.Engine){

	conn, err := grpc.Dial(*dgraph, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	gdb := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	router := gin.Default()
	router.GET("/graph", func(c *gin.Context) {
		fmt.Printf("Called /graph")
		graph := getGraph(gdb)
		fmt.Printf("%+v",graph)
		c.JSON(200, graph)
	})
	router.Run(":3000")
	return gdb, router

}

func main() {
	var (
		dgraph = flag.String("d", "127.0.0.1:9080", "Dgraph Alpha address")
	)
	flag.Parse()

	initialize(dgraph)


}

func getGraph(dgb *dgo.Dgraph) map[string]interface{}{
    
	resp, err := dgb.NewTxn().Query(context.Background(), `{
	  topic(func: has(storyline)) {
	    uid
	    topic_title
	    storyline {
	      name
	      event {
	        name
	        place {
	          name
	        }
	      }
	    }
	  }
	}`)
	
	if err != nil {
		log.Fatal(err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", response)
	return response
}
