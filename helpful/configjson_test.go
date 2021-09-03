package helpful

import "testing"

const defaultData = `
{
	"A": {
			"B":"test"
		},
	"arr": [
			{"val":"someval1"},{"val":"someval2"}
		]
}
`

func testJsonConfig(t *testing.T, cfg *jsonConfig) {

	res, err := cfg.Child("A").GetString("B")
	if err != nil {
		t.Fatalf("cant get value from hardcoded json config: %v", err)
	}
	if res != "test" {
		t.Errorf("invalid value %v instead of %v", res, "test")
	}

	res, err = cfg.Child("arr").Child("0").GetString("val")
	if err != nil {
		t.Errorf("cant get value from hardcoded json config: %v", err)
	}
	if res != "someval1" {
		t.Errorf("invalid value %v instead of %v", res, "test")
	}

	children, err := cfg.GetArray("arr")
	if err != nil {
		t.Errorf("cant get children in hardcoded json config")
	} else {
		testChildren(t, children)
	}
}

func TestJsonFromBytes(t *testing.T) {
	data := defaultData
	cfg, err := newJsonCfgFromBytes([]byte(data))
	if err != nil {
		t.Fatal("cant init json config from hardcode")
	}
	testJsonConfig(t, cfg)
}

func TestJsonFromBytesWithRefresh(t *testing.T) {
	//var data = new(string)
	data := defaultData
	cfg, f, err := newJsonCfgWithRefresh(func() ([]byte, error) {
		return []byte(data), nil
	})
	if err != nil {
		t.Fatalf("cant get value from hardcoded json config: %v", err)
	}
	testJsonConfig(t, cfg)

	data = `
{
	"A": "some"
}
`
	err = f()
	if err != nil {
		t.Fatalf("cant update hardcoded json config")
	}
	res, err := cfg.GetString("A")
	if err != nil {
		t.Fatalf("cant get val for A in updated hardcoded json config")
	}
	if res != "some" {
		t.Fatalf("incorrect value from updated config")
	}

}

func testChildren(t *testing.T, children []Config) {
	if len(children) != 2 {
		t.Errorf("incorrect array size in hardcoded json config")
		return
	}
	res, err := children[1].GetString("val")
	if err != nil {
		t.Errorf("cant get value from hardcoded json config: %v", err)
	}
	if res != "someval2" {
		t.Errorf("invalid value %v instead of %v", res, "test")
	}

	_, err = children[1].AsMap()
	if err != nil {
		t.Errorf("cant convert object to map")
	}
}
