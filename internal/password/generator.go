package password

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"datflux/internal/entropy"

	"github.com/nbutton23/zxcvbn-go"
)

type Generator struct {
	collector       *entropy.Collector
	minLength       int
	maxLength       int
	useSymbols      bool
	useNumbers      bool
	useUpper        bool
	useLower        bool
	paranoiaMode    bool
	paranoiaSamples int
}

type PasswordStrength struct {
	Score            int     // 0-4 (0=very weak, 4=very strong)
	EntropyBits      float64 // estimated entropy bits
	CrackTimeDesc    string  // human-readable crack time estimate
	CrackTimeSeconds float64 // estimated crack time in s
	Feedback         string  // feedback to improve the password
}

type AttackModelType int

const (
	OnlineRateLimited  AttackModelType = iota // default zxcvbn (100 guesses/sec)
	OfflineGPUCracking                        // serious offline attack (1 billion guesses/sec)
	QuantumComputing                          // theoretical future attack (10^15+ guesses/sec)
)

type AttackModel struct {
	Name          string
	Description   string
	GuessesPerSec float64
}

func GetAttackModels() []AttackModel {
	return []AttackModel{
		{
			Name:          "Online Rate-Limited",
			Description:   "Standard online attack with rate limiting (100 guesses/sec)",
			GuessesPerSec: 100, // 10ms per guess
		},
		{
			Name:          "Offline GPU Cracking",
			Description:   "Serious password file breach (1 billion guesses/sec)",
			GuessesPerSec: 1e9, // 1 billion
		},
		{
			Name:          "Quantum Computing",
			Description:   "State-level adversary with advanced tech (10^15 guesses/sec)",
			GuessesPerSec: 1e15, // 1 quadrillion
		},
	}
}

func NewGenerator(collector *entropy.Collector) *Generator {
	return &Generator{
		collector:       collector,
		minLength:       16,
		maxLength:       32,
		useSymbols:      true,
		useNumbers:      true,
		useUpper:        true,
		useLower:        true,
		paranoiaMode:    false,
		paranoiaSamples: 25,
	}
}

func (g *Generator) SetParanoiaMode(enabled bool, samples int) {
	g.paranoiaMode = enabled
	g.paranoiaSamples = max(1, samples)
}

func (g *Generator) GetParanoiaMode() (bool, int) {
	return g.paranoiaMode, g.paranoiaSamples
}

func (g *Generator) Generate() string {
	if g.paranoiaMode {
		return g.generateParanoid()
	}

	// standard
	seed := g.collector.GenerateSeed()
	return g.generateWithSeed(seed)
}

// out of 25 passwords, returns the one with highest entropy
func (g *Generator) generateParanoid() string {
	// pre-allocate candidates and corresponding entropy values
	candidates := make([]string, g.paranoiaSamples)
	entropies := make([]float64, g.paranoiaSamples)

	// goroutines to get candidates in parallel
	var wg sync.WaitGroup
	wg.Add(g.paranoiaSamples)

	for i := range make([]struct{}, g.paranoiaSamples) {
		go func(index int) {
			defer wg.Done()

			// 512-bit entropy for each candidate password
			entropyBytes := g.collector.GetRawEntropy512()
			// first 8 bytes for seed
			seed := int64(binary.LittleEndian.Uint64(entropyBytes[:8]))
			candidates[index] = g.generateWithSeed(seed)

			// analyze strength
			strength := g.AnalyzeStrength(candidates[index])
			entropies[index] = strength.EntropyBits
		}(i)
	}

	wg.Wait()

	// keep the one with highest entropy
	bestIndex := 0
	for i, entropy := range entropies {
		if entropy > entropies[bestIndex] {
			bestIndex = i
		}
	}

	return candidates[bestIndex]
}

// uses the given seed to generate a password
func (g *Generator) generateWithSeed(seed int64) string {
	// nosec G404 -- uses a securely seeded PRNG
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
		// nosec G404 -- uses a securely seeded PRNG
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

	// longer passwords in paranoia mode
	passLength := g.minLength + secureRand.Intn(g.maxLength-g.minLength+1)
	if g.paranoiaMode {
		// 48-80 chars in paranoia mode
		passLength = 48 + secureRand.Intn(33)
	}
	passLength = max(passLength, len(requiredChars))

	password := make([]byte, passLength)

	for i, char := range requiredChars {
		password[i] = char
	}

	for i := len(requiredChars); i < passLength; i++ {
		password[i] = allChars[secureRand.Intn(len(allChars))]
	}

	// Fisher-Yates shuffle
	// multiple passes in paranoia mode
	numPasses := 1
	if g.paranoiaMode {
		numPasses = 3 // more thorough shuffling in paranoia mode
	}

	for range make([]struct{}, numPasses) {
		for i := range password {
			j := secureRand.Intn(i + 1)
			password[i], password[j] = password[j], password[i]
		}
	}

	return string(password)
}

func (g *Generator) GenerateRandomChar() byte {
	var allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-="
	// nosec G404 -- only used for visual animation
	return allChars[rand.Intn(len(allChars))]
}

// uses zxcvbn to evaluate password strength
func (g *Generator) AnalyzeStrength(password string) PasswordStrength {
	result := zxcvbn.PasswordStrength(password, nil)

	// feedback based on score
	feedback := ""
	if result.Score < 3 {
		if len(password) < 12 {
			feedback = "Consider a longer password"
		} else {
			feedback = "Try adding more varied characters"
		}
	}

	return PasswordStrength{
		Score:            result.Score,
		EntropyBits:      float64(result.Entropy),
		CrackTimeDesc:    result.CrackTimeDisplay,
		CrackTimeSeconds: result.CrackTime,
		Feedback:         feedback,
	}
}

// crack time adjusted for a specific attack model
func (g *Generator) GetAdjustedCrackTime(baseSeconds float64, modelType AttackModelType) float64 {
	models := GetAttackModels()

	// zxcvbn's default is 10ms/guess
	defaultGuessesPerSec := 100.0
	modelGuessesPerSec := models[modelType].GuessesPerSec

	adjustmentFactor := modelGuessesPerSec / defaultGuessesPerSec
	return baseSeconds / adjustmentFactor
}

func (g *Generator) GetCrackTimeForModel(password string, modelType AttackModelType) string {
	result := zxcvbn.PasswordStrength(password, nil)

	// for quantum computing, use entropy directly
	if modelType == QuantumComputing {
		entropyBits := float64(result.Entropy)
		quantumEntropyBits := entropyBits / 2
		quantumCrackTime := math.Pow(2, quantumEntropyBits) / GetAttackModels()[QuantumComputing].GuessesPerSec
		return GetCrackTimeDescription(quantumCrackTime)
	}

	adjustedTime := g.GetAdjustedCrackTime(result.CrackTime, modelType)

	return GetCrackTimeDescription(adjustedTime)
}

// descriptions based on crack time
func GetCrackTimeDescription(seconds float64) string {
	// time constants
	minute := float64(60)
	hour := minute * 60
	day := hour * 24
	month := day * 30
	year := day * 365
	decade := year * 10
	century := year * 100

	universeAge := 13.8 * 1e9 * year

	switch {
	case seconds < 0.001:
		return "instant"
	case seconds < 1:
		return "< 1 second"
	case seconds < minute:
		return fmt.Sprintf("%d seconds", int(seconds))
	case seconds < hour:
		return fmt.Sprintf("%d minutes", int(seconds/minute))
	case seconds < day:
		return fmt.Sprintf("%d hours", int(seconds/hour))
	case seconds < 7*day:
		return fmt.Sprintf("%d days", int(seconds/day))
	case seconds < month:
		return fmt.Sprintf("%d weeks", int(seconds/(7*day)))
	case seconds < year:
		return fmt.Sprintf("%d months", int(seconds/month))
	case seconds < decade:
		return fmt.Sprintf("%d years", int(seconds/year))
	case seconds < century:
		return fmt.Sprintf("%d decades", int(seconds/decade))
	case seconds < 10*century:
		return fmt.Sprintf("%d centuries", int(seconds/century))
	case seconds < universeAge:
		billionYears := seconds / year / 1e9

		if billionYears < 0.1 {
			millennia := int(seconds / (10 * century))
			return fmt.Sprintf("%d millennia", millennia)
		}

		return fmt.Sprintf("%.1f billion years", billionYears)
	default:
		cosmicScale := seconds / universeAge
		if cosmicScale < 1000 {
			// at least 1.1x, to avoid strange values
			if cosmicScale < 1.1 {
				return "the age of the universe"
			}
			return fmt.Sprintf("%.1f× the age of the universe", cosmicScale)
		}
		return "until the heat death of the universe"
	}
}
