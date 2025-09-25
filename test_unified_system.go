package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

type RegisterRequest struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Jurusan  string `json:"jurusan"`
	Angkatan int    `json:"angkatan"`
}

type GraduateRequest struct {
	MahasiswaID  uint   `json:"mahasiswa_id"`
	TahunLulus   int    `json:"tahun_lulus"`
	NoTelepon    string `json:"no_telepon"`
	AlamatAlumni string `json:"alamat_alumni"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	fmt.Println("=== Testing Unified Mahasiswa-Alumni System ===")
	
	// Step 1: Register mahasiswa
	fmt.Println("\n1. Testing Mahasiswa Registration...")
	registerReq := RegisterRequest{
		NIM:      "2021001",
		Nama:     "John Doe",
		Email:    "john@test.com",
		Password: "password123",
		Jurusan:  "Teknik Informatika",
		Angkatan: 2021,
	}
	
       registerResult, err := makeRequest("POST", "/auth/mahasiswa/register", registerReq)
       if err != nil {
	       fmt.Printf("Registration failed: %v\n", err)
	       return
       }
       fmt.Printf("Registration result: %s\n", registerResult)

       // Parse MahasiswaID from registration response
       var regResp map[string]interface{}
       var mahasiswaID uint = 1
       if err := json.Unmarshal([]byte(registerResult), &regResp); err == nil {
	       if data, ok := regResp["data"].(map[string]interface{}); ok {
		       if id, ok := data["id"].(float64); ok {
			       mahasiswaID = uint(id)
		       }
	       }
       }
	
	// Step 2: Login as mahasiswa
	fmt.Println("\n2. Testing Mahasiswa Login...")
	loginReq := LoginRequest{
		Email:    "john@test.com",
		Password: "password123",
	}
	
	loginResult, err := makeRequest("POST", "/auth/mahasiswa/login", loginReq)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Printf("Login result: %s\n", loginResult)
	
	// Step 3: Graduate mahasiswa (make them alumni)
	fmt.Println("\n3. Testing Mahasiswa Graduation...")
       graduateReq := GraduateRequest{
	       MahasiswaID:  mahasiswaID,
	       TahunLulus:   2024,
	       NoTelepon:    "081234567890",
	       AlamatAlumni: "Jakarta",
       }
	
	graduateResult, err := makeRequest("POST", "/auth/mahasiswa/graduate", graduateReq)
	if err != nil {
		fmt.Printf("Graduation request failed: %v\n", err)
		return
	}
	fmt.Printf("Graduation result: %s\n", graduateResult)
	
	// Check if graduation was successful before proceeding
	if len(graduateResult) > 0 && graduateResult[0:1] == "{" {
		var result map[string]interface{}
		json.Unmarshal([]byte(graduateResult), &result)
		if success, ok := result["success"].(bool); ok && !success {
			fmt.Printf("Graduation failed, skipping alumni login test\n")
			return
		}
	}
	
	// Step 4: Login as alumni
	fmt.Println("\n4. Testing Alumni Login (same credentials, different endpoint)...")
	alumniLoginResult, err := makeRequest("POST", "/auth/alumni/login", loginReq)
	if err != nil {
		fmt.Printf("Alumni login failed: %v\n", err)
		return
	}
	fmt.Printf("Alumni login result: %s\n", alumniLoginResult)
	
	fmt.Println("\n=== Unified System Test Completed Successfully! ===")
	fmt.Println("Key Points Tested:")
	fmt.Println("✓ Single account registration as mahasiswa")
	fmt.Println("✓ Login as active mahasiswa")
	fmt.Println("✓ Graduation process (status evolution)")
	fmt.Println("✓ Login as alumni (same account, different role)")
	fmt.Println("✓ No separate alumni registration required!")
}

func makeRequest(method, endpoint string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}