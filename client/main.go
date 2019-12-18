package main

import (
	"example.com/grpc-test/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	g :=gin.Default()

	g.GET("/add/:a/:b", func(ctx *gin.Context) {
		a, err := strconv.ParseInt(ctx.Param("a"), 10 ,64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter a"})
			return
		}
		b, err := strconv.ParseInt(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter b"})
			return
		}
		req := &proto.Request{A:int64(a), B:int64(b)}

		if response, err := client.Add(ctx, req);err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Result)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	if err := g.Run(":8080"); err !=nil {
		log.Fatal("Failed to run server: %v",err)
	}
}
