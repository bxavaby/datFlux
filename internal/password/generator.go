package password

import (
	"math/rand"

	"datflux/internal/entropy"
)

type Generator struct {
	collector  *entropy.Collector
	minLength  int
	maxLength  int
	useSymbols bool
	useNumbers bool
	useUpper   bool
	useLower   bool
}

func NewGenerator(collector *entropy.Collector) *Generator {
	return &Generator{
		collector:  collector,
		minLength:  16,
		maxLength:  32,
		useSymbols: true,
		useNumbers: true,
		useUpper:   true,
		useLower:   true,
	}
}

// high-entropy password
func (g *Generator) Generate() string {
	seed := g.collector.GenerateSeed()

	source := rand.NewSource(seed)
	secureRand := rand.New(source)

	var lowercase = "abcdefghijklmnopqrstuvwxyz"
	var uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var numbers = "0123456789"
	var symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?/"

	var allChars string
	var requiredChars []byte

	if g.useLower {
		allChars += lowercase
		requiredChars = append(requiredChars, lowercase[secureRand.Intn(len(lowercase))])
	}

	if g.useUpper {
		allChars += uppercase
		requiredChars = append(requiredChars, uppercase[secureRand.Intn(len(uppercase))])
	}

	if g.useNumbers {
		allChars += numbers
		requiredChars = append(requiredChars, numbers[secureRand.Intn(len(numbers))])
	}

	if g.useSymbols {
		allChars += symbols
		requiredChars = append(requiredChars, symbols[secureRand.Intn(len(symbols))])
	}

	passLength := g.minLength + secureRand.Intn(g.maxLength-g.minLength+1)
	passLength = max(passLength, len(requiredChars))

	password := make([]byte, passLength)

	for i, char := range requiredChars {
		password[i] = char
	}

	for i := len(requiredChars); i < passLength; i++ {
		password[i] = allChars[secureRand.Intn(len(allChars))]
	}

	for i := range password {
		j := secureRand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}

func (g *Generator) GenerateRandomChar() byte {
	var allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-="
	return allChars[rand.Intn(len(allChars))]
}
