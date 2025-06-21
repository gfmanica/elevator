package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NumFloors     = 10
	NumPassengers = 5 // Menos passageiros para uma demonstração mais clara
)

// PassengerRequest representa uma solicitação de viagem simples.
type PassengerRequest struct {
	ID               int
	DestinationFloor int
}

// elevator é uma goroutine que simula um elevador simples.
// Ele começa no andar 0 e atende às solicitações na ordem em que chegam.
func elevator(requests <-chan PassengerRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Elevador iniciado no andar 0.")
	currentFloor := 0

	// O 'for range' processa cada solicitação que chega no canal.
	// O loop termina quando o canal 'requests' é fechado pela função main.
	for req := range requests {
		fmt.Printf("Elevador: Recebeu solicitação do Passageiro %d para o andar %d.\n", req.ID, req.DestinationFloor)

		// Move para o andar de destino
		for currentFloor != req.DestinationFloor {
			time.Sleep(500 * time.Millisecond) // Simula o tempo de viagem
			if currentFloor < req.DestinationFloor {
				currentFloor++
			} else {
				currentFloor--
			}
			fmt.Printf("Elevador: movendo... andar atual %d\n", currentFloor)
		}

		fmt.Printf("Elevador: Chegou ao andar %d. Passageiro %d desembarcou.\n", currentFloor, req.ID)
	}

	fmt.Println("Elevador: Todas as solicitações atendidas. Desligando.")
}

// passenger é uma goroutine que simula uma pessoa solicitando uma viagem.
func passenger(id int, requests chan<- PassengerRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	destination := rand.Intn(NumFloors)
	fmt.Printf("Passageiro %d: Pressionou o botão para o andar %d.\n", id, destination)

	// Envia a solicitação para o canal do elevador.
	requests <- PassengerRequest{ID: id, DestinationFloor: destination}
}

func main() {
	var elevatorWg sync.WaitGroup
	var passengerWg sync.WaitGroup

	// Canal para os passageiros enviarem suas solicitações de destino.
	requests := make(chan PassengerRequest, NumPassengers)

	// Inicia a goroutine do elevador.
	elevatorWg.Add(1)
	go elevator(requests, &elevatorWg)

	// Inicia as goroutines dos passageiros.
	for i := 1; i <= NumPassengers; i++ {
		passengerWg.Add(1)
		go passenger(i, requests, &passengerWg)
	}

	// Espera todos os passageiros terminarem de fazer suas solicitações.
	passengerWg.Wait()

	// Fecha o canal de solicitações para sinalizar ao elevador que não há mais passageiros.
	close(requests)

	// Espera o elevador terminar de processar todas as solicitações.
	elevatorWg.Wait()

	fmt.Println("Simulação concluída.")
}
