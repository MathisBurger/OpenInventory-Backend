package middleware

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
)

// Base struct
type TwoFactorPair struct {
	Username string
	Token    string
}

// Global available array of 2fa pairs
var TwoFactorPairs []TwoFactorPair

var TwoFactorCommunicationChannel chan TwoFactorPair

// entrypoint from main function
func TwoFactorService() {
	TwoFactorCommunicationChannel = make(chan TwoFactorPair)
	go twoFaService(TwoFactorCommunicationChannel)
}

// infinite running service to handle new 2fa sessions
func twoFaService(c chan TwoFactorPair) {

	sessions := actions.GetAll2FaSessions()
	for _, s := range sessions {
		TwoFactorPairs = append(TwoFactorPairs, TwoFactorPair{s.Owner, s.Secret})
	}

	for {
		val := <-c
		TwoFactorPairs = append(TwoFactorPairs, val)
	}
}
