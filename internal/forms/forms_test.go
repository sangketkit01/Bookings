package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestValidForm(t *testing.T){
	r , err := http.NewRequest("POST","/some-url",nil)
	if err != nil{
		t.Error("Can't not request to the url")
	}

	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("Form should be valid")
	}
}

func TestRequiredForm(t *testing.T){
	r,err := http.NewRequest("POST","/some-url",nil)
	if err != nil{
		t.Error("Can't not request to the url")
	}

	form := New(r.PostForm)
	form.Required("a","b","c")
	if form.Valid() {
		t.Error("Form should not be valid")
	}

	postData := url.Values{}
	postData.Add("a","a")
	postData.Add("b","b")
	postData.Add("c","c")
	r.PostForm = postData

	form = New(r.PostForm)
	if !form.Valid(){
		t.Error("Form should be valid")
	}
}

func TestHasForm(t *testing.T) {
	r, err := http.NewRequest("POST", "/some-url", nil)
	if err != nil {
		t.Error("Can't not request to the url")
	}

	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("The value is empty , Has function should return false")
	}

	postData := url.Values{}
	postData.Add("b", "hello")
	r.PostForm = postData
	form = New(r.PostForm)

	if !form.Has("b") {
		t.Error("The field has a value, Has function should return true", form.Errors)
	}
}

func TestMinLengthForm(t *testing.T){
	r,err := http.NewRequest("POST","/some-url",nil)
	if err != nil {
		t.Error("Can't not request to the url")
	}

	form := New(r.PostForm)
	if form.MinLength("x",10){
		t.Error("Form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == ""{
		t.Error("should have an error, but did not get one")
	}


	postData := url.Values{}
	postData.Add("somefield1","hello")
	form = New(postData)

	if !form.MinLength("somefield1",5){
		t.Error(fmt.Sprintf("Form must have %d minimum length characters but only has %d characters",5,len(form.Get("somefield1"))))
	}

	isError = form.Errors.Get("somefield1")
	if isError != ""{
		t.Error("shouldn't have an error, but did get one")
	}
}

func TestEmailForm(t *testing.T){
	r,err := http.NewRequest("POST","/some-url",nil)
	if err != nil {
		t.Error("Can't not request to the url")
	}

	postData := url.Values{}
	postData.Add("test","test")
	r.PostForm = postData
	form := New(r.PostForm)

	form.IsEmail("test")

	if form.Valid(){
		t.Error("The value that the dedicated field has contains is not an email",form.Errors.Get("test"))
	}

	postData.Add("email","thiraphat.sa@kkumail.com")
	form = New(postData)
	form.IsEmail("email")

	if !form.Valid(){
		t.Error("The value that the dedicated field has contains is an email",form.Errors.Get("email"))
	}
}
