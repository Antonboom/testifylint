package http

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ExpectMirroredRequest(t *testing.T, client interface{}, clientset interface{}, mirrorPods []BackendRef, path string) {
	for i, mirrorPod := range mirrorPods {
		if mirrorPod.Name == "" {
			t.Fatalf("Mirrored BackendRef[%d].Name wasn't provided in the testcase, this test should only check http request mirror.", i)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(mirrorPods))

	for _, mirrorPod := range mirrorPods {
		go func(mirrorPod BackendRef) {
			defer wg.Done()

			require.Eventually(t, func() bool { // want "go-require: require must only be used in the goroutine running the test function"
				mirrorLogRegexp := regexp.MustCompile(fmt.Sprintf("Echoing back request made to \\%s to client", path))
				// ...

				for _, log := range [][]byte{} {
					if mirrorLogRegexp.MatchString(string(log)) {
						return true
					}
				}
				return false
			}, 60*time.Second, time.Millisecond*100, fmt.Sprintf(`Couldn't find mirrored request in "%s/%s" logs`))
		}(mirrorPod)
	}

	wg.Wait()

	t.Log("Found mirrored request log in all desired backends")
}

type BackendRef struct {
	Name      string
	Namespace string
}
