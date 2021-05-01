package utils

import (
	"crypto/md5"
	"dcomicServer/database"
	"dcomicServer/model"
	"encoding/hex"
	"errors"
	"math/rand"
	"net"
	"time"
)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetNewToken(l int) string {
	token := GetRandomString(l)
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"token": token}).One(&user)
	if err == nil {
		return GetNewToken(l)
	}
	return token
}

func GetPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte("dcomic_salt"))
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetExternalIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, face := range interfaces {
		if face.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if face.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		address, err := face.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range address {
			ip := GetIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("cannot get ip")
}

func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}


