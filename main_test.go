package main

import (
	"encoding/json"
	"testing"

	"go_json_vs_proto/pb"

	"google.golang.org/protobuf/proto"
	"github.com/gomodule/redigo/redis"
)

type UserJson struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Address        string  `json:"address"`
	Pincode        string  `json:"pincode"`
	IsValid        bool    `json:"is_valid"`
	AccountBalance float64 `json:"account_balance"`
}

const userJsonKey = "user_json"
const userProtoKey = "user_proto"

func BenchmarkSetJsonStruct(b *testing.B) {
	pool := newPool()
	client := pool.Get()
	defer client.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newUser := &UserJson{
			FirstName:      "Jhon",
			LastName:       "Doe",
			Address:        "221B Baker Street, London",
			Pincode:        "NW1",
			IsValid:        true,
			AccountBalance: 109084493,
		}
		userJson, err := json.Marshal(newUser)
		if err != nil {
			b.Log("error ", err)
		}

		_, err = client.Do("SET", userJsonKey, string(userJson))
		if err != nil {
			b.Log("error ", err)
		}
	}

}

func BenchmarkSetProto(b *testing.B) {
	pool := newPool()
	client := pool.Get()
	defer client.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newUser := &pb.UserProto{
			FirstName:      "Jhon",
			LastName:       "Doe",
			Address:        "221B Baker Street, London",
			Pincode:        "NW1",
			IsValid:        true,
			AccountBalance: 109084493,
		}
		userProto, err := proto.Marshal(newUser)
		if err != nil {
			b.Log("error ", err)
		}

		_, err = client.Do("SET", userProtoKey, string(userProto))
		if err != nil {
			b.Log("error ", err)
		}
	}

}

func BenchmarkGetJsonStruct(b *testing.B) {
	pool := newPool()
	client := pool.Get()
	defer client.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bytes, err := redis.Bytes(client.Do("GET", userJsonKey))
		if err != nil {
			b.Log("error ", err)
		}

		newUser := &UserJson{}
		err = json.Unmarshal(bytes, newUser)
		if err != nil {
			b.Log("error ", err)
		}
	}

}

func BenchmarkGetProto(b *testing.B) {
	pool := newPool()
	client := pool.Get()
	defer client.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
	

		bytes, err := redis.Bytes(client.Do("GET", userProtoKey))
		if err != nil {
			b.Log("error ", err)
		}

		newUser := &pb.UserProto{}
		err = proto.Unmarshal(bytes, newUser)
		if err != nil {
			b.Log("error ", err)
		}
	}

}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
