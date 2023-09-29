package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const fixedPort = 5174 // Définir le port fixe ici

func main() {
	url := fmt.Sprintf("http://10.49.122.144:%d/ping", fixedPort)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur avec le port %d: %v\n", fixedPort, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Le port %d ne répond pas avec le statut 200 OK\n", fixedPort)
		return
	}

	// Ici on fait une requête POST sur /signup.
	signupURL := fmt.Sprintf("http://10.49.122.144:%d/signup", fixedPort)
	user := map[string]string{"user": "Malthus"}
	userJson, _ := json.Marshal(user)
	resp, err = http.Post(signupURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête signup:", err)
		return
	}
	fmt.Println("Réponse du /signup:", resp.Status)

	// Ici on fait une requête POST sur /check.
	checkURL := fmt.Sprintf("http://10.49.122.144:%d/check", fixedPort)
	resp, err = http.Post(checkURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête check:", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /check:", string(body))

	getUserSecretURL := fmt.Sprintf("http://10.49.122.144:%d/getUserSecret", fixedPort)
	resp, err = http.Post(getUserSecretURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête getUserSecret:", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /getUserSecret:", string(body))

	// Extrait le "Secret" de la réponse.
	var secretMap map[string]string
	err = json.Unmarshal(body, &secretMap)
	if err != nil {
		fmt.Println("Erreur lors de la déserialisation du Secret:", err)
		return
	}
	secret := secretMap["Secret"]

	// Ici on fait une requête POST sur /getUserLevel.
	getUserLevelURL := fmt.Sprintf("http://10.49.122.144:%d/getUserLevel", fixedPort)
	levelUser := map[string]string{"User": "Malthus", "Secret": secret}
	levelUserJson, _ := json.Marshal(levelUser)
	resp, err = http.Post(getUserLevelURL, "application/json", bytes.NewBuffer(levelUserJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête getUserLevel:", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /getUserLevel:", string(body))

	// Effectuer une requête POST sur /getUserPoints.
	getUserPointsURL := fmt.Sprintf("http://10.49.122.144:%d/getUserPoints", fixedPort)
	userPointsData := map[string]string{"User": "Malthus", "Secret": secret}
	userPointsJson, _ := json.Marshal(userPointsData)
	resp, err = http.Post(getUserPointsURL, "application/json", bytes.NewBuffer(userPointsJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête getUserPoints:", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /getUserPoints:", string(body))
}

// le code ci-dessous est le même code que le code plus haut, mais avec la recheche aléatoire du port.
// Il fonctionne bien, mais mets un peu de temps à rechercher le bon port.
// Contrairement à celui du haut qui fonctionne directement avec le bon port

/*
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

	getUserSecretURL := fmt.Sprintf("http://10.49.122.144:%d/getUserSecret", port)
	resp, err = http.Post(getUserSecretURL, "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête getUserSecret:", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /getUserSecret:", string(body))

	// Extrait le "Secret" de la réponse.
	var secretMap map[string]string
	err = json.Unmarshal(body, &secretMap)
	if err != nil {
		fmt.Println("Erreur lors de la déserialisation du Secret:", err)
		return
	}
	secret := secretMap["Secret"]

	// Ici on fait une requête POST sur /getUserLevel.
	getUserLevelURL := fmt.Sprintf("http://10.49.122.144:%d/getUserLevel", port)
	levelUser := map[string]string{"User": "Malthus", "Secret": secret}
	levelUserJson, _ := json.Marshal(levelUser)
	resp, err = http.Post(getUserLevelURL, "application/json", bytes.NewBuffer(levelUserJson))
	if err != nil {
		fmt.Println("Erreur lors de la requête getUserLevel:", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Réponse du /getUserLevel:", string(body))
}

*/
