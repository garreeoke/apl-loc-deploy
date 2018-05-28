package work

import (
	"testing"
	propeller "applariat.io/propeller/types"
	"fmt"
)

func TestClusterMgr(t *testing.T) {

	err := ClusterMgr(&propeller.RestData{})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

}