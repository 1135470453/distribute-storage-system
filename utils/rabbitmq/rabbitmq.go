package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	//s 为mq的virtualhost
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}
	//创建Channel
	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}
	//创建queue,之后向这个queue发送消息
	q, e := ch.QueueDeclare(
		"",
		false,
		true, //autoDelete:当有消费者连接后，所有消费者都选择断开时，该queue自动删除
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name //q不是没有name吗?
	return mq
}

//将队列与exchange绑定,routingkey为空
func (q *RabbitMQ) Bind(exchange string) {
	e := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil)
	if e != nil {
		panic(e)
	}
	q.exchange = exchange
}

//queue表示要发往的消费者
func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

// Consume q从绑定的exchange获取信息
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, e := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	return c
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
}
