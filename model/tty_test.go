package model

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestAddDomainName(t *testing.T) {
	type args struct {
		domainName string
		client     *redis.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddDomainName(tt.args.domainName, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("AddDomainName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterConsulParam_RegisterToConsul(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		r       *RegisterConsulParam
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.RegisterToConsul(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("RegisterConsulParam.RegisterToConsul() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConsulNodeInformation(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConsulNodeInformation(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConsulNodeInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConsulNodeInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUUIDExists(t *testing.T) {
	type args struct {
		name   string
		client *redis.Client
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsUUIDExists(tt.args.name, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUUIDExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUUIDExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
