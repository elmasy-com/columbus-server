package db

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mrand "math/rand"
	"time"

	sdk "github.com/elmasy-com/columbus-sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ApiKeyLength = 48
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNameEmpty  = errors.New("name is empty")
	ErrUserKeyEmpty   = errors.New("key is empty")
	ErrUserNameTaken  = errors.New("name is taken")
	ErrUserNil        = errors.New("user is nil")
	ErrUserNotDeleted = errors.New("not deleted")
)

// IsNameTaken check whether the given user name is already taken.
//
// Returns ErrUserNameEmpty if name is empty.
func IsNameTaken(name string) (bool, error) {

	if name == "" {
		return false, ErrUserNameEmpty
	}

	n, err := Users.CountDocuments(context.TODO(), bson.M{"name": name})
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

// genAPIKey generates API key  and check is the generated key is exist in the database.
// If crypto/rand fail, use math/rand.
func genAPIKey() (string, error) {

	// The chances of collision is **very very** low
	for {

		key := make([]byte, ApiKeyLength)

		n, err := crand.Read(key)
		if err == nil && n == ApiKeyLength {
			goto check
		}

		// crypto/rand failed, fallback to math/rand
		mrand.Seed(time.Now().UnixMilli())

		mrand.Read(key)

	check:
		keyStr := hex.EncodeToString(key)
		c, err := Users.CountDocuments(context.TODO(), bson.M{"key": keyStr})
		if err != nil {
			return "", err
		}
		if c == 0 {
			return keyStr, nil
		}
	}
}

// UserCreate creates a new user in the users collection and returns it.
// The API key automatically generated.
// If admin is true, the new user will be an admin.
// If user name is already taken returns ErrUserNameTaken error.
func UserCreate(name string, admin bool) (sdk.User, error) {

	if name == "" {
		return sdk.User{}, ErrUserNameEmpty
	}

	if taken, err := IsNameTaken(name); err != nil {
		return sdk.User{}, fmt.Errorf("failed to check user name: %w", err)
	} else if taken {
		return sdk.User{}, ErrUserNameTaken
	}

	key, err := genAPIKey()
	if err != nil {
		return sdk.User{}, fmt.Errorf("failed to generate key: %w", err)
	}

	user := sdk.User{
		Key:   key,
		Name:  name,
		Admin: admin,
	}

	_, err = Users.InsertOne(context.TODO(), bson.D{{Key: "key", Value: user.Key}, {Key: "name", Value: user.Name}, {Key: "admin", Value: user.Admin}})

	return user, err
}

// UserGetKey returns the user from the db based on API key.
//
// If user not found, returns ErrUserNotFound error.
//
// If key is empty, returns ErrUserKeyEmpty error.
func UserGetKey(key string) (sdk.User, error) {

	if key == "" {
		return sdk.User{}, ErrUserKeyEmpty
	}

	var u sdk.User

	r := Users.FindOne(context.TODO(), bson.M{"key": key})
	if r.Err() != nil {
		if errors.Is(r.Err(), mongo.ErrNoDocuments) {
			return u, ErrUserNotFound
		}
		return u, fmt.Errorf("failed to find key: %w", r.Err())
	}

	err := r.Decode(&u)
	if err != nil {
		return u, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

// UserGetName returns the user from the db based on name.
//
// If user not found, returns ErrUserNotFound error.
//
// If name is empty, returns ErrUserNameEmpty error.
func UserGetName(name string) (sdk.User, error) {

	if name == "" {
		return sdk.User{}, ErrUserNameEmpty
	}

	var u sdk.User

	r := Users.FindOne(context.TODO(), bson.M{"name": name})
	if r.Err() != nil {
		if errors.Is(r.Err(), mongo.ErrNoDocuments) {
			return u, ErrUserNotFound
		}
		return u, fmt.Errorf("failed to find name: %w", r.Err())
	}

	err := r.Decode(&u)
	if err != nil {
		return u, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

// UserDelete delete user based on key+name.
//
// If user not found, returns ErrUserNotFound error.
//
// If key is empty, returns ErrUserKeyEmpty error.
//
// If DeleteOne() deletes 0 user, returns ErrUserNotDeleted.
func UserDelete(key, name string) error {

	if key == "" {
		return ErrUserKeyEmpty
	}
	if name == "" {
		return ErrUserNameEmpty
	}

	r, err := Users.DeleteOne(context.TODO(), bson.M{"key": key, "name": name})
	if err != nil {
		return err
	}

	if r.DeletedCount == 0 {
		return ErrUserNotDeleted
	}

	return nil
}

// UserChangeKey update the API key for the given user, and change the key in user.
//
// If user nil, returns ErrUserNil.
//
// If key/name is empty, returns ErrUserKeyEmpty/ErrUserNameEmpty.
func UserChangeKey(user *sdk.User) error {

	if user == nil {
		return ErrUserNil
	}
	if user.Key == "" {
		return ErrUserKeyEmpty
	}
	if user.Name == "" {
		return ErrUserNameEmpty
	}

	newKey, err := genAPIKey()
	if err != nil {
		return fmt.Errorf("failed to generate API key: %w", err)
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": user.Key, "name": user.Name}, bson.M{"$set": bson.M{"key": newKey}})
	if err != nil {
		return err
	}
	if r.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	user.Key = newKey

	return nil
}

// UserChangeName update the name for the given user, and change the name in user.
//
// If user nil, returns ErrUserNil.
//
// If key/name is empty, returns ErrUserKeyEmpty/ErrUserNameEmpty.
//
// If username is taken, returns ErrUserNameTaken.
func UserChangeName(user *sdk.User, newName string) error {

	if user == nil {
		return ErrUserNil
	}
	if user.Key == "" {
		return ErrUserKeyEmpty
	}
	if user.Name == "" {
		return ErrUserNameEmpty
	}

	if taken, err := IsNameTaken(newName); err != nil {
		return fmt.Errorf("failed to check name: %w", err)
	} else if taken {
		return ErrUserNameTaken
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": user.Key, "name": user.Name}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return err
	}
	if r.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	user.Name = newName

	return nil
}

// UserChangeAdmin update the admin field for the given user, and change the admin value in user.
// If key is empty, returns ErrUserKeyEmpty.
// If name is empty, returns ErrUserNameEmpty.
func UserChangeAdmin(user *sdk.User, newValue bool) error {

	if user == nil {
		return ErrUserNil
	}
	if user.Key == "" {
		return ErrUserKeyEmpty
	}
	if user.Name == "" {
		return ErrUserNameEmpty
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": user.Key, "name": user.Name}, bson.M{"$set": bson.M{"admin": newValue}})
	if err != nil {
		return err
	}
	if r.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	user.Admin = newValue

	return nil
}

// UserList returns a list of every users.
func UserList() ([]sdk.User, error) {

	cursor, err := Users.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users := make([]sdk.User, 0)

	for cursor.Next(context.TODO()) {

		u := sdk.User{}

		err := cursor.Decode(&u)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}

		users = append(users, u)
	}

	if cursor.Err() != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}
