package main

import (
	"social-backend/config"
	"social-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()


	r := gin.Default()
	routes.UserRoutes(r)
	routes.PostRoutes(r)
	routes.ImageRoutes(r)
	routes.RegisterRoutes(r)
	routes.InteractionRoutes(r)
	
	// routes.PostRoutes(r)


	r.Run("0.0.0.0:8080")
}

// package main

// import (
//     "fmt"
//     "golang.org/x/crypto/bcrypt"
// )

// func main() {
//     hashedPassword := "$2a$10$CukLi5RCBmFT6j0dBPbZze3NQUxY9KUwCxQAA/zPtcITImeyimzri"
//     password := "password123"

//     err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//     if err != nil {
//         fmt.Println("Password does NOT match")
//     } else {
//         fmt.Println("Password matches")
//     }
// }

// package main

// import (
// 	"fmt"
// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	password := "password123" // Change this to the correct password
// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	fmt.Println(string(hashedPassword))
// }
