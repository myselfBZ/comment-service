package main

import (
	"log"
	"github.com/IBM/sarama"
)

const( 
    BROKER_URL = "localhost:9092"
    TOPIC = "comments"
)

var CONSUMER_CONN sarama.PartitionConsumer

type Consumer struct{
    err     chan error
    msg     chan string
    conn    sarama.PartitionConsumer 
}


func(c *Consumer) run() {

    for{
        select{
            case msg := <- c.msg:
                log.Println("message recieved: ", msg)
            case err := <- c.err:
                log.Println("error recieved: ", err)
        }
    }
}

func (c *Consumer) errorLoop(){
    log.Print("error loop is working")
    for err := range c.conn.Errors(){
        c.err <- err.Err
    }
}

func (c *Consumer) messageLoop(){
    log.Print("message loop is working")
    for msg := range c.conn.Messages(){
        log.Print("waiting for ya")
        c.msg <- string(msg.Value)
    }
}

func NewConsumer(conn sarama.PartitionConsumer) *Consumer{
    return &Consumer{
        conn: conn,
        err: make(chan error),
        msg: make(chan string),
    }
}

func init(){
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    conn, err := sarama.NewConsumer([]string{BROKER_URL}, config)
    if err != nil{
        log.Fatal("error connecting to shit: ", err)
    }
    partConsumer, err := conn.ConsumePartition(TOPIC, 0, sarama.OffsetOldest)
    if err != nil{
        log.Fatal("error creating consumer partition: ", err)
    }
    CONSUMER_CONN = partConsumer
}

func main(){
    c := NewConsumer(CONSUMER_CONN)
    done := make(chan struct{})
    go c.run()
    go c.errorLoop()
    go c.messageLoop()   
    <- done
    log.Print("consumer started")
}
