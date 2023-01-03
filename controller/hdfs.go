package main

import (
	"fmt"
	"log"

	"github.com/colinmarc/hdfs"
)

type HdfsConnection struct {
	Client *hdfs.Client
}

func (conn *HdfsConnection) Init() {
	client, err := hdfs.New("namenode:9000")
	conn.Client = client
	if err != nil {
		log.Fatal("Error establishing hdfs connection", err)
	}
}

func (conn *HdfsConnection) CreateFile(query string, filename string, content string) {
	dir := fmt.Sprintf("/input/%s", query)
	err := conn.Client.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal("Failed to mkdir in hdfs", err)
	}
	file := fmt.Sprintf("%s/%s", dir, filename)
	writer, err := conn.Client.Create(file)

	if err != nil {
		log.Fatal("Failed to create hdfs file", err)
	}

	_, err = writer.Write([]byte(content))
	if err != nil {
		log.Fatal("Failed to write to hdfs file", err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal("Failed to close hdfs file", err)
	}
}

func (conn *HdfsConnection) EnsureFolderStructure(structure string) {
	err := conn.Client.MkdirAll(structure, 0777)
	if err != nil {
		log.Fatal("Failed to ensure mkdir in hdfs", err)
	}
}
