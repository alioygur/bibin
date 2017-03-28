package main

import (
	"errors"
	"log"
	"testing"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/goutil"
)

func TestAliko(t *testing.T) {
	var i domain.Image
	ss, err := goutil.NewSQLStruct(&i)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(ss.Columns())
}

func ali() (i interface{}, err error) {
	err = nil
	defer func() {
		log.Println(err)
	}()
	if 1 == 1 {
		err := errors.New("velii")
		return nil, err
	}
	return nil, err
}
