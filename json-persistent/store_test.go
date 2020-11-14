package store

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

type Human struct {
	Name   string
	Height float64
}

func testFile() *os.File {
	f, err := ioutil.TempFile(".", "store")
	if err != nil {
		panic(err)
	}
	return f
}

func TestOpen(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())

	ioutil.WriteFile(f.Name(), []byte(`{"hello", "world}`), 0644)
	ks, err := Open(f.Name())
	if err != nil {
		t.Error(err)
	}

	if len(ks.Data) != 1 {
		t.Errorf("expected length of data to be %d, but got %d", 1, len(ks.Data))
	}

	if world, ok := ks.Data["hello"]; !ok || string(world) != "world" {
		t.Errorf("expected value of key \"hello\" to be %s, but got %s", "world", world)
	}
}

func TestGeneral(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())

	ks := new(Store)

	err := ks.Set("hello", "world")
	if err != nil {
		t.Errorf("failed when testing Set: %s", err.Error())
	}

	if err = Save(ks, f.Name()); err != nil {
		t.Errorf("failed when saving Set: %s", err.Error())
	}

	ks2, _ := Open(f.Name())

	var a, b string

	ks.Get("hello", &a)
	ks2.Get("hello", &b)

	if a != b {
		t.Errorf("inconsistent Get value, %v should be the same as %v", a, b)
	}

	ks.Set("human:1", Human{"Dante", 5.4})
	Save(ks, "test2.json.gz")
	Save(ks, "test2.json")

	defer os.Remove("test2.json.gz")
	defer os.Remove("test2.json")

	ks2, err = Open("test2.json")
	if err != nil {
		t.Error(err.Error())
	}

	var human Human
	ks.Get("human:1", &human)
	if human.Height != 5.4 {
		t.Errorf("expected value of '%v' but got '%v'", 5.4, human.Height)
	}

	ks2, err = Open("test2.json.gz")
	if err != nil {
		t.Error(err.Error())
	}

	ks2.Get("human:1", &human)
	if human.Height != 5.4 {
		t.Errorf("expected value of '%v' but got '%v'", 5.4, human.Height)
	}

}

func TestRegex(t *testing.T) {
	f := testFile()

	defer os.Remove(f.Name())

	ks := new(Store)
	ks.Set("hello:1", "world1")
	ks.Set("hello:2", "world2")
	ks.Set("hello:3", "world3")
	ks.Set("world:1", "hello1")

	if len(ks.GetAll(regexp.MustCompile(`hello`))) != len(ks.Keys())-1 {
		t.Errorf("problem on getting all the values")
	}

	if len(ks.GetAll(nil)) != len(ks.Keys()) {
		t.Errorf("problem on getting all the values")
	}
}
