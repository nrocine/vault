package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestZookeeperBackend(t *testing.T) {
	addr := os.Getenv("ZOOKEEPER_ADDR")
	if addr == "" {
		t.SkipNow()
	}

	client, _, err := zk.Connect([]string{addr}, time.Second)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("/vault-%d", time.Now().Unix())
	acl := zk.WorldACL(zk.PermAll)
	_, err = client.Create(randPath, []byte("hi"), int32(0), acl)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	defer func() {
		client.Delete(randPath, -1)
	}()

	b, err := NewBackend("zookeeper", map[string]string{
		"address": addr + "," + addr,
		"path":    randPath,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
