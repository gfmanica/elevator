package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NumeroDeSapos    = 5
	DistanciaCorrida = 50
)

// sapo é uma goroutine que simula um sapo na corrida.
func sapo(nome string, linhaDeChegada chan<- string, wg *sync.WaitGroup, declareWinnerOnce *sync.Once) {
	// garante que o WaitGroup será notificado quando a goroutine terminar.
	defer wg.Done()

	distanciaPecorrida := 0

	for distanciaPecorrida < DistanciaCorrida {
		// Simula um pulo com uma pausa aleatória.
		salto := rand.Intn(5) + 1

		distanciaPecorrida += salto

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		
		fmt.Printf("Sapo %s pulou para a posição %d\n", nome, distanciaPecorrida)
	}

	fmt.Printf("--- Sapo %s cruzou a linha de chegada! ---\n", nome)

	// usa sync.Once para garantir que o vencedor seja declarado apenas uma vez.
	// apenas a primeira goroutine a chamar Do() executará a função.
	declareWinnerOnce.Do(func() {
		linhaDeChegada <- nome
	})
}

func main() {
	var wg sync.WaitGroup
	var declareWinnerOnce sync.Once // variável para garantir uma única execução.

	// canal para comunicar o vencedor. O buffer de 1 garante que
	// o primeiro sapo a chegar não fique bloqueado.
	linhaDeChegada := make(chan string, 1)

	fmt.Println("Começou a corrida!")

	// inicia as goroutines dos sapos.
	for i := 1; i <= NumeroDeSapos; i++ {
		wg.Add(1)
		nomeSapo := fmt.Sprintf("%d", i)
		go sapo(nomeSapo, linhaDeChegada, &wg, &declareWinnerOnce)
	}

	// espera o primeiro sapo cruzar a linha de chegada.
	vencedor := <-linhaDeChegada
	fmt.Printf("\n🎉 O Sapo %s venceu a corrida! 🎉\n\n", vencedor)

	// espera todos os outros sapos terminarem a corrida.
	wg.Wait()

	fmt.Println("Todos os sapos terminaram a corrida.")
}
