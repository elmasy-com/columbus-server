package db

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ApiKeyLength = 48
)

// IsNameTaken check whether the given user name is already taken.
//
// Returns fault.ErrNameEmpty if name is empty.
func IsNameTaken(name string) (bool, error) {

	if name == "" {
		return false, fault.ErrNameEmpty
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
//
// If user name is already taken returns fault.ErrNameTaken error.
//
// If name is empty returns fault.ErrNameEmpty.
func UserCreate(name string, admin bool) (user.User, error) {

	if name == "" {
		return user.User{}, fault.ErrNameEmpty
	}

	if taken, err := IsNameTaken(name); err != nil {
		return user.User{}, fmt.Errorf("failed to check user name: %w", err)
	} else if taken {
		return user.User{}, fault.ErrNameTaken
	}

	key, err := genAPIKey()
	if err != nil {
		return user.User{}, fmt.Errorf("failed to generate key: %w", err)
	}

	user := user.User{
		Key:   key,
		Name:  name,
		Admin: admin,
	}

	_, err = Users.InsertOne(context.TODO(), bson.D{{Key: "key", Value: user.Key}, {Key: "name", Value: user.Name}, {Key: "admin", Value: user.Admin}})

	return user, err
}

// UserGetKey returns the user from the db based on API key.
//
// If user not found, returns fault.ErrUserNotFound error.
//
// If key is empty, returns fault.ErrMissingAPIKey error.
func UserGetKey(key string) (user.User, error) {

	if key == "" {
		return user.User{}, fault.ErrMissingAPIKey
	}

	var u user.User

	r := Users.FindOne(context.TODO(), bson.M{"key": key})
	if r.Err() != nil {
		if errors.Is(r.Err(), mongo.ErrNoDocuments) {
			return u, fault.ErrUserNotFound
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
// If user not found, returns fault.ErrUserNotFound error.
//
// If name is empty, returns fault.ErrNameEmpty error.
func UserGetName(name string) (user.User, error) {

	if name == "" {
		return user.User{}, fault.ErrNameEmpty
	}

	var u user.User

	r := Users.FindOne(context.TODO(), bson.M{"name": name})
	if r.Err() != nil {
		if errors.Is(r.Err(), mongo.ErrNoDocuments) {
			return u, fault.ErrUserNotFound
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
// If user not found, returns fault.ErrNameEmpty error.
//
// If key is empty, returns fault.ErrMissingAPIKey error.
//
// If DeleteOne() deletes 0 user, returns fault.ErrUserNotDeleted.
func UserDelete(key, name string) error {

	if key == "" {
		return fault.ErrMissingAPIKey
	}
	if name == "" {
		return fault.ErrNameEmpty
	}

	r, err := Users.DeleteOne(context.TODO(), bson.M{"key": key, "name": name})
	if err != nil {
		return err
	}

	if r.DeletedCount == 0 {
		return fault.ErrUserNotDeleted
	}

	return nil
}

// UserChangeKey update the API key for the given user, and change the key in user.
//
// If user nil, returns fault.ErrUserNil.
//
// If u.Key/u.Name is empty, returns fault.ErrMissingAPIKey/fault.ErrNameEmpty.
//
// If document not modified returns fault.ErrNotModified.
func UserChangeKey(u *user.User) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}
	if u.Name == "" {
		return fault.ErrNameEmpty
	}

	newKey, err := genAPIKey()
	if err != nil {
		return fmt.Errorf("failed to generate API key: %w", err)
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": u.Key, "name": u.Name}, bson.M{"$set": bson.M{"key": newKey}})
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	if r.ModifiedCount == 0 {
		return fault.ErrNotModified
	}

	u.Key = newKey

	return nil
}

// UserChangeName update the name for the given user, and change the name in user.
//
// If user nil, returns fault.ErrUserNil.
//
// If u.Key/u.Name is empty, returns fault.ErrMissingAPIKey/fault.ErrNameEmpty.
//
// If name is taken, returns fault.ErrNameTaken.
func UserChangeName(u *user.User, name string) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}
	if u.Name == "" {
		return fault.ErrNameEmpty
	}

	if taken, err := IsNameTaken(name); err != nil {
		return fmt.Errorf("failed to check name: %w", err)
	} else if taken {
		return fault.ErrNameTaken
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": u.Key, "name": u.Name}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		return err
	}
	if r.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	u.Name = name

	return nil
}

// UserChangeAdmin update the admin field for the given user, and change the admin value in user.
// If user nil, returns fault.ErrUserNil.
//
// If key/name is empty, returns fault.ErrMissingAPIKey/fault.ErrNameEmpty.
//
// If username is taken, returns fault.ErrNameTaken.
func UserChangeAdmin(u *user.User, newValue bool) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}
	if u.Name == "" {
		return fault.ErrNameEmpty
	}

	r, err := Users.UpdateOne(context.TODO(), bson.M{"key": u.Key, "name": u.Name}, bson.M{"$set": bson.M{"admin": newValue}})
	if err != nil {
		return err
	}
	if r.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	u.Admin = newValue

	return nil
}

// UserList returns a list of every users.
func UserList() ([]user.User, error) {

	cursor, err := Users.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users := make([]user.User, 0)

	for cursor.Next(context.TODO()) {

		u := user.User{}

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
