package queries

import (
	"log"
	"context"
	"time"
	"errors"
	"database/sql"

	"github.com/novadinn/social-network/model"

	"github.com/jmoiron/sqlx"
	_  "github.com/go-sql-driver/mysql"
)

const (
	// TODO: add posts count, followers count, followed count
	createUsersTable = `CREATE TABLE IF NOT EXISTS users (
id VARCHAR(20) PRIMARY KEY, 
email VARCHAR(255) NOT NULL UNIQUE, 
username VARCHAR(50) NOT NULL UNIQUE, 
created_at DATETIME NOT NULL DEFAULT now());`
	createPostsTable = `CREATE TABLE IF NOT EXISTS posts (
id VARCHAR(20) PRIMARY KEY, 
user_id VARCHAR(20) NOT NULL, 
content TEXT NOT NULL, 
created_at TIMESTAMP NOT NULL DEFAULT now(), 
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE);`
	createCommentsTable = `CREATE TABLE IF NOT EXISTS comments (
id VARCHAR(20) PRIMARY KEY, 
user_id VARCHAR(20) NOT NULL, 
post_id VARCHAR(20) NOT NULL, 
content TEXT NOT NULL, 
created_at TIMESTAMP NOT NULL DEFAULT now(), 
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE, 
FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE ON UPDATE CASCADE);`
	createUserFollowsTable = `CREATE TABLE IF NOT EXISTS user_follows (
follower_id VARCHAR(20) NOT NULL, 
followed_id VARCHAR(20) NOT NULL, 
created_at DATETIME NOT NULL DEFAULT now(), 
FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE, 
FOREIGN KEY (followed_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE, 
PRIMARY KEY(follower_id, followed_id))`
	
	createUser = `INSERT IGNORE INTO users (id, email, username, created_at) 
VALUES (?, ?, ?, ?);`
	createPost = `INSERT IGNORE INTO posts (id, user_id, content, created_at) 
VALUES (?, ?, ?, ?);`
	createComment = `INSERT IGNORE INTO comments(id, user_id, post_id, content, created_at) 
VALUES (?, ?, ?, ?, ?)`
	createUserFolow = `INSERT IGNORE INTO user_follows (follower_id, followed_id, created_at) 
VALUES (?, ?, ?);`

	getUserByID = `SELECT * FROM users 
WHERE id = ?`
	getUserByEmail = `SELECT * FROM users 
WHERE email = ?`
	getUserByUsername = `SELECT * FROM users 
WHERE username = ?`
	
	getPosts = `SELECT posts.*, users.username 
FROM posts 
INNER JOIN users ON posts.user_id = users.id 
ORDER BY posts.id DESC`
	getFollowingPosts = `SELECT posts.*, users.username 
FROM posts 
INNER JOIN users ON posts.user_id = users.id 
WHERE posts.user_id IN (?)
ORDER BY posts.id DESC;`
	
	getPostsByUsername = `SELECT posts.*, users.username 
FROM posts 
INNER JOIN users ON posts.user_id = users.id 
WHERE users.username = ? 
ORDER BY posts.id DESC`
	getPostByID = `SELECT posts.*, users.username 
FROM posts 
INNER JOIN users ON posts.user_id = users.id 
WHERE posts.id = ? 
ORDER BY posts.id DESC`

	getCommentsByPostID = `SELECT comments.*, users.username 
FROM comments 
INNER JOIN users ON comments.user_id = users.id 
WHERE comments.post_id = ? 
ORDER BY comments.id DESC`

	getUserFollow = `SELECT * 
FROM user_follows 
WHERE follower_id = ? AND followed_id = ?`
	deleteUserFollow = `DELETE FROM user_follows 
WHERE follower_id = ? AND followed_id = ?`
	getFollowingByUsername = `SELECT users.* 
FROM user_follows 
INNER JOIN users ON user_follows.followed_id = users.id 
WHERE follower_id = (SELECT users.id FROM users 
WHERE username = ?)
ORDER BY user_follows.followed_id;`
	getFollowersByUsername = `SELECT users.* 
FROM user_follows 
INNER JOIN users ON user_follows.follower_id = users.id 
WHERE followed_id = (SELECT users.id FROM users 
WHERE username = ?)
ORDER BY user_follows.follower_id;`
)

type Queries struct {
	db *sql.DB
}

func New() *Queries {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/social_network?parseTime=true")
	if err != nil {
		log.Fatal("Database error: ", err)
		panic(err)
	}

	q := &Queries{db: db}
	if err = q.createTable(createUsersTable); err != nil {
		log.Fatal("Error when creating users table: ", err)
		panic(err)
	}
	if err = q.createTable(createPostsTable); err != nil {
		log.Fatal("Error when creating posts table: ", err)
		panic(err)
	}
	if err = q.createTable(createCommentsTable); err != nil {
		log.Fatal("Error when creating comments table: ", err)
		panic(err)
	}
	if err = q.createTable(createUserFollowsTable); err != nil {
		log.Fatal("Error when creating user follows table: ", err)
		panic(err)
	}
	
	return q
}

func (q *Queries) UserExistsByID(id string) (bool, error) {
	rows, err := q.db.Query(getUserByID, id)
	if err != nil {
		return false, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (q *Queries) UserExistsByEmail(email string) (bool, error) {
	rows, err := q.db.Query(getUserByEmail, email)
	if err != nil {
		return false, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (q *Queries) UserExistsByUsername(username string) (bool, error) {
	rows, err := q.db.Query(getUserByUsername, username)
	if err != nil {
		return false, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (q *Queries) UserFollowExists(followerID, followedID string) (bool, error) {
	rows, err := q.db.Query(getUserFollow, followerID, followedID)
	if err != nil {
		return false, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (q *Queries) DeleteUserFollow(followerID, followedID string) error {
	_, err := q.db.Exec(deleteUserFollow, followerID, followedID)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := q.db.QueryRow(getUserByEmail, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (q *Queries) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := q.db.QueryRow(getUserByUsername, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (q *Queries) CreateUser(id, email, username string, createdAt time.Time) error {
	var enough bool
	err := q.db.QueryRow(createUser, id, email, username, createdAt).Scan(&enough)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		
		return err
	}

	return nil
}

func (q *Queries) CreatePost(id, userID, content string, createdAt time.Time) error {
	var enough bool
	err := q.db.QueryRow(createPost, id, userID, content, createdAt).Scan(&enough)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		
		return err
	}

	return nil
}

func (q *Queries) CreateComment(id, userID, postID, content string, createdAt time.Time) error {
	var enough bool
	err := q.db.QueryRow(createComment, id, userID, postID, content, createdAt).Scan(&enough)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		
		return err
	}

	return nil
}

func (q *Queries) CreateUserFollow(followerID, followedID string, createdAt time.Time) error {
	var enough bool
	err := q.db.QueryRow(createUserFolow, followerID, followedID, createdAt).Scan(&enough)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		
		return err
	}

	return nil
}

func (q *Queries) GetPosts() ([]model.Post, error) {
	rows, err := q.db.Query(getPosts)
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt,
			&post.Username); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return posts, err
    }

	return posts, nil
}

func (q *Queries) GetFollowingPosts(ids []interface{}) ([]model.Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	
	que, args, err := sqlx.In(getFollowingPosts, ids)
	
	rows, err := q.db.Query(que, args...)
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt,
			&post.Username); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return posts, err
    }

	return posts, nil
}

func (q *Queries) GetPostsByUsername(username string) ([]model.Post, error) {
	rows, err := q.db.Query(getPostsByUsername, username)
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt,
			&post.Username); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return posts, err
    }

	return posts, nil
}

func (q *Queries) GetFollowingByUsername(username string) ([]model.User, error) {
	rows, err := q.db.Query(getFollowingByUsername, username)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
        return users, err
    }

	return users, nil
}

func (q *Queries) GetFollowersByUsername(username string) ([]model.User, error) {
	rows, err := q.db.Query(getFollowersByUsername, username)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
        return users, err
    }

	return users, nil
}

func (q *Queries) GetPostByID(id string) (model.Post, error) {
	var post model.Post
	err := q.db.QueryRow(getPostByID, id).Scan(
		&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post, nil
		}
		return model.Post{}, err
	}

	return post, nil
}

func (q *Queries) GetComments(postID string) ([]model.Comment, error) {
	rows, err := q.db.Query(getCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content,
			&comment.CreatedAt, &comment.Username); err != nil {
			return comments, err
		}
		
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
        return comments, err
    }

	return comments, nil
}

func (q *Queries) createTable(query string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  
	defer cancel()

	_, err := q.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
