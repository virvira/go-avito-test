package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Segment struct {
    ID    int    `json:"id"`
    Slug  string `json:"slug"`
}

type SegmentSlug struct {
    Slug  string `json:"slug"`
}

func main() {
    //connect to database
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    //create the table if it doesn't exist
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(255), email VARCHAR(64), created_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL, updated_at TIMESTAMP, deleted_at TIMESTAMP);")
    _, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UNIQUE_users_name ON users (email);")
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS segments (id SERIAL PRIMARY KEY, slug VARCHAR(255), created_at TIMESTAMP default now(), updated_at TIMESTAMP, deleted_at TIMESTAMP);")
    _, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UNIQUE_segments_name ON segments (slug);")
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS user_segment (id SERIAL PRIMARY KEY, user_id INTEGER not null REFERENCES users (id), segment_id INTEGER not null REFERENCES segments (id), created_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL, deleted_at TIMESTAMP);")
    _, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS name ON user_segment (user_id, segment_id) WHERE deleted_at is NULL;")

    if err != nil {
        log.Fatal(err)
    }

    //create router
    router := mux.NewRouter()
    router.HandleFunc("/users", getUsers(db)).Methods("GET")
    router.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
    router.HandleFunc("/users", createUser(db)).Methods("POST")
    router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
    router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

    router.HandleFunc("/users/{id}/segments", getUserSegments(db)).Methods("GET")
    router.HandleFunc("/users/{id}/segments", addUserSegments(db)).Methods("POST")
    router.HandleFunc("/users/{id}/segments", deleteUserSegments(db)).Methods("DELETE")

    router.HandleFunc("/segments", getSegments(db)).Methods("GET")
    router.HandleFunc("/segments/{id}", getSegment(db)).Methods("GET")
    router.HandleFunc("/segments", createSegment(db)).Methods("POST")
    router.HandleFunc("/segments/{id}", updateSegment(db)).Methods("PUT")
    router.HandleFunc("/segments/{id}", deleteSegment(db)).Methods("DELETE")

    //start server
    log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, name, email FROM users WHERE deleted_at IS NULL")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        users := []User{}
        for rows.Next() {
            var u User
            if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(users)
    }
}

// get user by id
func getUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u User
        err := db.QueryRow("SELECT * FROM users WHERE id = $1 WHERE deleted_at IS NULL", id).Scan(&u.ID, &u.Name)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(u)
    }
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u User
        json.NewDecoder(r.Body).Decode(&u)

        err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u User
        json.NewDecoder(r.Body).Decode(&u)

        vars := mux.Vars(r)
        id := vars["id"]

        _, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

// delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u User
        err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            _, err := db.Exec("UPDATE users SET deleted_at = now() WHERE id = $1", id)
            if err != nil {
                //todo : fix error handling
                w.WriteHeader(http.StatusNotFound)
                return
            }

            json.NewEncoder(w).Encode("User deleted")
        }
    }
}

func getSegments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, slug FROM segments WHERE deleted_at IS NULL")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        users := []Segment{}
        for rows.Next() {
            var u Segment
            if err := rows.Scan(&u.ID, &u.Slug); err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(users)
    }
}

func getSegment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u Segment
        err := db.QueryRow("SELECT * FROM segments WHERE id = $1 WHERE deleted_at IS NULL", id).Scan(&u.ID, &u.Slug)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(u)
    }
}

func createSegment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u Segment
        json.NewDecoder(r.Body).Decode(&u)

        err := db.QueryRow("INSERT INTO Segment (slug) VALUES ($1, $2) RETURNING id", u.Slug).Scan(&u.ID)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

func updateSegment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u Segment
        json.NewDecoder(r.Body).Decode(&u)

        vars := mux.Vars(r)
        id := vars["id"]

        _, err := db.Exec("UPDATE segments SET slug = $1 WHERE id = $3", u.Slug, id)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

func deleteSegment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u Segment
        err := db.QueryRow("SELECT * FROM segments WHERE id = $1", id).Scan(&u.ID, &u.Slug)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            _, err := db.Exec("DELETE FROM segments WHERE id = $1", id)
            if err != nil {
                //todo : fix error handling
                w.WriteHeader(http.StatusNotFound)
                return
            }

            json.NewEncoder(w).Encode("Segment deleted")
        }
    }
}

func addUserSegments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userId := vars["id"]

        var u []SegmentSlug
        json.NewDecoder(r.Body).Decode(&u)

        for _, n := range u {
            var segmentId int
            row := db.QueryRow("SELECT id FROM segments WHERE slug = $1", n.Slug)

            err := row.Scan(&segmentId)

            if err != nil {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                var existsId int
                err := db.QueryRow("SELECT id FROM user_segment WHERE user_id = $1 AND segment_id = $2 AND deleted_at IS NULL", userId, segmentId).Scan(&existsId)
                if err != nil && existsId == 0 {
                    _, err := db.Exec("INSERT INTO user_segment (user_id, segment_id) VALUES ($1, $2)", userId, segmentId)

                    if err != nil {
                        return
                    }
                }

                json.NewEncoder(w).Encode("Segment added")
            }
        }

        json.NewEncoder(w).Encode(u)
    }
}

func deleteUserSegments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userId := vars["id"]

        var u []SegmentSlug
        json.NewDecoder(r.Body).Decode(&u)

        for _, n := range u {
            var segmentId int
            row := db.QueryRow("SELECT id FROM segments WHERE slug = $1", n.Slug)

            err := row.Scan(&segmentId)

            if err != nil {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                _, err := db.Exec("UPDATE user_segment SET deleted_at = now() WHERE user_id = $1 AND segment_id = $2", userId, segmentId)

                if err != nil {
                    w.WriteHeader(http.StatusNotFound)
                    return
                }

                json.NewEncoder(w).Encode("Segment deleted")
            }
        }

        json.NewEncoder(w).Encode(u)
    }
}

func getUserSegments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userId := vars["id"]

        rows, err := db.Query("SELECT s.id, s.slug FROM segments s JOIN user_segment us ON us.segment_id = s.id WHERE s.deleted_at IS NULL AND us.deleted_at IS NULL AND us.user_id = $1", userId)
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        users := []Segment{}
        for rows.Next() {
            var u Segment
            if err := rows.Scan(&u.ID, &u.Slug); err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(users)
    }
}