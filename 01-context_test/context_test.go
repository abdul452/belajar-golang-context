package contexttest

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	backround := context.Background()
	fmt.Println(backround)

	todo := context.TODO()
	fmt.Println(todo)
}
