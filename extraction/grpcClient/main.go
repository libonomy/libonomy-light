package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"

	"github.com/libonomy/libonomy-light/extraction/proto"

	"google.golang.org/grpc"

	"github.com/libonomy/libonomy-light/extraction/constants"
)

func main() {

	// fmt.Print("Gin Context", gin.Context{})
	fmt.Print("Background Context", context.Background())

	conn, err := grpc.Dial("localhost"+constants.Server.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("CDonnot Dial", err.Error())
	}

	client := proto.NewCheckingServicesClient(conn)
	fmt.Println(client, conn)

	router := mux.NewRouter()
	router.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		customMessage := &proto.Request{Body: "Binding"}
		responseRPC, err := client.GetResponse(context.Background(), customMessage)
		if err != nil {

			w.Write([]byte("Something Went Wrong"))
			return
		}
		w.Write([]byte(responseRPC.Body))

	}).Methods("GET")
	router.HandleFunc("/uploadFile", func(w http.ResponseWriter, r *http.Request) {
		name := "testing.txt"
		fmt.Println("Printing Name Param From ", name)
		// w.Write([]byte("File Name You Enter is \t" + name))
		fileContent, err := ioutil.ReadFile("Files/testing.txt")
		if err != nil {
			w.Write([]byte("File Name You Enter is \t" + err.Error()))
			return
		}
		fmt.Println(fileContent)
		fileParams := &proto.Files{
			FileName:    name,
			FileContent: fileContent,
		}
		res, err := client.UploadFile(context.Background(), fileParams)
		if err != nil {
			w.Write([]byte("File Name You Enter is \t" + err.Error()))
			return
		}
		w.Write([]byte("Response From Server\t" + res.GetMsg()))
	}).Methods("GET")
	http.ListenAndServe(":4700", router)
}
