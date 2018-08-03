package redistore

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const MaxLength = 4096

// saves value by key | SETEX command
var Setex = func(key, value string, age int) error {
	if MaxLength != 0 && len(value) > MaxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", key, age, value)
	return err
}

// loads value by key | GET command
// returns empty string if there was an error
var GET = func(key string) (string, error) {
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return "", err
	}
	data, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil // no data was associated with this key
	}
	b, err := redis.Bytes(data, err)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// delete removes keys from redis | DEL command
var Del = func(key string) error {
	conn := store.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

// get keys from redis by pattern | KEYS command
var Keys = func(pattern string) (arr []string, err error) {
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	res, err := conn.Do("KEYS", pattern)
	redisResult, ok := res.([]interface{})
	if !ok {
		return nil, nil
	}
	for _, item := range redisResult {
		switch reflect.TypeOf(item).Kind() {
		case reflect.Slice:
			{
				s := reflect.ValueOf(item)
				result := []string{}
				for i := 0; i < s.Len(); i++ {
					result = append(result, string(s.Index(i).Interface().(uint8)))
				}
				key := strings.Join(result, "")
				arr = append(arr, key)
			}
		default:
			return arr, err
		}
	}
	return arr, err
}

// get keys from redis by pattern | scan command
// WARNING! this function only makes up to 1000 iterations (so if you need to exctract more than 4k values at one call please modify maxIterations)
// @deprecated. Use Keys instead
var Scan = func(pattern string) ([]string, error) {
	conn := store.Pool.Get()
	defer conn.Close()

	iteratorId := 0
	var result []string

	//this iteration is only to prevent possible infinite loop.
	maxIterations := 1000
	for noOfGetRequests := 0; noOfGetRequests < maxIterations; noOfGetRequests++ {
		//scan not always return all the results in one request.
		//scan returns [ [iteratorId] [ [][][]... ] ]
		//where iteratorId is a number(link) to next set of result, and [ [][][]... ]  - array of array of bytes
		res, err := conn.Do("SCAN", iteratorId, "match", pattern)
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, errors.New("GetByPattern: no data was associated with this key")
		}
		redisResult, ok := res.([]interface{})
		if !ok {
			return nil, errors.New("GetByPattern: can not convert interface to processable results")
		}
		if len(redisResult) <= 1 {
			return nil, errors.New("GetByPattern: not expected output from redis")
		}
		newIterator := redisResult[0].([]uint8)
		iteratorId, err = strconv.Atoi(string(newIterator))
		if err != nil {
			return nil, errors.New("GetByPattern: can not get iterator")
		}

		for _, item := range redisResult[1:] {
			switch reflect.TypeOf(item).Kind() {
			case reflect.Slice:
				{
					s := reflect.ValueOf(item)
					for i := 0; i < s.Len(); i++ {
						result = append(result, string(s.Index(i).Interface().([]uint8)))
					}
				}
			default:
				return nil, errors.New("GetByPattern: not supported data type")
			}
		}
		if iteratorId == 0 {
			return result, nil
		}
	}
	return nil, errors.New("too many request to redis detected. Check redis for pattern " + pattern + " or increase maxIterations")
}
