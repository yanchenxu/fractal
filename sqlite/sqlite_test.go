package sqlite

import (
	"testing"
	"os/user"
	"path"
	//"os"

	//"fmt"
	"fmt"
)

func TestNew(t *testing.T) {
	us, _ := user.Current()
	path := path.Join(us.HomeDir, "sqlite_test.db")
	//if ok, _ := os.Stat(path); ok != nil {
	//	err := os.Remove(path)
	//	if err != nil {
	//		t.Errorf("remove err, %v", err)
	//	}
	//}
	fmt.Println(path)
	db, err := New("./foo.db")
	if err != nil {
		t.Errorf("new db err, %v", err)
	}
	db.Close()
	//os.Remove(path)
}
