package redistore

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var store, storeDeprecating *RediStore

func connect(newStore *RediStore) {
	r, err := NewRediStore(10, "tcp", defaultRedisHost+":"+defaultRedisPort, "", []byte(""))
	if err != nil {
		fmt.Println(err)
		panic("can not connect to redis")
	}
	r.SetMaxAge(30 * 24 * 3600)
	store = r
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

//generate random string based on latin alphabet
var randString = func(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

var _ = Describe("deprecated store cookie options", func() {
	connect(store)

	Context("redis SETEX, GET, DEL, SCAN, KEYS commands", func() {
		It("should not let store too big value", func() {
			key := randString(10)
			value := randString(MaxLength + 1)
			age := 20
			err := Setex(key, value, age)
			Expect(err).NotTo(BeNil())
		})
		It("should set value for 1 sec", func() {
			key := randString(10)
			value := randString(10)
			//saving value for 1 sec
			age := 1
			err := Setex(key, value, age)
			Expect(err).To(BeNil())

			//check if it was saved
			newValue := ""
			newValue, err = Get(key)
			Expect(err).To(BeNil())
			Expect(newValue).To(BeEquivalentTo(value))

			//sleep for 2 sec
			time.Sleep(2 * time.Second)
			newValue, err = Get(key)
			Expect(err).To(BeNil())
			Expect(newValue).NotTo(BeEquivalentTo(value))
		})
		It("should get value", func() {
			key := randString(10)
			value := randString(10)
			err := Setex(key, value, 20)
			Expect(err).To(BeNil())

			newValue := ""
			newValue, err = Get(key)
			Expect(err).To(BeNil())
			Expect(newValue).To(BeEquivalentTo(value))
		})
		It("should get values by pattern and ignore not matching ones (scan command)", func() {
			pattern := randString(5)

			values := []string{}
			//checking 200 keys with same pattern
			amountOfValuesWithSamePattern := 200
			for i := 0; i < amountOfValuesWithSamePattern; i++ {
				values = append(values, pattern+string(i))
			}
			//adding few more non-pattern values
			values = append(values, randString(3))
			values = append(values, randString(3))

			for _, key := range values {
				err := Setex(key, "1", 100)
				Expect(err).To(BeNil())
			}
			res, err := Scan(pattern + "*")
			Expect(err).To(BeNil())
			Expect(len(res)).To(Equal(amountOfValuesWithSamePattern))
		})
		It("should get no values by pattern if there is no matching (scan command)", func() {
			pattern := randString(5)

			res, err := Scan(pattern + "*")
			Expect(err).To(BeNil())
			Expect(len(res)).To(Equal(0))
		})
		It("should delete value", func() {
			key := randString(10)
			value := randString(10)
			//save
			err := Setex(key, value, 20)
			Expect(err).To(BeNil())
			//check if it was saved
			newValue := ""
			newValue, err = Get(key)
			Expect(err).To(BeNil())
			Expect(newValue).To(BeEquivalentTo(value))

			//deleting from store
			err = Del(key)
			Expect(err).To(BeNil())

			//checking if it was realy deleted
			newValue, err = Get(key)
			Expect(err).To(BeNil())
			Expect(newValue).NotTo(BeEquivalentTo(value))
		})
		It("should not tolerate negative MaxLength", func() {
			isPositive := (MaxLength >= 0)
			Expect(isPositive).To(BeTrue())
		})
		It("should get values by pattern and ignore not matching ones(KEYS command)", func() {
			pattern := randString(5)

			values := []string{}
			//checking 200 keys with same pattern
			amountOfValuesWithSamePattern := 200
			for i := 0; i < amountOfValuesWithSamePattern; i++ {
				values = append(values, pattern+string(i))
			}
			//adding few more non-pattern values
			values = append(values, randString(3))
			values = append(values, randString(3))

			for _, key := range values {
				err := Setex(key, "1", 100)
				Expect(err).To(BeNil())
			}
			res, err := Keys(pattern + "*")
			Expect(err).To(BeNil())
			Expect(len(res)).To(Equal(amountOfValuesWithSamePattern))
		})
		It("should get no values by pattern if there is no matching(KEYS command)", func() {
			pattern := randString(5)

			res, err := Keys(pattern + "*")
			Expect(err).To(BeNil())
			Expect(len(res)).To(Equal(0))
		})
	})
})
