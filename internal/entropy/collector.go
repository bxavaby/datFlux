package entropy

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"
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

	// metrics for entropy estimation
	sampleVariance map[string]float64
	prevValues     map[string]float64
	entropyScore   float64
}

func NewCollector(samplingRate time.Duration, maxSamples int) *Collector {
	if maxSamples <= 0 {
		maxSamples = 10
	}

	return &Collector{
		samples:        make([]EntropySource, 0, maxSamples),
		maxSamples:     maxSamples,
		samplingRate:   samplingRate,
		lastEntropy:    make([]byte, 32),
		sampleVariance: make(map[string]float64),
		prevValues:     make(map[string]float64),
		entropyScore:   0.0,
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

	if len(c.samples) == 0 {
		// use current time as fallback
		return time.Now().UnixNano()
	}

	var entropyData string
	for _, sample := range c.samples {
		entropyData += fmt.Sprintf(
			"%f|%f|%f|%f|%d|",
			sample.CPU,
			sample.Memory,
			sample.NetworkRx,
			sample.NetworkTx,
			sample.Timestamp,
		)
	}

	// extra entropy sources
	entropyData += fmt.Sprintf("|%d|%d|%.6f|",
		time.Now().UnixNano(),
		c.lastSeedValue,
		c.entropyScore,
	)

	// 2x hash for better distribution
	h1 := sha256.New()
	h1.Write([]byte(entropyData))
	firstHash := h1.Sum(nil)

	h2 := sha256.New()
	h2.Write(firstHash)
	c.lastEntropy = h2.Sum(nil)

	seed := int64(binary.LittleEndian.Uint64(c.lastEntropy[:8]))
	c.lastSeedValue = seed

	return seed
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
