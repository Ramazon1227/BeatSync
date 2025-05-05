package hrv

import (
	"fmt"
	"math"
	"time"

	"github.com/Ramazon1227/BeatSync/models"
	"github.com/scientificgo/fft"
)

// SDNN: Standard deviation of NN intervals
func CalculateSDNN(rr []float64) float64 {
	if len(rr) == 0 {
		return 0
	}
	mean := mean(rr)
	var sumSquares float64
	for _, r := range rr {
		sumSquares += math.Pow(r-mean, 2)
	}
	return math.Sqrt(sumSquares / float64(len(rr)))
}

// RMSSD: Root mean square of successive differences
func CalculateRMSSD(rr []float64) float64 {
	if len(rr) < 2 {
		return 0
	}
	var sum float64
	for i := 1; i < len(rr); i++ {
		diff := rr[i] - rr[i-1]
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(rr)-1))
}

// NN50: Number of pairs with diff > 50ms
func CalculateNN50(rr []float64) int {
	count := 0
	for i := 1; i < len(rr); i++ {
		if math.Abs(rr[i]-rr[i-1]) > 50 {
			count++
		}
	}
	return count
}

// PNN50: NN50 count as percentage of total pairs
func CalculatePNN50(rr []float64) float64 {
	if len(rr) < 2 {
		return 0
	}
	return float64(CalculateNN50(rr)) / float64(len(rr)-1) * 100
}

// SD1/SD2: From Poincare plot geometry
func CalculateSD1SD2(rr []float64) (float64, float64) {
	if len(rr) < 2 {
		return 0, 0
	}
	var diffs []float64
	for i := 1; i < len(rr); i++ {
		diffs = append(diffs, rr[i]-rr[i-1])
	}
	sd1 := math.Sqrt(0.5) * stddev(diffs)
	sd2 := math.Sqrt(2*math.Pow(stddev(rr), 2) - math.Pow(sd1, 2))
	return sd1, sd2
}


// --- Helper Functions ---

func mean(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func stddev(data []float64) float64 {
	m := mean(data)
	var sum float64
	for _, v := range data {
		sum += math.Pow(v-m, 2)
	}
	return math.Sqrt(sum / float64(len(data)))
}


// Extract RR intervals from PPG data (in milliseconds)
func ExtractRR(ppg []models.PPGData, minPeakDistanceMs int, minPeakHeight float64) []float64 {
	if len(ppg) < 3 {
		fmt.Println("Error: Insufficient PPG data for peak detection")
		return nil
	}

	var rr []float64
	var lastPeakTime *time.Time

	// Validate PPG data
	for i := 1; i < len(ppg); i++ {
		if ppg[i].TimeStamp.Before(*ppg[i-1].TimeStamp) {
			fmt.Println("Error: PPG timestamps are not sorted")
			return nil
		}
	}

	for i := 1; i < len(ppg)-1; i++ {
		// Peak detection: value greater than neighbors and threshold
		if ppg[i].Value > ppg[i-1].Value &&
			ppg[i].Value > ppg[i+1].Value &&
			ppg[i].Value >= minPeakHeight {
			if lastPeakTime != nil {
				delta := ppg[i].TimeStamp.Sub(*lastPeakTime).Milliseconds()
				if delta >= int64(minPeakDistanceMs) && delta > 0 {
					rr = append(rr, float64(delta))
					lastPeakTime = ppg[i].TimeStamp
				} else {
					fmt.Printf("Warning: Skipped peak at %v, delta %d ms too small or invalid\n", ppg[i].TimeStamp, delta)
				}
			} else {
				lastPeakTime = ppg[i].TimeStamp
			}
		}
	}

	if len(rr) == 0 {
		fmt.Println("Warning: No RR intervals detected. Check minPeakHeight or minPeakDistanceMs")
	}
	return rr
}


func CalculateFrequencyDomain(rr []float64) (lf, hf, vlf, lfHfRatio float64) {
	n := len(rr)
	if n < 2 {
		fmt.Println("Warning: Not enough RR intervals for frequency domain analysis.")
		return 0, 0, 0, 0
	}

	// 1. Estimate Instantaneous Heart Rate (IHR)
	ihr := make([]float64, n)
	for i, r := range rr {
		if r > 0 {
			ihr[i] = 1.0 / r // Frequency in Hz
		} else {
			ihr[i] = 0 // Handle zero RR intervals
		}
	}

	// 2. Estimate the sampling rate of the IHR signal
	var totalTime float64
	for _, r := range rr {
		totalTime += r
	}
	var samplingRate float64
	if totalTime > 0 && n > 1 {
		samplingRate = float64(n-1) / totalTime
	} else if n > 0 {
		samplingRate = 1.0 // Default if only one RR interval
	} else {
		fmt.Println("Warning: Cannot estimate sampling rate with no RR intervals.")
		return 0, 0, 0, 0
	}

	// 3. Detrend the IHR signal (remove the mean)
	meanIHR := 0.0
	validPoints := 0
	for _, val := range ihr {
		if val > 0 {
			meanIHR += val
			validPoints++
		}
	}
	if validPoints > 0 {
		meanIHR /= float64(validPoints)
	}
	detrendedIHR := make([]float64, n)
	for i, val := range ihr {
		if val > 0 {
			detrendedIHR[i] = val - meanIHR
		}
	}

	// 4. Apply a windowing function (e.g., Hamming) to reduce spectral leakage
	windowedIHR := make([]float64, n)
	for i := 0; i < n; i++ {
		hamming := 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		windowedIHR[i] = detrendedIHR[i] * hamming
	}

	// 5. Perform FFT
	complexResult := fft.Fft(complexify(windowedIHR), false)

	// 6. Calculate the Power Spectral Density (PSD)
	psd := make([]float64, n/2+1)
	for i := 0; i <= n/2; i++ {
		magnitudeSquared := real(complexResult[i])*real(complexResult[i]) + imag(complexResult[i])*imag(complexResult[i])
		psd[i] = magnitudeSquared / float64(n)
	}

	// 7. Define frequency bands (in Hz)
	vlfBand := struct{ low, high float64 }{low: 0.003, high: 0.04}
	lfBand := struct{ low, high float64 }{low: 0.04, high: 0.15}
	hfBand := struct{ low, high float64 }{low: 0.15, high: 0.4}

	// 8. Calculate frequency resolution
	freqResolution := samplingRate / float64(n)

	// 9. Integrate PSD within the frequency bands
	var vlfPower, lfPower, hfPower float64
	for i := 1; i <= n/2; i++ { // Start from 1 to avoid the DC component
		freq := float64(i) * freqResolution
		power := psd[i]

		if freq >= vlfBand.low && freq < vlfBand.high {
			vlfPower += power
		}
		if freq >= lfBand.low && freq < lfBand.high {
			lfPower += power
		}
		if freq >= hfBand.low && freq < hfBand.high {
			hfPower += power
		}
	}

	lf = lfPower
	hf = hfPower
	vlf = vlfPower
    // fmt.Println("LF:", lf, "HF:", hf, "VLF:", vlf)
	if hf > 0 {
		lfHfRatio = lf / hf
	} else {
		lfHfRatio = 0.0
	}

	return lf, hf, vlf, lfHfRatio
}

func complexify(real []float64) []complex128 {
	complexData := make([]complex128, len(real))
	for i, r := range real {
		complexData[i] = complex(r, 0)
	}
	return complexData
}

// func main() {
// 	// Example usage with a synthetic RR interval series (in seconds)
// 	rrIntervals := []float64{0.8, 0.82, 0.78, 0.85, 0.81, 0.79, 0.83, 0.8, 0.82, 0.78}
// 	lf, hf, vlf, lfHf := CalculateFrequencyDomainPrecise(rrIntervals)
// 	fmt.Printf("Precise LF Power: %.4f\n", lf)
// 	fmt.Printf("Precise HF Power: %.4f\n", hf)
// 	fmt.Printf("Precise VLF Power: %.4f\n", vlf)
// 	fmt.Printf("Precise LF/HF Ratio: %.4f\n", lfHf)

// 	// Example with more data
// 	longRR := make([]float64, 100)
// 	for i := 0; i < 100; i++ {
// 		longRR[i] = 0.8 + 0.05*math.Sin(2*math.Pi*float64(i)/50) + 0.02*math.Sin(2*math.Pi*float64(i)/10)
// 	}
// 	lfLong, hfLong, vlfLong, lfHfLong := CalculateFrequencyDomainPrecise(longRR)
// 	fmt.Printf("\nPrecise Long RR - LF Power: %.4f\n", lfLong)
// 	fmt.Printf("Precise Long RR - HF Power: %.4f\n", hfLong)
// 	fmt.Printf("Precise Long RR - VLF Power: %.4f\n", vlfLong)
// 	fmt.Printf("Precise Long RR - LF/HF Ratio: %.4f\n", lfHfLong)
// }