package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	out := new(bytes.Buffer)
	er := new(bytes.Buffer)
	cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
	cmd.Stdout = out
	cmd.Stderr = er

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var person struct {
		Name string
		Age  int
	}
	// if err := json.NewDecoder(stdout).Decode(&person); err != nil {
	// 	log.Fatal(err)
	// }
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(out).Decode(&person); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is %d years old\n", person.Name, person.Age)
}
