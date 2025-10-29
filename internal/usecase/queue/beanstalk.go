package queue

import (
	"log"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type BeanstalkProducer struct {
	Conn *beanstalk.Conn
}

func NewBeanstalkProducer(addr string) (*BeanstalkProducer, error) {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &BeanstalkProducer{
		Conn: conn,
	}, nil
}

func (b *BeanstalkProducer) AddJob(data []byte, tubeName string) error {
	tube := &beanstalk.Tube{
		Conn: b.Conn,
		Name: tubeName,
	}
	_, err := tube.Put(data, 1, 0, 120*time.Second)
	if err != nil {
		return err
	}
	log.Printf("job added to tube %s\n", tube.Name)
	return nil
}
