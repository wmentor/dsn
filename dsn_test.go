package dsn

import (
	"testing"
)

func TestDSN(t *testing.T) {

	tP := func(str string, wait map[string]string, err error) {
		dsn, e := New(str)
		if e != nil {
			if e != err {
				t.Fatal("Invalid error for: " + str)
			}
		}

		if len(dsn) != len(wait) {
			t.Fatal("Invalid map size for: " + str)
		}

		for k, v := range wait {
			if rv, h := dsn[k]; h {
				if rv != v {
					t.Fatal("Invalid map " + k + "=" + rv + " for: " + str)
				}
			} else {
				t.Fatal("Key not found for: " + str)
			}
		}
	}

	tP("", map[string]string{}, nil)
	tP(" ", map[string]string{}, nil)
	tP("=", map[string]string{"": ""}, nil)
	tP("1=", map[string]string{"1": ""}, nil)
	tP("1=2 3=4 empty= addr=127.0.0.1", map[string]string{"1": "2", "3": "4", "empty": "", "addr": "127.0.0.1"}, nil)
	tP("test=true", map[string]string{"test": "true"}, nil)
	tP("test=t=rue", nil, ErrParse)
	tP("test==true", nil, ErrParse)
	tP("t", nil, ErrParse)
	tP("tes", nil, ErrParse)
	tP("1= test=true 1=12", nil, ErrKeyAlreadyExists)
	tP("1= test=true 1=12 ", nil, ErrKeyAlreadyExists)
	tP(`val=1 example=1+2\=3`, map[string]string{"example": "1+2=3", "val": "1"}, nil)
	tP(`\=v\=al=1 example=\=1+2\=3`, map[string]string{"example": "=1+2=3", "=v=al": "1"}, nil)
	tP(`message=Hello,\ World!`, map[string]string{"message": "Hello, World!"}, nil)
	tP(`message=Hello,\sWorld! answer=Hi`, map[string]string{"message": "Hello, World!", "answer": "Hi"}, nil)
	tP(`message=Hello,\sWorld! answer=Hi\`, nil, ErrParse)
	tP(`message=Hello,\sWorld! answer`, nil, ErrParse)
	tP(`message=Hello,\sWorld! message= `, nil, ErrKeyAlreadyExists)
	tP(`message=Hello,\sWorld! ans\wer=Hi`, nil, ErrParse)

	tN := func(str string) DSN {
		dsn, e := New(str)
		if e != nil {
			t.Fatal("Failed: " + str)
		}
		return dsn
	}

	dsn := tN("server=127.0.0.1 port=80 sslmode=true single=1 keepalive=false domain=test.ru timeout=1.5")

	if dsn.GetString("server", "!!!") != "127.0.0.1" {
		t.Fatal("ivalid return value")
	}

	if dsn.GetString("server1", "!!!") != "!!!" {
		t.Fatal("ivalid return value")
	}

	dsn.SetString("server", "192.168.1.1")

	if dsn.GetString("server", "!!!") != "192.168.1.1" {
		t.Fatal("ivalid return value")
	}

	if dsn.GetInt("port", 443) != 80 {
		t.Fatal("Invalid GetInt")
	}

	if dsn.GetInt("timeout", 2) != 2 {
		t.Fatal("Invalid timeout")
	}

	dsn.SetInt("port", 8080)

	if dsn.GetInt("port", 443) != 8080 {
		t.Fatal("Invalid SetInt")
	}

	if !dsn.GetBool("sslmode", false) {
		t.Fatal("GetBool failed")
	}

	if dsn.GetBool("server", false) {
		t.Fatal("GetBool failed")
	}

	if dsn.GetBool("keepalive", true) {
		t.Fatal("GetBool failed")
	}

	if !dsn.GetBool("single", false) {
		t.Fatal("GetBool failed")
	}

	dsn.SetBool("single", false)

	if dsn.GetBool("single", true) {
		t.Fatal("SetBool failed")
	}

	if dsn.GetFloat("timeout", 2) != 1.5 {
		t.Fatal("GetFloat failed")
	}

	if dsn.GetFloat("server", 12) != 12 {
		t.Fatal("GetFloat failed")
	}

	dsn.SetFloat("timeout", 5)
	if dsn.GetFloat("timeout", 2) != 5 {
		t.Fatal("GetFloat failed")
	}

	if dsn.GetInt64("port", 443) != 8080 {
		t.Fatal("Invalid GetInt")
	}

	if dsn.GetInt64("timeout", 2) != 2 {
		t.Fatal("Invalid timeout")
	}

	dsn.SetInt64("port", 80)

	if dsn.GetInt64("port", 443) != 80 {
		t.Fatal("Invalid SetInt")
	}

	tSR := func(src string) {
		dsn = tN(src)
		str := dsn.String()
		ndsn, err := New(str)
		if err != nil {
			t.Fatal("New failed for: " + dsn.String())
		}

		if len(ndsn) != len(dsn) {
			t.Fatal("Invalid dsn size for: " + dsn.String())
		}

		for k, v := range dsn {
			if rv, h := ndsn[k]; h {
				if rv != v {
					t.Fatal("Invalid key for: " + dsn.String())
				}
			} else {
				t.Fatal("Key not found for: " + dsn.String())
			}
		}
	}

	tSR(``)
	tSR(`data=`)
	tSR(`msg=1+2\=3`)
	tSR(`\=1=2 \=msg=hello\sworld`)
}
