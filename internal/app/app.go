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

	//3) Lists
	//pushing a new element into a list at the head(left) - lpush
	result := app.Redis.Lpush("work:queue", "lpush")
	log.Println("lpush result - number of elements in the list", result)

	// pushing a new element into a list at the tail(right) - rpush
	result = app.Redis.Rpush("work:queue", "rpush")
	log.Println("rpush result - number of elements in the list", result)

	// popping an element from a head(left side) of list - lpop
	resultPop := app.Redis.Lpop("work:queue")
	log.Println("lpop result - popped element of the list", resultPop)

	//we also can pushing multiple elements
	//todo example

	// popping an element from a tail(right sife) of list - rpop
	resultPop = app.Redis.Rpop("work:queue")
	log.Println("rpop result - popped element of the list", resultPop)

	// getting all list data - lrange
	rangeResulst := app.Redis.Lrange("work:queue", 0, -1)
	log.Println("lrange result - slice of elements in the list", rangeResulst)

	//if we trying to pop element from an empty list, result will be null
	//todo another example

	//if we want to store element as a capped collection - ltrim. ltrim is similar to lrange,
	//but instead of displaying elements of specified range, it sets this range as the new list value
	// elements in the list now: [lpush test rpush], so lets try to leave only the first element, and to delete rest
	lTrimStatus := app.Redis.Ltrim("work:queue", 0, 0)
	log.Println("Result of ltrim", lTrimStatus)
	//and now lets try to see that list
	rangeResulst = app.Redis.Lrange("work:queue", 0, -1)
	log.Println("Lrange after ltrim result:", rangeResulst)

	//try it again
	// for a pattern pubsub redis implememnts command brpop and blpop
	// blpop - popping element from head of list, if there are no elements, it will wait timeout.. .
	result = app.Redis.Lpush("work:queue:2", "value") //adding a new element into list
	blPopResult := app.Redis.BLPop(5, []string{"work:queue:2"})
	log.Println("BLPop result:", blPopResult)
	// brpop - popping element from tail of list, if there are no elements, it will wait timeout..
	//todo example

	//todo automatic creation and removal of keys

	//-----------
	//3) hashes -
	// hset - sets multiple(o one) fields of the hash
	hSetResult := app.Redis.HSet("hashed", "one", "two", "free", 12)
	log.Println("HSet result", hSetResult)

	// hget - gets values of set
	hGetResult := app.Redis.HGet("hashed", "one")
	log.Println("HGet result", hGetResult)

	return nil
}
