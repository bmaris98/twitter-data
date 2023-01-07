package main

import (
	"fmt"
	"log"
	"os/exec"
)

func ExecHadoopJob(query string, id string) {
	hdfsConn := HdfsConnection{}
	hdfsConn.Init()
	outputFolder := fmt.Sprintf("/output/%s", query)
	hdfsConn.EnsureFolderStructure(outputFolder)
	err := hdfsConn.Client.Close()
	if err != nil {
		log.Fatal("Failed to close tmp hdfs conn", err)
	}
	wrapperCommand := fmt.Sprintf("echo '$HADOOP_HOME/bin/hadoop jar $JAR_FILEPATH $CLASS_TO_RUN %s %s' | nc -u -w 3 namenode 4444", query, id)

	_, err = exec.Command("/bin/sh", "-c", wrapperCommand).Output()

	if err != nil {
		log.Fatal("Failed to send netcat udp command package", err)
	}
}
