package main 

import (
    "context"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
  )
func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}
func main() {
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
failOnError(err, "Failed to connect to RabbitMQ")
defer conn.Close()
ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()
q, err := ch.QueueDeclare(
  "hello", // name
  false,   // durable
  false,   // delete when unused
  false,   // exclusive
  false,   // no-wait
  nil,     // arguments
)
failOnError(err, "Failed to declare a queue")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

app := fiber.New()

app.Use(logger.New(),)
app.Get("/send", func(c *fiber.Ctx) error {
        // Create a message to publish.
        message := amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(c.Query("msg")),
        }
err = ch.PublishWithContext(ctx,
  "",     // exchange
  q.Name, // routing key
  false,  // mandatory
  false,  // immediate
  message,  
)
failOnError(err, "Failed to publish a message")
return nil
})
log.Fatal(app.Listen(":3000"))
}

