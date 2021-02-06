package controller

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("before test portfolioControler_test.go")
	code := m.Run()
	fmt.Println("after test portfolioControler_test.go")
	os.Exit(code)
}

func TestFetchMailSendSelect(t *testing.T) {
}

func TestFetchMailBatchStatus(t *testing.T) {
}

func TestFetchAllMembers(t *testing.T) {
}

func TestFetchAllWorker(t *testing.T) {
}

func TestFetchLoginInfo(t *testing.T) {
}

func TestFetchCreditInfoRegist(t *testing.T) {
}

func TestFetchSendMailRegist(t *testing.T) {
}

func TestFetchProfileInfo(t *testing.T) {
}

func TestFetchRegistAccount(t *testing.T) {
}

func TestFetchRegistAccountMail(t *testing.T) {
}

func TestFetchMailAdrInfo(t *testing.T) {
}

func TestFetchSignUpAccountMail(t *testing.T) {
}
