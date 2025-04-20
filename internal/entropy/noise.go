package entropy

import (
	"math"
	"math/rand"
	"net"
	"sync"
	"time"
)

type NoiseGenerator struct {
	collector     *Collector
	stopChan      chan struct{}
	wg            sync.WaitGroup
	systemMetrics chan EntropySource
}

func NewNoiseGenerator() *NoiseGenerator {
	ng := &NoiseGenerator{
		collector:     NewCollector(200*time.Millisecond, 10),
		stopChan:      make(chan struct{}),
		systemMetrics: make(chan EntropySource, 10),
	}

	ng.wg.Add(3)
	go ng.generateCPUNoise()
	go ng.generateRAMNoise()
	go ng.generateNetworkNoise()

	ng.wg.Add(1)
	go ng.collectSystemMetrics()

	ng.wg.Add(1)
	go ng.samplingRoutine()

	return ng
}

func (ng *NoiseGenerator) Collector() *Collector {
	return ng.collector
}

func (ng *NoiseGenerator) Stop() {
	close(ng.stopChan)
	ng.wg.Wait()
}

func (ng *NoiseGenerator) generateCPUNoise() {
	defer ng.wg.Done()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ng.stopChan:
			return
		default:
			x := r.Float64() * 10000

			for i := 0; i < 10000; i++ {
				x = math.Sin(x) * math.Cos(x) * math.Tan(x)
				x = math.Sqrt(math.Abs(x)) + math.Log(math.Abs(x+1))
			}
		}
	}
}

func (ng *NoiseGenerator) generateRAMNoise() {
	defer ng.wg.Done()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memoryChunks := make([][]byte, 0, 20)

	for {
		select {
		case <-ng.stopChan:
			for i := range memoryChunks {
				memoryChunks[i] = nil
			}
			memoryChunks = nil
			return

		default:
			// clear old allocations periodically to prevent OOM
			if len(memoryChunks) > 15 {
				memoryChunks = memoryChunks[8:]
			}

			// 5MB-50MB
			size := r.Intn(45*1024*1024) + (5 * 1024 * 1024)
			data := make([]byte, size)

			for i := range data {
				if i%1024 == 0 { // operate on every 1024th byte
					data[i] = byte(r.Intn(256))
				}
			}

			memoryChunks = append(memoryChunks, data)

			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ng *NoiseGenerator) generateNetworkNoise() {
	defer ng.wg.Done()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	targetPorts := []string{
		"127.0.0.1:7",    // echo port
		"127.0.0.1:13",   // daytime
		"127.0.0.1:37",   // time
		"127.0.0.1:8080", // HTTP port
		"127.0.0.1:8888", // alternative HTTP port
	}

	// local listener setup to respond to generated traffic
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:8888")
		if err == nil {
			defer listener.Close()

			for {
				conn, err := listener.Accept()
				if err != nil {
					break
				}

				// echo back received data
				go func(c net.Conn) {
					defer c.Close()
					buffer := make([]byte, 1024)
					c.Read(buffer)
					c.Write(buffer)
				}(conn)
			}
		}
	}()

	for {
		select {
		case <-ng.stopChan:
			return

		default:
			// transport protocol: TCP = reliable, UDP = speed
			proto := "tcp"
			if r.Intn(2) == 0 {
				proto = "udp"
			}

			target := targetPorts[r.Intn(len(targetPorts))]

			// packets: 1KB-100KB
			size := r.Intn(99*1024) + 1024
			data := make([]byte, size)
			r.Read(data)

			conn, err := net.DialTimeout(proto, target, 100*time.Millisecond)
			if err == nil {
				conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
				conn.Write(data)

				// try to read response (good for entropy)
				responseBuffer := make([]byte, 1024)
				conn.Read(responseBuffer)
				conn.Close()
			}

			// < sleep = > network activity
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (ng *NoiseGenerator) collectSystemMetrics() {
	defer ng.wg.Done()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ng.stopChan:
			return
		case <-ticker.C:
			stats := getSystemStats()

			select {
			case ng.systemMetrics <- stats:
			default:
				// channel buffer full, move on
			}
		}
	}
}

func (ng *NoiseGenerator) samplingRoutine() {
	defer ng.wg.Done()

	for {
		select {
		case <-ng.stopChan:
			return
		case stats := <-ng.systemMetrics:
			ng.collector.AddSample(stats)
		}
	}
}

func getSystemStats() EntropySource {
	// actual system stats in monitor/system.go
	// for now just return placeholder values
	return EntropySource{
		CPU:       rand.Float64() * 100,
		Memory:    rand.Float64() * 100,
		NetworkRx: rand.Float64() * 1024 * 1024,
		NetworkTx: rand.Float64() * 1024 * 1024,
		Timestamp: time.Now().UnixNano(),
	}
}
