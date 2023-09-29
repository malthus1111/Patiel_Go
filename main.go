package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Icic maxGoroutines définit le nombre maximal de goroutines à exécuter en même temps.
const maxGoroutines = 100

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxGoroutines)
	foundPort := make(chan int)

	// Ceci fait une Boucle à travers les ports et lance une goroutine pour chaque port.
	for port := 1025; port <= 8192; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(port int) {
			defer wg.Done()
			defer func() { <-sem }()

			url := fmt.Sprintf("http://10.49.122.144:%d/ping", port)
			resp, err := http.Get(url)
			if err != nil {
				//Ceci affiche l'erreur et sort de la goroutine.
				fmt.Printf("Erreur avec le port %d: %v\n", port, err)
				return
			}
			defer resp.Body.Close()

			// Si le statut de la réponse est 200, signale le port trouvé.
			if resp.StatusCode == http.StatusOK {
				foundPort <- port
			}
		}(port)
	}

	go func() {
		wg.Wait()
		close(foundPort)
	}()

	// Ici on attends un port valide.
	port := <-foundPort

	// Ici on fait une requête POST sur /signup.
	signupURL := fmt.Sprintf("http://10.49.122.144:%d/signup", port)
	user := map[string]string{"user": "Malthus"}
	userJson, _ := json.Marshal(user)
	resp, err := http.Post(signupURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête signup:", err)
		return
	}
	fmt.Println("Réponse du /signup:", resp.Status)

	// Ici on fait une requête POST sur /check.
	checkURL := fmt.Sprintf("http://10.49.122.144:%d/check", port)
	resp, err = http.Post(checkURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête check:", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /check:", string(body))
}
