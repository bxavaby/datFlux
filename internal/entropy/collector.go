package entropy

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/dchest/blake2s"
	"github.com/seehuhn/fortuna"
)

type EntropySource struct {
	CPU       float64
	Memory    float64
	NetworkRx float64
	NetworkTx float64
	Timestamp int64
}

type Collector struct {
	mu            sync.Mutex
	samples       []EntropySource
	maxSamples    int
	samplingRate  time.Duration
	lastEntropy   []byte
	lastSeedValue int64

	// Fortuna RNG
	rng         *fortuna.Accumulator
	entropySink chan<- []byte

	// metrics for entropy estimation
	sampleVariance map[string]float64
	prevValues     map[string]float64
	entropyScore   float64
}

func NewCollector(samplingRate time.Duration, maxSamples int) *Collector {
	if maxSamples <= 0 {
		maxSamples = 10
	}

	homeDir, _ := os.UserHomeDir()
	seedFile := homeDir + "/.datflux_seed"

	rng, err := fortuna.NewRNG(seedFile)
	if err != nil {
		// if the file-backed RNG cannot be created, use an in-memory one
		rng, _ = fortuna.NewRNG("")
	}

	sink := rng.NewEntropyDataSink()

	return &Collector{
		samples:        make([]EntropySource, 0, maxSamples),
		maxSamples:     maxSamples,
		samplingRate:   samplingRate,
		lastEntropy:    make([]byte, 32),
		sampleVariance: make(map[string]float64),
		prevValues:     make(map[string]float64),
		entropyScore:   0.0,
		rng:            rng,
		entropySink:    sink,
	}
}

func (c *Collector) Close() {
	if c.entropySink != nil {
		close(c.entropySink)
		c.entropySink = nil
	}

	if c.rng != nil {
		c.rng.Close()
	}
}

func (c *Collector) AddSample(source EntropySource) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.samples = append(c.samples, source)
	if len(c.samples) > c.maxSamples {
		c.samples = c.samples[1:]
	}

	c.updateEntropyEstimate(source)

	// add sample data to Fortuna
	sampleBytes := []byte(fmt.Sprintf(
		"%f|%f|%f|%f|%d|%d",
		source.CPU,
		source.Memory,
		source.NetworkRx,
		source.NetworkTx,
		source.Timestamp,
		time.Now().UnixNano(),
	))

	// send to entropy sink
	if c.entropySink != nil {
		select {
		case c.entropySink <- sampleBytes:
		// sent successfully
		default:
			// channel full, skip sample
		}
	}
}

func (c *Collector) updateEntropyEstimate(source EntropySource) {
	metrics := map[string]float64{
		"cpu":       source.CPU,
		"memory":    source.Memory,
		"networkRx": source.NetworkRx,
		"networkTx": source.NetworkTx,
	}

	totalEntropy := 0.0
	metricCount := len(metrics)

	for name, value := range metrics {
		prev, exists := c.prevValues[name]
		if !exists {
			prev = value
		}

		delta := math.Abs(value - prev)

		variance, exists := c.sampleVariance[name]
		if !exists {
			variance = delta
		} else {
			// exponential moving average with 0.7 weight for new values
			variance = variance*0.3 + delta*0.7
		}
		c.sampleVariance[name] = variance

		c.prevValues[name] = value

		// normalized entropy contribution
		// > variance = > entropy
		metricEntropy := math.Min(1.0, variance/100.0)
		totalEntropy += metricEntropy
	}

	// average entropy (scale to 0-1)
	c.entropyScore = math.Min(1.0, totalEntropy/float64(metricCount))
}

func (c *Collector) GenerateSeed() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rng == nil {
		// use current time as fallback
		return time.Now().UnixNano()
	}

	// get 8 bytes from Fortuna
	randomBytes := c.rng.RandomData(8)
	seed := int64(binary.LittleEndian.Uint64(randomBytes))
	c.lastSeedValue = seed

	// get 32 bytes for entropy quality
	c.lastEntropy = c.rng.RandomData(32)

	return seed
}

// returns the full 32 bytes (256 bits) of entropy
func (c *Collector) GetRawEntropy() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rng == nil {
		// fallback to Blake2s hash of system time
		h := blake2s.New256()
		h.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
		return h.Sum(nil)
	}

	// return 32 bytes (256 bits) of entropy from Fortuna
	return c.rng.RandomData(32)
}

func (c *Collector) GetEntropyQuality() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. sample buffer fullness
	bufferFill := float64(len(c.samples)) / float64(c.maxSamples)

	// 2. dynamic score based on metric variance
	dynamicScore := c.entropyScore

	// 3. time-based factor (entropy improves over time as more events occur)
	// gives users a sense of improvement with time
	timeFactor := math.Min(1.0, float64(len(c.samples))*0.1)

	// weighted combination
	qualityScore := (bufferFill * 0.3) + (dynamicScore * 0.6) + (timeFactor * 0.1)

	return math.Min(1.0, qualityScore)
}
