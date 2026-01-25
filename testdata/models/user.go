package 

type User struct {
    id string `json:"id"`
    email string `json:"email"`
    created_at time.Time `json:"created_at"`
}
