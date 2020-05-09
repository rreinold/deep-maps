package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"encoding/json"    
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
    
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"str/config"
)

func main() {
	var (
		dgraph = flag.String("d", "127.0.0.1:9080", "Dgraph Alpha address")
	)
	flag.Parse()

	initialize(dgraph)


}

func initialize(dgraph *string) (*dgo.Dgraph, *gin.Engine){
	configuration := config.GetConfig()
	fmt.Println(configuration)
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
	router.GET("/locations", func(c *gin.Context) {
		latString := c.DefaultQuery("lat", "30.268071")
		lat, errLat := strconv.ParseFloat(latString, 64)
		lngString := c.DefaultQuery("lng", "-97.742802")
		lng, errLng := strconv.ParseFloat(lngString, 64)
		if errLng != nil || errLat != nil {
			fmt.Printf("Invalid coordinates provided")
			c.JSON(400, "Invalid coordinates provided")
		}
		fmt.Printf("Called /location, using lat %v lng %v",lat, lng)
		locations := searchLocations(gdb, lat, lng)
		fmt.Printf("%+v",locations)
		c.JSON(200, locations)
	})
	router.Run(":3000")
	return gdb, router

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

func searchLocations(dgb *dgo.Dgraph, lat float64 , lng float64) map[string]interface{}{
    query := fmt.Sprintf(`{
	  locations(func: near(location, [%v,%v], 4000) ) {
	    name
	    location
	  }
	}`,lng, lat)
	fmt.Println(query)
	resp, err := dgb.NewTxn().Query(context.Background(), query)
	
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
