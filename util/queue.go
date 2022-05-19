package util

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"strconv"
	"syscall"
)

const (
	segmentMaxSize = 1024 * 1024 * 4
	//taskMaxSize    = 1024 * 4
)

var (
	prefix = "/tmp/"
)

type index struct {
	a uint16 //文件名index
	b uint32 //文件内的index
}

type DurableQueue struct {
	readCh     chan []byte
	writeCh    chan []byte
	ctx        context.Context
	cancelFunc context.CancelFunc
	name       string
	meta       []byte
	curRead    []byte
	curWrite   []byte
	consumer   index
	producer   index
	//writeLock    *sync.Mutex
	//readLock     *sync.Mutex
	producerFile *os.File
	consumerFile *os.File
}

func (d *DurableQueue) Poll() (rt []byte, err error) {
	select {
	case rt = <-d.readCh:
	default:
	}
	if len(rt) == 0 {
		err = errors.New("")
	}
	return
}

func (d *DurableQueue) poll() {
	//d.readLock.Lock()
	//defer d.readLock.Unlock()
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			//返回任务之后再移动index
			if d.checkForConsumer() {
				var l uint32
				buffer := bytes.NewBuffer(d.curRead[d.consumer.b:])
				err := binary.Read(buffer, binary.BigEndian, &l)
				if err != nil {
					log.Printf("read data block %d error at %d", d.consumer.a, d.consumer.b)
					log.Fatalln(err)
				}
				rt := make([]byte, l)
				copy(rt, d.curRead[d.consumer.b+4:])
				d.consumer.b += 4 + l
				d.readCh <- rt
				d.setConsumerIndex()
			} else {
				//err = errors.New("temporally nothing new")
				//暂时没有新的task
			}
		}
	}
}

func (d *DurableQueue) Offer(b []byte) {
	d.writeCh <- b
}

func (d *DurableQueue) offer() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case b := <-d.writeCh:
			if len(b) == 0 {
				continue
			}
			//d.writeLock.Lock()
			//defer d.writeLock.Unlock()
			d.checkForProducer(b)
			l := uint32(len(b))
			var buffer bytes.Buffer
			//d.curWrite[d.producer.b:]
			err := binary.Write(&buffer, binary.BigEndian, l)
			if err != nil {
				log.Fatalln(err)
			}
			temp := buffer.Bytes()
			//println(len(temp))
			//先写内容再写长度 避免消费者读取到不完整消息
			copy(d.curWrite[d.producer.b+4:], b)
			copy(d.curWrite[d.producer.b:], temp)
			//buffer.Write(b)
			d.producer.b += 4 + l
			d.setProducerIndex()
			//fmt.Printf("%x\n", d.getProducerIndex())
		}
	}
}

func (d *DurableQueue) load() error {
	//打开meta
	meta, err := os.OpenFile(prefix+d.name+".meta", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_ = meta.Truncate(12)
	if err != nil {
		return err
	}
	d.meta, err = syscall.Mmap(int(meta.Fd()), 0, 12, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	d.consumer = d.getConsumerIndex()
	d.producer = d.getProducerIndex()
	//打开data
	d.consumerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.consumer.a))+".data", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = d.consumerFile.Truncate(segmentMaxSize)
	if err != nil {
		return err
	}
	d.curRead, err = syscall.Mmap(int(d.consumerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	d.producerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.producer.a))+".data", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = d.producerFile.Truncate(segmentMaxSize)
	if err != nil {
		return err
	}
	d.curWrite, err = syscall.Mmap(int(d.producerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	return nil
}

func (d *DurableQueue) Test() {
	//d.meta[100] = 'n'
}

func (d *DurableQueue) Stop() {

}

func (d *DurableQueue) getConsumerIndex() index {
	buffer := bytes.NewBuffer(d.meta)
	var a uint16
	_ = binary.Read(buffer, binary.BigEndian, &a)
	var b uint32
	_ = binary.Read(buffer, binary.BigEndian, &b)
	return index{a, b}
}

func (d *DurableQueue) getProducerIndex() index {
	buffer := bytes.NewBuffer(d.meta[6:])
	var a uint16
	_ = binary.Read(buffer, binary.BigEndian, &a)
	var b uint32
	_ = binary.Read(buffer, binary.BigEndian, &b)
	return index{a, b}
}

func (d *DurableQueue) setConsumerIndex() {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, d.consumer.a)
	_ = binary.Write(&buffer, binary.BigEndian, d.consumer.b)
	//d.readLock.Lock()
	//defer d.readLock.Unlock()
	copy(d.meta, buffer.Bytes())
}

func (d *DurableQueue) setProducerIndex() {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, d.producer.a)
	_ = binary.Write(&buffer, binary.BigEndian, d.producer.b)
	//d.writeLock.Lock()
	//defer d.writeLock.Unlock()
	copy(d.meta[6:], buffer.Bytes())
}

func (d *DurableQueue) checkForProducer(b []byte) {
	//println(segmentMaxSize - int(d.producer.b) - 4 - len(b))
	if int(d.producer.b)+4+len(b) >= segmentMaxSize {
		if d.producer.b+4 >= segmentMaxSize {
			//啥也不做
			//log.Println("left < 4 bytes")
		} else {
			//标记一下是补充的用于对其的字节
			//（mmap对其有用吗？？）
			//log.Println("left > 4 bytes")
			//log.Println("mark as end")
			var buffer bytes.Buffer
			var n uint32 = 0xFFFFFFFF
			_ = binary.Write(&buffer, binary.BigEndian, &n)
			copy(d.curWrite[d.producer.b:], buffer.Bytes())
		}
		err := syscall.Munmap(d.curWrite)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("unmap %d.data success", d.producer.a)

		err = d.producerFile.Close()
		if err != nil {
			log.Printf("producer close block %d fail", d.producer.a)
			log.Println(err)
		}

		d.producer.a++
		d.producer.b = 0
		d.producerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.producer.a))+".data", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("produce new segment %d", d.producer.a)
		err = d.producerFile.Truncate(segmentMaxSize)
		if err != nil {
			log.Fatalln(err)
		}
		d.curWrite, err = syscall.Mmap(int(d.producerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	}
}

/**
返回true：还有task
返回false：没有task了
*/
func (d *DurableQueue) checkForConsumer() bool {
	if d.consumer.b+4 >= segmentMaxSize {
		//换下一块
		err := syscall.Munmap(d.curRead)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("unmap %d.data success", d.consumer.a)
		err = d.consumerFile.Close()
		if err != nil {
			log.Printf("consumer close block %d fail", d.consumer.a)
			log.Println(err)
		}
		d.deleteOutdated()
		d.consumer.a++
		d.consumer.b = 0
		d.consumerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.consumer.a))+".data", os.O_RDWR, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("consumer new segment %d", d.consumer.a)
		err = d.consumerFile.Truncate(segmentMaxSize)
		if err != nil {
			log.Fatalln(err)
		}
		d.curRead, err = syscall.Mmap(int(d.consumerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	}
	//读取task长度
	buffer := bytes.NewBuffer(d.curRead[d.consumer.b:])
	var l uint32
	err := binary.Read(buffer, binary.BigEndian, &l)
	if err != nil {
		log.Printf("read data block %d error at %d", d.consumer.a, d.consumer.b)
		log.Fatalln(err)
	}
	//此处有3种情况
	//1是这个块到底了(读到长度是0xFFFFFFFF) 2是暂时没有新的可以消费(读到长度为0) 3正常往下读
	switch l {
	case 0xFFFFFFFF:
		err := syscall.Munmap(d.curRead)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("unmap %d.data success", d.consumer.a)
		d.deleteOutdated()
		d.consumer.a++
		d.consumer.b = 0
		d.consumerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.consumer.a))+".data", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("consumer open new block %d", d.consumer.a)
		err = d.consumerFile.Truncate(segmentMaxSize)
		if err != nil {
			log.Fatalln(err)
		}
		d.curRead, err = syscall.Mmap(int(d.consumerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		return d.checkForConsumer()
	case 0:
		//log.Println("nothing left")
		return false
	default:
		return true
	}
	/*if l == 0 {
		//换下一块
		err := syscall.Munmap(d.curRead)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("unmap %d.data success", d.consumer.a)
		d.deleteOutdated()
		d.consumer.a++
		d.consumer.b = 0
		d.consumerFile, err = os.OpenFile(prefix+strconv.Itoa(int(d.consumer.a))+".data", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("consumer open new block %d", d.consumer.a)
		err = d.consumerFile.Truncate(segmentMaxSize)
		if err != nil {
			log.Fatalln(err)
		}
		d.curRead, err = syscall.Mmap(int(d.consumerFile.Fd()), 0, segmentMaxSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		return d.checkForConsumer()
		//32个二进制的1
	} else if l == 4294967295 {
		//没有task了
		log.Println("nothing left")
		return false
	}*/
}

func (d *DurableQueue) deleteOutdated() {
	log.Printf("delete block %d", d.consumer.a)
	err := os.Remove(prefix + strconv.Itoa(int(d.consumer.a)) + ".data")
	if err != nil {
		log.Println(err)
	}
}

func NewQueue(name string) (*DurableQueue, error) {
	q := &DurableQueue{
		readCh:  make(chan []byte),
		writeCh: make(chan []byte),
		name:    name,
		//readLock:  &sync.Mutex{},
		//writeLock: &sync.Mutex{},
	}
	q.ctx, q.cancelFunc = context.WithCancel(context.Background())
	if err := q.load(); err != nil {
		return nil, err
	}
	/*	a, b := q.getProducerIndex()
		println(a)
		println(b)*/
	go q.offer()
	go q.poll()
	return q, nil
}
