package app

import (
	"Redis/pkg/redis"
	"log"
	"strconv"
	"time"
)

type App struct {
	Version string
	Redis   *redis.Redis
}

func New() *App {
	var app *App = new(App)

	app.Version = "0.9.1"
	app.Redis = redis.New("localhost:6389", "", 0)
	return app
}

func (app *App) Run() error {
	//Redis data types:
	//1) Strings - simpliest type of value what can be associated with a Redis key
	//setting value by key - set
	_, err := app.Redis.Exec("set", "firstVar", 10)
	if err != nil {
		return err
	}
	//getting value by key - get
	val, err := app.Redis.Exec("get", "firstVar", nil)
	if err != nil {
		return err
	}

	log.Println(val) //10

	//incrementing value - incr
	val, err = app.Redis.Exec("incr", "firstVar", nil)
	if err != nil {
		return err
	}
	log.Println("Variable after increment", val)

	//decrementing value - decr
	val, err = app.Redis.Exec("decr", "firstVar", nil)
	if err != nil {
		return err
	}
	log.Println("Variable after decrement", val)

	var recordTtl int = 4 //4 seconds of ttl

	//setting tempVariable with ttl
	_ = app.Redis.Set("tempVariable", "It will be destroyed after "+strconv.Itoa(recordTtl)+" seconds", 4)
	if err != nil {
		return err
	}
	//getting ttl of tempVariable
	ttl, err := app.Redis.Exec("ttl", "tempVariable", nil)
	if err != nil {
		return err
	}
	log.Println("ttl of tempVariable:", ttl)
	//getting existing tempVariable
	val = app.Redis.Get("tempVariable")
	log.Println("tempVariable", val)
	//getting expired tempVariable
	time.Sleep(time.Duration(recordTtl) * time.Second)
	val = app.Redis.Get("tempVariable")
	log.Println("Expired temp variable", val)

	//getting value and setting new one - getset
	oldVal, err := app.Redis.Exec("getset", "firstVar", 100500)
	log.Println("Old value of firstVar", oldVal)
	log.Println("New value of firstVar", app.Redis.Get("firstVar"))

	//setting multiple value - mset
	_, err = app.Redis.MSet([]string{"key1", "key2", "key3"}, []string{"value1", "value2", "value4"})
	if err != nil {
		return err
	}

	//get multiple variables
	log.Println("Getting multiple values", app.Redis.MGet([]string{"key1", "key2", "key3"}...))

	//--------------------
	//2) Working with keys (it s not a type but it helps to interact with keys)
	//checking key existance - exists
	app.Redis.Set("testExistance", "exists", 0)
	val, err = app.Redis.Exec("exists", "testExistance", nil)
	log.Println("Checking key existance", val)

	//deleting value by key - del
	app.Redis.Exec("del", "testExistance", nil)
	val, err = app.Redis.Exec("exists", "testExistance", nil)
	log.Println("Checking key existance after removing", val)

	//getting type of value stored by key
	app.Redis.Set("intVar", 123, 0) //type is string
	val, err = app.Redis.Exec("type", "intVar", nil)
	log.Println("Getting type of int var", val)

	//getting key list by pattern (here - *)
	keys, err := app.Redis.Exec("keys", "*", nil)
	if err != nil {
		return err
	}
	log.Println("all keys", keys)
	return nil

	//3) Lists
	//pushing a new element into a list at the head(left) - lpush
}
