package core

import (
	"testing"
	"fmt"
)

// --- TestManualWiredSingleTon ---
type IService interface {
	great(name string)
	get() string
}

type Service struct {
	Word string
}

func (s *Service) great(name string) {
	s.Word = name
}
func (s Service) get() string {
	return s.Word

}

func NewService() interface{} {
	return &Service{}
}

type IConsumer interface {
	call() (string)
	setService(srv IService)
}

type Consumer struct {
	GreatService IService
}

func (c Consumer) call() string {
	return "Hi!" + c.GreatService.get()
}

func (c *Consumer) setService(srv IService) {
	c.GreatService = srv
}

func NewConsumer() interface{} {
	return &Consumer{}
}

var gimanager GInstanceManager

func TestManualWiredSingleTon(t *testing.T) {
	consumer := gimanager.Get("Consumer").(IConsumer) 
	great := gimanager.Get("Service").(IService)
	consumer.setService(great)
	great.great(" Yuki!")
	str := consumer.call()
	if str != "Hi! Yuki!" {
		t.Error("Not equal 'Hi! Yuki!' : '%v%v',", "Hi!", great.get())
	}
	consumer2 := gimanager.Get("Consumer").(IConsumer)
	great2 :=  gimanager.Get("Service").(IService)
	if great2.get() != " Yuki!" {
		t.Error("Not equal ' Yuki!' : '%v'", great2.get())
	}
	str2 := consumer2.call()
	if str2 != "Hi! Yuki!" {
		t.Error("2nd cosumer great Not equal 'Hi! Yuki!' : '%v%v',", "Hi!", great2.get())
	}
}

// --- TestAutoWiredSingleton ---

type IWireConsumer interface{
	call() (string)
	getService() (IService)
}

type WireConsumer struct{
	GreatService IService `wired:"Service"`
}

func (c WireConsumer) call() string {
	return "Ya!!!" + c.GreatService.get()
}

func (c WireConsumer) getService() IService{
	return c.GreatService
}

func NewWireConsumer() interface{} {
	return &WireConsumer{}
}

func  TestAutoWiredSingleton(t *testing.T){
	consumer := gimanager.Get("WireConsumer").(IWireConsumer)
	consumer.getService().great(" Akiko!!!")
	str := consumer.call()
	if str != "Ya!!! Akiko!!!" {
		t.Error("Not equal 'Ya!!! Akiko!!!' : %v",str)
	}
	consumer2 := gimanager.Get("Consumer").(IConsumer)
	str2 := consumer2.call()
	if str2 != "Hi! Akiko!!!" {
		t.Error("Not equal 'Hi! Akiko!!!' : '%v'", str2)
		t.Error("consumer.service:",consumer)
		t.Error("consumer.service:",consumer2)
	}
}

// --- TestAutoInjectPrototype　---
type IInjectConsumer interface{
	call() (string)
	getService() (IService)
}

type InjectConsumer struct{
	GreatService IService `inject:"Service"`
}

func (c InjectConsumer) call() string{
	return "Hello!!!"+c.GreatService.get()
}

func (c InjectConsumer) getService() IService{
	return c.GreatService
}

func NewInjectConsumer() interface{} {
	return &InjectConsumer{}
}

func TestAutoInjectPrototype(t *testing.T) {
	consumer := gimanager.New("InjectConsumer").(IInjectConsumer)
	consumer.getService().great(" Motoki!!!")
	str := consumer.call()
	if str!= "Hello!!! Motoki!!!" {
		t.Error("Not equal 'Hello!!! Motoki!!!:", str)	
	}
	consumer2 := gimanager.New("InjectConsumer").(IInjectConsumer)
	consumer2.getService().great(" Kenta!!!")
	str2 := consumer2.call()
	if str2 != "Hello!!! Kenta!!!" {
		t.Error("Not equal 'Hello!!! Kenta!!!:", str2)
	}
	str3 := consumer.call()
	if str3 != str {
		t.Error("str != str3 :", str, ":", str3)
	}
}

// --- TestDuplicateFactoryName ---
func TestDuplicateFactoryName(t *testing.T) {
	defer func(){
		if err := recover(); err != nil {
			if err == "GIManager.Duplicate.Factory" {
				fmt.Println("Panic happened:", err)
				return
			} else {
				t.Error("Wrong Panic!:", err)
				return
			}
        	}
		t.Error("Not happen panic!")
		return
	}()
	gimanager.AddFactoryMethod("Consumer", NewConsumer);
	t.Error("Not happen panic!")
}

// -- TestInjectProxy --
type IProxyService interface {
	great(name string)
	get() string
}

type ProxyService struct {
	Word string
}

func (s *Service) great(name string) {
	s.Word = name
}
func (s Service) get() string {
	return s.Word

}

func NewProxyService() interface{} {
	return &Service{}
}

type IProxyConsumer interface {
	call() (string)
}

type ProxyConsumer struct {
	GreatService IProxyService `proxy:"Service"`
}

func (c ProxyConsumer) call() string {
	return "Hey!" + c.GreatService.get()
}


func NewProxyConsumer() interface{} {
	return &ProxyConsumer{}
}
func TestInjectProxy( t *testing.T ){
	
}

// GInstanceManagerの初期化とファクトリーメソッドの登録
func init() {
	gimanager = GetGInstanceManager()
	gimanager.AddFactoryMethod("Consumer", NewConsumer)
	gimanager.AddFactoryMethod("Service", NewService)
        gimanager.AddFactoryMethod("WireConsumer", NewWireConsumer)
	gimanager.AddFactoryMethod("InjectConsumer", NewInjectConsumer)
}

