package main

import (
  // "fmt"
  "flag"
  // "log"
  "encoding/json"

  "github.com/kitacode/gomon"
  "github.com/kitacode/go-sync-store/models"
  "github.com/nsqio/go-nsq"
  "github.com/getsentry/raven-go"
)

func main() {
  raven.SetDSN("http://76ac8462c1794c8bb3a3898311f75324:79cc0177e8e04e719a8b2ea70c74e5a7@localhost:9000/4")

  mode := flag.String("m", "marketplace", "mode of pusher")
  flag.Parse()

  switch *mode {
    case "marketplace": {
      nsqConfig := nsq.NewConfig()
      w, errNsq := nsq.NewProducer("127.0.0.1:4150", nsqConfig)
      if errNsq != nil {
        raven.CaptureErrorAndWait(errNsq, nil)
      }

      mongoConn := new(gomon.Mongo)
      mongoConn.Init([]string{"127.0.0.1:27017",}, "mystore", "marketplaces")

      marketplaces := mongoConn.Find().Iter()
      marketplace := models.Marketplace{}

    	for marketplaces.Next(&marketplace) {
        marketplaceStr, err := json.Marshal(marketplace)
        if err != nil {
          raven.CaptureErrorAndWait(err, nil)
        }

        err = w.Publish("marketplaces", []byte(string(marketplaceStr)))
        if err != nil {
          raven.CaptureErrorAndWait(err, nil)
        }
    	}

      w.Stop()
    }
  }
}
