package main

import (
	"fmt"
	"log"

	"com.lisitede.backstage.gitea/internal/gitea"
)

func main() {
	adapter, err := gitea.NewAdapter()
	if err != nil {
		log.Fatal(err)
	}

	giteaVersion, err := adapter.GetGiteaVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Gitea Version: %s\n", giteaVersion)
	fmt.Printf("CLI Version: %s\n", cliVersion)
}
