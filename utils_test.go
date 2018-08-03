package redistore

import (
	"fmt"
	"os"
	"smarpshare/util"

	"."

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var store, storeDeprecating *redistore.RediStore

var updateStore = func(newStore *redistore.RediStore) {
	store = newStore
}

const (
	defaultRedisHost = "127.0.0.1"
	defaultRedisPort = "6379"
)

func setup() string {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = defaultRedisHost
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = defaultRedisPort
	}

	return fmt.Sprintf("%s:%s", addr, port)
}

var _ = Describe("deprecated store cookie options", func() {

	Context("redis save/load/delete func", func() {
		FIt("should say value is too big", func() {
			// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
			setup()
			key := util.RandString(10)
			value := util.RandString(MaxLength + 1)
			age := 20
			err := Setex(key, value, age)
			Expect(err).NotTo(BeNil())
		})
		// It("should save value for 1 sec", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
		// 	Connect2SessionStorage()
		// 	key := util.RandString(10)
		// 	value := util.RandString(10)
		// 	//saving value for 1 sec
		// 	age := 1
		// 	err := SaveByKey(key, value, age)
		// 	Expect(err).To(BeNil())

		// 	//check if it was saved
		// 	newValue := ""
		// 	newValue, err = LoadByKey(key)
		// 	Expect(err).To(BeNil())
		// 	Expect(newValue).To(BeEquivalentTo(value))

		// 	//sleep for 2 sec
		// 	time.Sleep(2 * time.Second)
		// 	newValue, err = LoadByKey(key)
		// 	Expect(err).To(BeNil())
		// 	Expect(newValue).NotTo(BeEquivalentTo(value))
		// })
		// It("should load value", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis

		// 	err := Connect2SessionStorage()
		// 	Expect(err).To(BeNil())
		// 	key := util.RandString(10)
		// 	value := util.RandString(10)
		// 	err = SaveByKey(key, value, 20)
		// 	Expect(err).To(BeNil())

		// 	newValue := ""
		// 	newValue, err = LoadByKey(key)
		// 	Expect(err).To(BeNil())
		// 	Expect(newValue).To(BeEquivalentTo(value))
		// })
		// It("should load values by pattern and ignore not matching ones", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
		// 	Connect2SessionStorage()
		// 	pattern := util.RandString(5)

		// 	values := []string{}
		// 	//checking 200 keys with same pattern
		// 	amountOfValuesWithSamePattern := 200
		// 	for i := 0; i < amountOfValuesWithSamePattern; i++ {
		// 		values = append(values, pattern+string(i))
		// 	}
		// 	//adding few more non-pattern values
		// 	values = append(values, util.RandString(3))
		// 	values = append(values, util.RandString(3))

		// 	for _, key := range values {
		// 		err := SaveByKey(key, "1", 100)
		// 		Expect(err).To(BeNil())
		// 	}
		// 	res, err := GetByPattern(pattern + "*")
		// 	Expect(err).To(BeNil())
		// 	Expect(len(res)).To(Equal(amountOfValuesWithSamePattern))
		// })
		// It("should load 0 values by pattern if there is no matching", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
		// 	Connect2SessionStorage()
		// 	pattern := util.RandString(5)

		// 	res, err := GetByPattern(pattern + "*")
		// 	Expect(err).To(BeNil())
		// 	Expect(len(res)).To(Equal(0))
		// })

		// It("should delete value", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect"
		// 	err := Connect2SessionStorage()
		// 	Expect(err).To(BeNil())
		// 	key := util.RandString(10)
		// 	value := util.RandString(10)
		// 	//save
		// 	err = SaveByKey(key, value, 20)
		// 	Expect(err).To(BeNil())
		// 	//check if it was saved
		// 	newValue := ""
		// 	newValue, err = LoadByKey(key)
		// 	Expect(err).To(BeNil())
		// 	Expect(newValue).To(BeEquivalentTo(value))

		// 	//deleting from store
		// 	err = DeleteByKey(key)
		// 	Expect(err).To(BeNil())

		// 	//checking if it was realy deleted
		// 	newValue, err = LoadByKey(key)
		// 	Expect(err).To(BeNil())
		// 	Expect(newValue).NotTo(BeEquivalentTo(value))
		// })

		// It("MaxLength should not be negative", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect"
		// 	isPositive := (MaxLength >= 0)
		// 	Expect(isPositive).To(BeTrue())
		// })
		// It("should load values by pattern and ignore not matching ones", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
		// 	Connect2SessionStorage()
		// 	pattern := util.RandString(5)

		// 	values := []string{}
		// 	//checking 200 keys with same pattern
		// 	amountOfValuesWithSamePattern := 200
		// 	for i := 0; i < amountOfValuesWithSamePattern; i++ {
		// 		values = append(values, pattern+string(i))
		// 	}
		// 	//adding few more non-pattern values
		// 	values = append(values, util.RandString(3))
		// 	values = append(values, util.RandString(3))

		// 	for _, key := range values {
		// 		err := SaveByKey(key, "1", 100)
		// 		Expect(err).To(BeNil())
		// 	}
		// 	res, err := GetByPatternKeys(pattern + "*")
		// 	Expect(err).To(BeNil())
		// 	Expect(len(res)).To(Equal(amountOfValuesWithSamePattern))
		// })
		// It("should load 0 values by pattern if there is no matching", func() {
		// 	// Connect2SessionStorage will most probably return error "missing address" or "cannot connect". Test it with already set up redis
		// 	Connect2SessionStorage()
		// 	pattern := util.RandString(5)

		// 	res, err := GetByPatternKeys(pattern + "*")
		// 	Expect(err).To(BeNil())
		// 	Expect(len(res)).To(Equal(0))
		// })
	})
})
