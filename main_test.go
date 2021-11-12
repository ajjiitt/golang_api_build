package main

import (
	controller "golang_api_build/controllers"
	"reflect"
	"testing"
)

func TestHashPassword(t *testing.T) {
	var test = []struct {
		userPassword     string
		providedPassword string
		checkOP          bool
	}{
		{userPassword: "tija4321", providedPassword: "$2a$14$r1EKEy7qZue6VadatT9oUutoeGEEkR/bDrtzAbEx.hn/NvGD/r0nm", checkOP: true},{userPassword: "tija42", providedPassword: "$2a$14$r1EKEy7qZue6VadatT9oUutoeGEEkR/bDrtzAbEx.hn/NvGD/r0nm", checkOP: true},
	}

	for _, value := range test {
		if check, msg := controller.VerifyPassword(value.userPassword, value.providedPassword); check == value.checkOP && len(msg) > 0 {
			t.Error("Test failed")
		}
	}
}

func TestHashingPassword(t *testing.T) {
	var test = []struct {
		password string
	}{
		{password: "tija4321"},
	}

	for _,value := range test {
		if v := controller.HashPassword(value.password); reflect.TypeOf(v) != reflect.TypeOf("temp") {
			t.Error("Test failed: func not returning string value")
		}
	}
}
