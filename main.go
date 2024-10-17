package main

import (
    "database/sql"
    "encoding/hex"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "regexp"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/argon2"
)

const (
    dbName      = "./users.db"
    saltLength  = 32
    usernameMax = 15
    passwordMin = 8
)

type User struct {
    ID       int
    Username string
    Password string
}

func GenerateRandomSalt(length int) (string, error) {
    if length <= 0 {
        return "", errors.New("panjang garam harus lebih besar dari 0")
    }

    salt := make([]byte, length)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    return hex.EncodeToString(salt), nil
}

func HashPassword(password, salt string) string {
    hashed := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
    return hex.EncodeToString(hashed)
}

func ValidateUsername(username string) error {
    if len(username) == 0 || len(username) > usernameMax {
        return errors.New("username harus antara 1 dan 15 karakter")
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
        return errors.New("username hanya boleh mengandung huruf dan angka")
    }
    return nil
}

func ValidatePassword(password string) error {
    if len(password) < passwordMin {
        return errors.New("password harus terdiri dari setidaknya 8 karakter")
    }
    if !(regexp.MustCompile(`[0-9]`).MatchString(password) &&
        regexp.MustCompile(`[A-Za-z]`).MatchString(password) &&
        regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)) {
        return errors.New("password harus mengandung huruf, angka, dan karakter khusus")
    }
    return nil
}

func InitializeDatabase() {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    sqlStmt := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        salt TEXT NOT NULL,
        password TEXT NOT NULL
    );`
    if _, err = db.Exec(sqlStmt); err != nil {
        log.Fatalf("%q: %s\n", err, sqlStmt)
    }
}

func Register(username, password string) (string, error) {
    if err := ValidateUsername(username); err != nil {
        return "", err
    }
    if err := ValidatePassword(password); err != nil {
        return "", err
    }

    salt, err := GenerateRandomSalt(saltLength)
    if err != nil {
        return "", fmt.Errorf("kesalahan saat menghasilkan garam: %v", err)
    }

    hashedPassword := HashPassword(password, salt)

    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return "", err
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO users (username, salt, password) VALUES (?, ?, ?)", username, salt, hashedPassword)
    if err != nil {
        if err.Error() == "UNIQUE constraint failed: users.username" {
            return "", errors.New("username ini sudah diambil")
        }
        return "", errors.New("kesalahan saat memasukkan pengguna ke database")
    }

    return "Berhasil: Pengguna terdaftar!", nil
}

func Authenticate(username, password string) (string, error) {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return "", err
    }
    defer db.Close()

    var storedSalt, storedPassword string
    row := db.QueryRow("SELECT salt, password FROM users WHERE username = ?", username)
    if err := row.Scan(&storedSalt, &storedPassword); err != nil {
        return "", errors.New("username tidak ditemukan")
    }

    hashedPassword := HashPassword(password, storedSalt)
    if hashedPassword != storedPassword {
        return "", errors.New("password salah")
    }

    return "Berhasil: Pengguna terautentikasi!", nil
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
    var filename string
    switch r.URL.Path {
    case "/register":
        filename = "register.html"
    case "/authenticate":
        filename = "authenticate.html"
    default:
        http.NotFound(w, r)
        return
    }

    http.ServeFile(w, r, filename)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Metode tidak valid", http.StatusMethodNotAllowed)
        return
    }

    username, password := r.FormValue("username"), r.FormValue("password")

    message, err := Register(username, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Fprintln(w, message)
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Metode tidak valid", http.StatusMethodNotAllowed)
        return
    }

    username, password := r.FormValue("username"), r.FormValue("password")

    message, err := Authenticate(username, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    fmt.Fprintln(w, message)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    InitializeDatabase()

    // Menyajikan halaman HTML dan CSS
    http.HandleFunc("/register", serveHTML)
    http.HandleFunc("/authenticate", serveHTML)
    http.HandleFunc("/register_process", registerHandler) // Ubah nama endpoint
    http.HandleFunc("/authenticate_process", authenticateHandler) // Ubah nama endpoint

    log.Println("Server dimulai di :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
