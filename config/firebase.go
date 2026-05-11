package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseApp adalah instance Firebase global
var FirebaseApp *firebase.App

// FirebaseAuth adalah instance Firebase Auth global
var FirebaseAuth *auth.Client

// InitFirebase menginisialisasi Firebase Admin SDK
// menggunakan file service account dari Firebase project Dannisa Sweet
func InitFirebase() {
	// Ambil path file credential dari .env
	// Default: firebase-service-account.json di root project
	credentialFile := os.Getenv("FIREBASE_CREDENTIAL_FILE")
	if credentialFile == "" {
		credentialFile = "firebase-service-account.json"
	}

	// Cek apakah file credential ada
	if _, err := os.Stat(credentialFile); os.IsNotExist(err) {
		log.Fatalf(
			"[Firebase] File credential tidak ditemukan: %s\n"+
				"Pastikan kamu sudah:\n"+
				"1. Buka Firebase Console → Project Dannisa Sweet\n"+
				"2. Project Settings → Service Accounts\n"+
				"3. Generate new private key → download JSON\n"+
				"4. Rename & taruh di root project sebagai: %s",
			credentialFile, credentialFile,
		)
	}

	// Inisialisasi Firebase App dengan credential Dannisa Sweet
	opt := option.WithCredentialsFile(credentialFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("[Firebase] Gagal inisialisasi Firebase App: %v", err)
	}
	FirebaseApp = app

	// Inisialisasi Firebase Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("[Firebase] Gagal inisialisasi Firebase Auth: %v", err)
	}
	FirebaseAuth = authClient

	log.Println("[Firebase] Firebase Dannisa Sweet berhasil diinisialisasi ✓")
}

// VerifyFirebaseToken memverifikasi Firebase ID token dari Flutter
// Dipanggil di middleware untuk setiap request yang butuh autentikasi
func VerifyFirebaseToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()
	token, err := FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}