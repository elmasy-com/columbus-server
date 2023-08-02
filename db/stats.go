package db

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MaxStatsEntry = 100
)

// StatsCountTotal returns the total number of entries in "domain" collection.
func StatsCountTotal() (int64, error) {

	return Domains.CountDocuments(context.TODO(), bson.M{})
}

// StatsCountUpdated returns the total number of entries that updated in "domain" collection.
func StatsCountUpdated() (int64, error) {

	return Domains.CountDocuments(context.TODO(), bson.M{"updated": bson.M{"$exists": true}})
}

// StatsCountValid returns the total number of entries that has at least on valid record in the "records" field in "domain" collection.
func StatsCountValid() (int64, error) {

	return Domains.CountDocuments(context.TODO(), bson.M{"records": bson.M{"$exists": true}})
}

// StatsInsert get the stats and insert a new entry in the "stats" collection.
//
// This function is **very** slow!
func StatsInsert() error {

	s := new(StatSchema)
	var err error

	s.Total, err = StatsCountTotal()
	if err != nil {
		return fmt.Errorf("failed to count total: %w", err)
	}

	s.Updated, err = StatsCountUpdated()
	if err != nil {
		return fmt.Errorf("failed to count updated: %w", err)
	}

	s.Valid, err = StatsCountValid()
	if err != nil {
		return fmt.Errorf("failed to count valid: %w", err)
	}

	s.Scanners, err = ScannerGets()
	if err != nil {
		return fmt.Errorf("failed to get scanners: %w", err)
	}

	s.Date = time.Now().Unix()

	_, err = Stats.InsertOne(context.TODO(), *s)

	return err
}

// StatsInsertWorker insert a new stat entry at the beginning and at a random time in an infinite loop.
//
// This function is designed to run as a goroutine in the background.
// The errors are printed to STDERR.
func StatsInsertWorker() {

	err := StatsInsert()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert new stat entry: %s\n", err)
	}

	for {

		time.Sleep(time.Duration(rand.Int63n(7200)+7200) * time.Second)

		err := StatsInsert()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to insert new stat entry: %s\n", err)
		}
	}
}

// StatsCleanWorker removes entries beyond MaxStatsEntry number.
//
// This function is designed to run as a goroutine in the background.
// The errors are printed to STDERR.
func StatsCleanWorker() {

	// Random sleep
	time.Sleep(time.Duration(rand.Int63n(7200)) * time.Second)

	t := time.Tick(3600 * time.Second)

	for range t {

		n, err := Stats.CountDocuments(context.TODO(), bson.M{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "StatsRemoveOldEntries(): Failed to count total stat entries: %s\n", err)
			continue
		}

		if n <= MaxStatsEntry {
			continue
		}

		i := 0

		cursor, err := Stats.Find(context.TODO(), bson.M{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "StatsRemoveOldEntries(): Failed to find stat entries: %s\n", err)
			continue
		}

		for cursor.Next(context.TODO()) {

			if i <= MaxStatsEntry {
				i++
				continue
			}

			s := new(StatSchema)

			err = cursor.Decode(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "StatsRemoveOldEntries(): Failed to decode: %s\n", err)
				continue
			}

			_, err := Stats.DeleteOne(context.TODO(), *s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "StatsRemoveOldEntries(): Failed to remove stat entry (date: %d, total: %d, updated: %d, valid: %d): %s\n", s.Date, s.Total, s.Updated, s.Valid, err)
			}
		}

		err = cursor.Err()
		if err != nil {
			fmt.Fprintf(os.Stderr, "StatsRemoveOldEntries(): Cursor failed: %s\n", err)
		}

		cursor.Close(context.TODO())

	}
}

// StatsGetNewest returns the newest entry from the "stats" collection.
func StatsGetNewest() (StatSchema, error) {

	s := new(StatSchema)

	err := Stats.FindOne(context.TODO(), bson.M{}, options.FindOne().SetSort(bson.M{"date": -1})).Decode(s)

	return *s, err
}
