package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"github.com/libonomy/libonomy-light/extraction/constants"
	"github.com/libonomy/libonomy-light/extraction/proto"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	listener, err := net.Listen(constants.Server.Protocol, constants.Server.Port)
	if err != nil {
		log.Fatalln("Error in listening. \n", err.Error())
	}

	srv := grpc.NewServer()
	// grpc.WithDialer
	proto.RegisterCheckingServicesServer(srv, &server{})
	// proto.RegisterCheckingServicesServer(srv, &server{})
	// reflection.Register(srv)

	err = srv.Serve(listener)
	if err != nil {
		log.Fatalln("Error is serving grpc", err.Error())
	}

}

func (s *server) GetResponse(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	requestBody := request.GetBody()

	responseMessage := "Request Body is \t" + requestBody + "\nSending This Response \tWorking"
	return &proto.Response{Body: responseMessage}, nil
}

func (s *server) UploadFile(ctx context.Context, params *proto.Files) (*proto.FileResponse, error) {
	fileName := params.GetFileName()
	fileContent := params.GetFileContent()

	err := ioutil.WriteFile("grpc/uploaded_"+fileName, fileContent, 0777)
	if err != nil {
		return &proto.FileResponse{}, err
	}
	return &proto.FileResponse{Msg: "Successfully Uploaded File on Server"}, nil
}

func (s *server) DownloadFile(ctx context.Context, name *proto.DownloadFileName) (*proto.Files, error) {
	return &proto.Files{}, nil
}
