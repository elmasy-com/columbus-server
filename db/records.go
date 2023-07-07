package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-server/config"
	"github.com/elmasy-com/elnet/dns"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Update the time in type t and value v record in d domain.
func recordsUpdateRecordTime(d *DomainSchema, t uint16, v string) error {

	if d == nil {
		return fmt.Errorf("DomainSchema is nil")
	}

	filter := bson.D{{Key: "domain", Value: d.Domain}, {Key: "tld", Value: d.TLD}, {Key: "sub", Value: d.Sub}, {Key: "records.type", Value: t}, {Key: "records.value", Value: v}}

	up := bson.D{{Key: "$set", Value: bson.D{{Key: "records.$.time", Value: time.Now().Unix()}}}}

	result, err := Domains.UpdateOne(context.TODO(), filter, up)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	return err
}

// Append a new record with type t and value v to domain d.
func recordsAppendRecord(d *DomainSchema, t uint16, v string) error {

	if d == nil {
		return fmt.Errorf("DomainSchema is nil")
	}

	filter := bson.D{{Key: "domain", Value: d.Domain}, {Key: "tld", Value: d.TLD}, {Key: "sub", Value: d.Sub}}

	up := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "records", Value: RecordSchema{Type: t, Value: v, Time: time.Now().Unix()}}}}}

	result, err := Domains.UpdateOne(context.TODO(), filter, up)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	return err
}

// Update type t records for d.
// This function updates the DB.
func recordsUpdateRecord(d *DomainSchema, t uint16) error {

	if d == nil {
		return fmt.Errorf("DomainSchema is nil")
	}

	var (
		r   []string
		err error
	)

	switch t {
	case dns.TypeA:
		r, err = dns.QueryARetryStr(d.String())
	case dns.TypeAAAA:
		r, err = dns.QueryAAAARetryStr(d.String())
	case dns.TypeTXT:
		r, err = dns.QueryTXTRetry(d.String())
	case dns.TypeCNAME:
		r, err = dns.QueryCNAMERetry(d.String())
	case dns.TypeMX:
		r, err = dns.QueryMXRetryStr(d.String())
	case dns.TypeNS:
		r, err = dns.QueryNSRetry(d.String())
	case dns.TypeCAA:
		r, err = dns.QueryCAARetryStr(d.String())
	case dns.TypeSRV:
		r, err = dns.QuerySRVRetryStr(d.String())
	default:
		return fmt.Errorf("invalid type: %d", t)
	}

	if err != nil {
		return err
	}

outerLoop:
	for i := range r {

		for ii := range d.Records {

			if d.Records[ii].Type == t && d.Records[ii].Value == r[i] {

				err = recordsUpdateRecordTime(d, t, r[i])
				if err != nil {
					return fmt.Errorf("failed to update time for %d %s record: %w", t, r[i], err)
				}
				continue outerLoop
			}
		}

		err = recordsAppendRecord(d, t, r[i])
		if err != nil {
			return fmt.Errorf("failed to append %d %s: %w", t, r[i], err)
		}
	}

	return nil
}

// RecordsUpdate updates the records field for domain d.
// If the same record found, updates the time in element.
// If new recotrd found, append it to te Records field.
// This function updates the records in the database.
func RecordsUpdate(d *DomainSchema) error {

	if d == nil {
		return fmt.Errorf("DomainSchema is nil")
	}

	err := recordsUpdateRecord(d, dns.TypeA)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update A records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeAAAA)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update AAAA records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeTXT)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update TXT records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeCNAME)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update CNAME records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeMX)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update MX records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeNS)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update NS records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeCAA)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update CAA records: %w", err)
	}

	err = recordsUpdateRecord(d, dns.TypeSRV)
	if err != nil && !errors.Is(err, dns.ErrName) &&
		!errors.Is(err, dns.ErrServerFailure) &&
		strings.Contains(err.Error(), "i/o timeout") {

		return fmt.Errorf("failed to update SRV records: %w", err)
	}

	return nil
}

func recordsUpdaterRoutine(doms <-chan *DomainSchema, wg *sync.WaitGroup) {

	defer wg.Done()

	for dom := range doms {

		err := recordsUpdateRecord(dom, dns.TypeA)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update A records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeAAAA)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update AAAA records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeTXT)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update TXT records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeCNAME)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update CNAME records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeMX)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update MX records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeNS)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update MX records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeCAA)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update CAA records for %s: %s\n", dom, err)
		}

		err = recordsUpdateRecord(dom, dns.TypeSRV)
		if err != nil && !errors.Is(err, dns.ErrName) {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to update SRV records for %s: %s\n", dom, err)
		}
	}

}

// RecordUpdater is a function created to run as goroutine in the background.
// Select random domains and update the DNS records.
// This function uses concurrent goroutines and ignores any error.
func RecordsUpdater() {

	for {

		var (
			start = time.Now()
			wg    = new(sync.WaitGroup)
			doms  = make(chan *DomainSchema, 100)
		)

		cursor, err := Domains.Aggregate(context.TODO(), bson.A{bson.M{"$sample": bson.M{"size": 1000}}}, options.Aggregate().SetBatchSize(100))
		if err != nil {
			fmt.Fprintf(os.Stderr, "RecordUpdater() failed to find: %s\n", err)
			// Wait before the next try
			time.Sleep(600 * time.Second)
			continue
		}

		for i := 0; i < config.DNSWorker; i++ {
			wg.Add(1)
			go recordsUpdaterRoutine(doms, wg)
		}

	domainLoop:
		for cursor.Next(context.TODO()) {

			d := new(DomainSchema)

			err = cursor.Decode(d)
			if err != nil {
				fmt.Fprintf(os.Stderr, "RecordUpdater() failed to find: %s\n", err)
				break domainLoop
			}

			doms <- d
		}

		if err = cursor.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "RecordUpdater() cursor failed: %s\n", err)
		}

		cursor.Close(context.TODO())
		close(doms)
		wg.Wait()
		fmt.Printf("Finished updating 1000 domain records in %s\n", time.Since(start))

	}
}
