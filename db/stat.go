package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// GetStat returns the total number of domains, the total number of updated domains and the total number of domains with valid record.
func GetStat() (total int64, updated int64, valid int64, scanners []ScannerSchema, err error) {

	total, err = Domains.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to count total: %w", err)
	}

	updated, err = Domains.CountDocuments(context.TODO(), bson.M{"updated": bson.M{"$exists": true}})
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to count updated: %w", err)
	}

	valid, err = Domains.CountDocuments(context.TODO(), bson.M{"records": bson.M{"$exists": true}})
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to count valid: %w", err)
	}

	scanners, err = ScannerGetIndexes()
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to get scanners: %w", err)
	}

	return total, updated, valid, scanners, err
}
