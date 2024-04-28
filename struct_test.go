package main

import "testing"

func TestWriteConfig(t *testing.T) {
	m := make(map[string]Item)

	m["/hello:GET"] = Item{
		Method: "GET",
		Path:   "/hello",
		Response: Response{
			Status: 200,
			Body:   "Hi guys ivde ellam ok aan",
		}}

	cf := &Config{
		Routes: m,
	}
	WriteConfig(cf)
}
