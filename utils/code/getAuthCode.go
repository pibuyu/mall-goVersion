package code

import (
	"fmt"
	"math/rand"
	"time"
)

func GetAuthCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
}
