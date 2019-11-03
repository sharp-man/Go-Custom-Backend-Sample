package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/victorsteven/fullstack/api/controllers"
	"github.com/victorsteven/fullstack/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	user := models.User{
		Username: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Kenny",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndPostTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOnePost() (models.Post, error) {

	user := models.User{
		Username: "Sam",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Post, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Post{}, err
	}
	var users = []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Magu",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var posts = []models.Post{
		models.Post{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		models.Post{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}

func refreshUserPostAndLikeTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}, &models.Like{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Like{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user, post and like tables")
	return nil
}

func seedUsersPostsAndLikes() (models.Post, []models.User, []models.Like, error) {
	// The idea here is: two users can like one post
	var err error
	var users = []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Magu",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	post := models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}
	var likes = []models.Like{
		models.Like{
			UserID: 1,
			PostID: post.ID,
		},
		models.Like{
			UserID: 2,
			PostID: post.ID,
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		err = server.DB.Model(&models.Like{}).Create(&likes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed likes table: %v", err)
		}
	}
	return post, users, likes, nil
}

func refreshUserPostAndCommentTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}, &models.Comment{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user, post and comment tables")
	return nil
}

func seedUsersPostsAndComments() (models.Post, []models.User, []models.Comment, error) {
	// The idea here is: two users can comment one post
	var err error
	var users = []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Magu",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	post := models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}
	var comments = []models.Comment{
		models.Comment{
			Body:   "user 1 made this comment",
			UserID: 1,
			PostID: post.ID,
		},
		models.Comment{
			Body:   "user 2 made this comment",
			UserID: 2,
			PostID: post.ID,
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		err = server.DB.Model(&models.Like{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed comments table: %v", err)
		}
	}
	return post, users, comments, nil
}
