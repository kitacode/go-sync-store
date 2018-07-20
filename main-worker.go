package main

import (
  "log"
  "flag"
  "sync"
  "encoding/json"

  "github.com/nsqio/go-nsq"
  "github.com/getsentry/raven-go"

  "github.com/kitacode/go-sync-store/models"
  "github.com/kitacode/go-store/shopify"
)

func main() {
  raven.SetDSN("http://76ac8462c1794c8bb3a3898311f75324:79cc0177e8e04e719a8b2ea70c74e5a7@localhost:9000/4")

  mode := flag.String("m", "new-orders", "mode of worker")
  flag.Parse()

  switch *mode {
    case "new-orders": {
      log.Println("new-orders")
    }
  }

  wg := &sync.WaitGroup{}
  wg.Add(1)

  decodeConfig := nsq.NewConfig()

  c, err := nsq.NewConsumer("marketplaces", "new-orders", decodeConfig)
  if err != nil {
    raven.CaptureErrorAndWait(err, nil)
  }

  marketplace := models.Marketplace{}

  c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
    err = json.Unmarshal([]byte(message.Body), &marketplace)
    if err != nil {
      raven.CaptureErrorAndWait(err, nil)
    }

    switch marketplace.Id {
      case "shopify": {
        spf := shopify.Shopify{
          Apikey: "86d42b670bd5822538b165b8d8846e8b",
          ApiPassword: "a9e28fc74fdef411d823bcbcc2207c43",
        }
        orders := spf.GetOrders()
        a, err := json.Marshal(orders)
        if err != nil {
          raven.CaptureErrorAndWait(err, nil)
        }
        log.Println(string(a))
      }
    }

    return nil
  }))

  err = c.ConnectToNSQD("127.0.0.1:4150")
  if err != nil {
    raven.CaptureErrorAndWait(err, nil)
  }
  wg.Wait()
}
