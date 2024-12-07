package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/satyendra001/mdm-oauth/utils"
)

type UserStore struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

func NewUserStore(conn *pgxpool.Pool, ctx context.Context) *UserStore {
	return &UserStore{
		conn: conn,
		ctx:  ctx,
	}
}

func (store *UserStore) GetAllUserInfo() {

	log.Println("Extracting all User Info")
	rows, err := store.conn.Query(store.ctx, "SELECT id, username, first_name, last_name FROM kraft_canada.auth_user")
	if err != nil {
		log.Println("Error in Querying DB ==> ", err.Error())
	}

	dbUser := new(utils.DBUser)

	log.Println("Rows data ==> ", rows)
	for rows.Next() {
		err := rows.Scan(
			&dbUser.Id,
			&dbUser.Username,
			&dbUser.FirstName,
			&dbUser.LastName,
		)

		if err != nil {
			log.Println("Error while scanning the rows ==> ", err.Error())
			return
		}

		log.Println("User Data ==> ", dbUser)
	}

	rows.Close()
}

func (store *UserStore) GetUserByEmail(email string) (user *utils.DBUser, err error) {
	conn := store.conn
	checkUserQuery := "SELECT id, username, email FROM kraft_canada.auth_user WHERE email=$1"

	rows, err := conn.Query(store.ctx, checkUserQuery, email)

	if err != nil {
		log.Println("User doesn't exists with email", email, "Error => ", err.Error())
		return nil, err
	}

	defer rows.Close()

	dbUser := new(utils.DBUser)

	for rows.Next() {
		err := rows.Scan(
			&dbUser.Id,
			&dbUser.Username,
			&dbUser.Email,
		)

		if err != nil {
			log.Println("Error scanning user with email", email, " Error ==>", err.Error())
			return nil, err
		}
	}

	// Handle the case when no user found in DB. In that case DB retruns id as 0
	if dbUser.Id == 0 {
		log.Println("No user found in DB. DB returned ID as 0")
		return nil, fmt.Errorf("user not found in db with email %s", email)
	}

	log.Println("Found a user in DB. Returning the user with email => ", email)
	return dbUser, nil
}

func (store *UserStore) CreateUser(OauthUser *goth.User) (int, error) {
	log.Println("Creating new User...")
	var userId int

	userCreationQuery := `INSERT INTO kraft_canada.auth_user
						(username, email, first_name, last_name, password, is_superuser, is_staff, is_active, date_joined)
						VALUES ($1, $2, $3, $4, $5, False, False, True, $6) RETURNING id`

	values := []interface{}{
		OauthUser.Name,
		OauthUser.Email,
		OauthUser.FirstName,
		OauthUser.LastName,
		"1qwerty",
		time.Now(),
	}

	row := store.conn.QueryRow(store.ctx, userCreationQuery, values...)

	err := row.Scan(&userId)

	if err != nil {
		log.Println("Error while creating new User ==>", err.Error())
		return 0, err
	}

	log.Println("Successfully created new User...")
	return userId, nil

}

func (store *UserStore) GetToken(userId int) string {
	var token string

	tokenQuery := "SELECT key FROM kraft_canada.authtoken_token WHERE user_id = $1"

	tokenRow := store.conn.QueryRow(store.ctx, tokenQuery, userId)

	tokenRow.Scan(&token)

	log.Println("Received token from the DB ==> ", token)

	return token
}

func (store *UserStore) CreateToken(userId int) (string, error) {
	log.Println("Started token creation for userId ", userId)

	token, err := utils.GenerateToken()

	if err != nil {
		log.Println("Error while creating a random token Err==>", err.Error())
		return "", err
	}
	tokenInsertQuery := `INSERT INTO kraft_canada.authtoken_token (key, user_id, created) VALUES ($1, $2, $3)`
	tokenQueryValues := []interface{}{token, userId, time.Now()}

	_, err = store.conn.Exec(store.ctx, tokenInsertQuery, tokenQueryValues...)

	if err != nil {
		log.Println("Error while Inserting token in DB ==> ", err.Error())
		return "", err
	}

	log.Println("Successfully created token for the user with ID ", userId)
	return token, nil
}
