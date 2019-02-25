package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math"
	"net/http"
	"strconv"
)

func main() {
	collector := AmatenCollector{}
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9042", nil)
	if err != nil {
		panic(err)
	}
}

type AmatenCollector struct {
}

func (AmatenCollector) Describe(chan<- *prometheus.Desc) {
}

func (AmatenCollector) Collect(ch chan<- prometheus.Metric) {
	for _, giftType := range []string{"amazon", "itunes", "google_play"} {
		res, err := GetPrise(giftType, 20)
		if err != nil {
			continue
		}

		allGiftCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "custom_amaten_all_gift_count",
			Help: "Number of listings",
			ConstLabels: map[string]string{
				"gift": giftType,
			},
		})
		allGiftCount.Set(float64(res.AllGiftCount))
		ch <- allGiftCount

		minRateGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "custom_amaten_rate_min",
			Help: "Minimum rate",
			ConstLabels: map[string]string{
				"gift": giftType,
			},
		})
		byRate100kGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "custom_amaten_rate_100k",
			Help: "Average rate when you bought 100,000 yen",
			ConstLabels: map[string]string{
				"gift": giftType,
			},
		})
		minRate := 100.0
		faceSum := 0
		priceSum := 0
		for i := range res.Gifts {
			rate, _ := strconv.ParseFloat(res.Gifts[i].Rate, 64)
			minRate = math.Min(minRate, rate)
			if faceSum < 100000 {
				faceSum += res.Gifts[i].FaceValue * res.Gifts[i].Cnt
				priceSum += res.Gifts[i].Price * res.Gifts[i].Cnt
			}
		}
		minRateGauge.Set(minRate / 100)
		byRate100kGauge.Set(float64(priceSum) / float64(faceSum))
		ch <- minRateGauge
		ch <- byRate100kGauge
	}
}
