package db

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elmasy-com/columbus-server/config"
	"github.com/elmasy-com/columbus-server/fault"
	"github.com/elmasy-com/elnet/dns"
	"github.com/elmasy-com/elnet/valid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	RecordsUpdaterDomainChan         chan string
	internalRecordsUpdaterDomainChan chan string
	totalUpdated                     atomic.Uint64
	startTime                        time.Time
)

// Update type t records for d.
// This function updates the DB.
func recordsUpdateRecord(d string, t uint16) error {

	if !valid.Domain(d) {
		return fault.ErrInvalidDomain
	}

	d = dns.Clean(d)

	p := dns.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return fault.ErrGetPartsFailed
	}

	var (
		r   []string
		err error
	)

	switch t {

	case dns.TypeA:
		// A
		r, err = dns.QueryARetryStr(d)

	case dns.TypeAAAA:
		// AAAA
		r, err = dns.QueryAAAARetryStr(d)

	case dns.TypeCAA:
		// CAA
		r, err = dns.QueryCAARetryStr(d)

	case dns.TypeCNAME:
		// CNAME
		r, err = dns.QueryCNAMERetry(d)

	case dns.TypeDNAME:
		//DNAME
		v, err2 := dns.QueryDNAMERetry(d)
		err = err2
		if v != "" {
			r = append(r, v)
		}

	case dns.TypeMX:
		// MX
		r, err = dns.QueryMXRetryStr(d)

	case dns.TypeNS:
		// NS
		r, err = dns.QueryNSRetry(d)

	case dns.TypeSOA:
		// SOA
		v, err2 := dns.QuerySOARetryStr(d)
		err = err2
		if v != "" {
			r = append(r, v)
		}

	case dns.TypeSRV:
		// SRV
		r, err = dns.QuerySRVRetryStr(d)

	case dns.TypeTXT:
		// TXT
		r, err = dns.QueryTXTRetry(d)

	default:
		return fmt.Errorf("invalid type: %d", t)
	}

	if err != nil {
		return err
	}

	for i := range r {

		// "records" field should contain only one element with "type" t and "value" v.
		// Try to update first!
		// If MatchedCount is 0, the record with "type" t and "value" r[i] is new and the new record will be appended to the array.
		// If MatchedCount is 1, only one record is exist with "type" t and "value" v and the time for the element is updated.
		// If MatchedCount is > 1, duplicate record found, ERROR!
		filter := bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "sub", Value: p.Sub}, {Key: "records.type", Value: t}, {Key: "records.value", Value: r[i]}}

		up := bson.D{{Key: "$set", Value: bson.D{{Key: "records.$.time", Value: time.Now().Unix()}}}}

		result, err := Domains.UpdateOne(context.TODO(), filter, up)
		if err != nil {
			return err
		}
		if result.MatchedCount > 1 {
			return fmt.Errorf("duplicate record found: %s", r[i])
		}
		if result.MatchedCount == 1 {
			continue
		}

		// Append new record to "records"
		filter = bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "sub", Value: p.Sub}}

		up = bson.D{{Key: "$addToSet", Value: bson.D{{Key: "records", Value: RecordSchema{Type: t, Value: r[i], Time: time.Now().Unix()}}}}}

		_, err = Domains.UpdateOne(context.TODO(), filter, up)
		if err != nil {
			return err
		}
	}

	return nil
}

// RecordsUpdate updates the records field for domain d.
// This function updates the records in the database.
// If the same record found, updates the "time" field in element.
// If new record found, append it to the "records" field.
//
// If ignore is true, common DNS errors are ignored.
//
// If domain d is invalid, returns fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns fault.ErrGetPartsFailed.
func RecordsUpdate(d string, ignore bool) error {

	err := recordsUpdateRecord(d, dns.TypeA)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update A: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update A: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeAAAA)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update AAAA: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update AAAA: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeCAA)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update CAA: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update CAA: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeCNAME)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update CNAME: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update CNAME: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeDNAME)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update DNAME: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update DNAME: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeMX)
	if err != nil {
		if !ignore {
			return fmt.Errorf("failed to update MX: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update MX: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeNS)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update NS: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update NS: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeSOA)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update SOA: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update SOA: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeSRV)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update SRV: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update SRV: %w", err)
		}
	}

	err = recordsUpdateRecord(d, dns.TypeTXT)
	if err != nil {

		if !ignore {
			return fmt.Errorf("failed to update TXT: %w", err)
		}

		if !errors.Is(err, dns.ErrName) && !errors.Is(err, dns.ErrServerFailure) &&
			!os.IsTimeout(err) && !errors.Is(err, dns.ErrRefused) {

			return fmt.Errorf("failed to update TXT: %w", err)
		}
	}

	return nil
}

// recordsUpdaterRoutine reads from DomainChan and internalDomainChan
// and updates the FQDN coming from the channel.
func recordsUpdaterRoutine(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		if totalUpdated.Load()%10000 == 0 {
			if totalUpdated.Load() != 0 {
				fmt.Printf("recordsUpdaterRoutine(): Updated %d domain records in %s\n", totalUpdated.Load(), time.Since(startTime))
			}
		}

		var d string

		select {
		case dom := <-RecordsUpdaterDomainChan:
			d = dom
			fmt.Printf("recordsUpdaterRoutine updating %s...\n", dom)

		case dom := <-internalRecordsUpdaterDomainChan:
			d = dom
		}

		if dns.HasSub(d) {

			totalUpdated.Add(1)

			// d is a FQDN
			err := RecordsUpdate(d, true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update DNS records for %s: %s\n", d, err)
			}

		} else {

			// If domain sent instead of FQDN, get every subdomain and updates it
			ds, err := LookupFull(d, -1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update DNS records for %s: %s\n", d, err)
				continue
			}

			totalUpdated.Add(uint64(len(ds)))

			for i := range ds {

				err := RecordsUpdate(ds[i], true)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update DNS records for %s: %s\n", ds[i], err)
				}
			}
		}
	}
}

// RandomDomainUpdater is a function created to run as goroutine in the background.
// Select random entries (FQDNs) and send it to internalDomainChan to update the records.
func RandomDomainUpdater(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		cursor, err := Domains.Aggregate(context.TODO(), bson.A{bson.M{"$sample": bson.M{"size": 1000}}}, options.Aggregate().SetBatchSize(1000))
		if err != nil {
			fmt.Fprintf(os.Stderr, "RandomDomainUpdater() failed to find toplist: %s\n", err)
			// Wait before the next try
			time.Sleep(600 * time.Second)
			continue
		}

		for cursor.Next(context.TODO()) {

			d := new(DomainSchema)

			err = cursor.Decode(d)
			if err != nil {
				fmt.Fprintf(os.Stderr, "RandomDomainUpdater() failed to find: %s\n", err)
				break
			}

			// Skip verified domain for now.
			// TODO: Remove later
			if len(d.Records) > 0 {
				continue
			}

			internalRecordsUpdaterDomainChan <- d.String()

		}

		if err = cursor.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "RandomDomainUpdater() cursor failed: %s\n", err)
		}

		cursor.Close(context.TODO())
	}
}

// TopListUpdater is a function created to run as goroutine in the background.
// Updates the domains and it subdomains in topList collection.
// This function uses concurrent goroutines and print only/ignores any error.
func TopListUpdater(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		time.Sleep(time.Duration(rand.Intn(49) * int(time.Hour)))

		start := time.Now()

		cursor, err := TopList.Find(context.TODO(), bson.M{}, options.Find().SetBatchSize(10).SetSort(bson.M{"count": -1}))
		if err != nil {
			fmt.Fprintf(os.Stderr, "TopListUpdater() failed to find toplist: %s\n", err)
			continue
		}

		for cursor.Next(context.TODO()) {

			d := new(TopListSchema)

			err = cursor.Decode(d)
			if err != nil {
				fmt.Fprintf(os.Stderr, "TopListUpdater() failed to find: %s\n", err)
				break
			}

			internalRecordsUpdaterDomainChan <- d.Domain

		}

		if err = cursor.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "TopListUpdater() cursor failed: %s\n", err)
		}

		cursor.Close(context.TODO())
		fmt.Printf("TopListUpdater(): Finished updating topList in %s\n", time.Since(start))
	}
}

func RecordsUpdater() {

	RecordsUpdaterDomainChan = make(chan string, config.DomainBuffer)
	internalRecordsUpdaterDomainChan = make(chan string, config.DomainBuffer)
	startTime = time.Now()

	wg := new(sync.WaitGroup)

	for i := 0; i < config.DomainWorker; i++ {
		wg.Add(1)
		go recordsUpdaterRoutine(wg)
	}

	wg.Add(1)
	go RandomDomainUpdater(wg)

	wg.Add(1)
	go TopListUpdater(wg)

	wg.Wait()
}
