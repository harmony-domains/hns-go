package onens

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestARegistration(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	// domain := unregisteredDomain()
	domain := "testjw.country"
	name, err := NewName(tclient, domain)
	require.Nil(t, err, "Failed to create name")

	// Register stage 1
	opts, err := generateTxOpts(registrant, registrantKey, "0")
	require.Nil(t, err, "Failed to generate transaction options")
	tx, secret, err := name.RegisterStageOne(registrant, tconfig.duration, opts)
	require.Nil(t, err, "Failed to send stage one transaction")
	// // Wait until mined
	waitForTransaction(tx.Hash())

	// Wait until ready to submit stage 2
	interval, err := name.RegistrationInterval()
	require.Nil(t, err, "Failed to obtain registration interval")
	time.Sleep(interval)
	// Sleep for 10 more seconds to ensure we are over the interval
	time.Sleep(10 * time.Second)

	// Register stage 2
	opts, err = generateTxOpts(registrant, registrantKey, "1200 Ether")
	require.Nil(t, err, "Failed to generate transaction options")
	tx, err = name.RegisterStageTwo(registrant, tconfig.duration, secret, opts)
	require.Nil(t, err, "Failed to send stage two transaction")
	// // Wait until mined
	waitForTransaction(tx.Hash())

	// Confirm registered
	registeredRegistrant, err := name.Registrant()
	require.Nil(t, err, "Failed to obtain registrant")
	assert.Equal(t, tconfig.NameWrapper, registeredRegistrant, "failed to register name")
}

func JWTestARecordCreateImport(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	domain := "test.country"
	// name := "a.test.country."
	// ip := 128.0.0.1
	// record imported from ens-deployer testing
	// TODO format our own records using https://pkg.go.dev/github.com/miekg/dns
	// blog https://miek.nl/2014/september/21/idn-and-private-rr-in-go-dns/
	// code to migrate https://www.npmjs.com/package/dns-js
	aRec := "0161047465737407636f756e747279000001000100000e10000480000001"
	aRecBytes, err := hex.DecodeString(aRec)
	// Get Resolver
	dnsresolver, err := NewDNSResolver(tconfig.client, domain)
	require.Nil(t, err, "Failed to create resolver for domain: %s", domain)

	// Clear Records
	opts, err := generateTxOpts(registrant, registrantKey, "0")
	tx, err := dnsresolver.ClearRecords(opts)
	require.Nil(t, err, "Failed to clear records")
	waitForTransaction(tx.Hash())

	// Check Records have been cleared
	existingRec, err := dnsresolver.Record("a.test.country.", 1)
	require.Nil(t, err, "Failed to get records")
	assert.Equal(t, existingRec, []byte{}, "Failed to clear existing records")

	// Set Records
	opts, err = generateTxOpts(registrant, registrantKey, "0")
	tx, err = dnsresolver.SetRecords(opts, aRecBytes)
	require.Nil(t, err, "Failed to update A record")
	waitForTransaction(tx.Hash())

	// Check A Record Created
	existingRec, err = dnsresolver.Record("a.test.country.", 1)
	require.Nil(t, err, "Failed to get records")
	assert.Equal(t, existingRec, []byte{}, "Failed to clear existing records")

}

func TestARecordCreateOnly(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	opts, err := generateTxOpts(registrant, registrantKey, "0")
	// opts, err := generateTxOpts(registrant, registrantKey, "1200 Ether")
	// fmt.Printf("registrant: %+v\n", registrant)
	// fmt.Printf("registrantKey: %+v\n", registrantKey)
	// fmt.Printf("opts: %+v\n", opts)
	// domain := unregisteredDomain()
	domain := "test.country"
	name := "a.test.country."

	aRec := new(dns.A)
	aRec.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
	aRec.A = []byte{128, 0, 0, 1}
	// aRecIn := name + "3600 IN A 1.2.3.4"
	// aRec, err := dns.NewRR(aRecIn)
	fmt.Printf("arec.Hdr: %+v\n", aRec.Hdr)
	require.Nil(t, err, "Failed to create a record  for domain")
	u := new(dns.RFC3597)
	u.ToRFC3597(aRec)
	fmt.Printf("u.ToRFC3597S: %v\n", u.String())
	fmt.Printf("u.ToRFC3597H: %v\n", u.Header())
	fmt.Printf("u.ToRFC3597H: %v\n", u)
	aRecStringHex := hex.EncodeToString([]byte(aRec.String()))
	aRecBytes, err := hex.DecodeString(aRecStringHex)
	// aRecBytes := []byte(aRec.String())
	// fmt.Printf("arec                                                   : %+v\n", aRec)
	// fmt.Printf("arec.Header()                                          : %+v\n", aRec.Header())
	// fmt.Printf("arec.String()                                          : %+v\n", aRec.String())
	// fmt.Printf("[]byte(arec.String())                                  : %+v\n", []byte(aRec.String()))
	// fmt.Printf("hex.EncodeToString([]byte(arec.String()))              : %+v\n", hex.EncodeToString([]byte(aRec.String())))
	// fmt.Printf("0x + hex.EncodeToString([]byte(arec.String()))         : %+v\n", ("0x" + hex.EncodeToString([]byte(aRec.String()))))
	// fmt.Printf("[]byte(0x + hex.EncodeToString([]byte(arec.String()))) : %+v\n", []byte("0x"+hex.EncodeToString([]byte(aRec.String()))))

	sampleInRec := []byte(hex.EncodeToString([]byte("atestcountry")))
	sampleInRec = append(sampleInRec, 128, 0, 0, 1)
	sampleRecord, err := hex.DecodeString(string(sampleInRec))
	// sampleRecord, err := hex.DecodeString("0161047465737407636f756e747279000001000100000e10000480000001")
	fmt.Printf("Sample InRec : %+v\n", sampleInRec)
	fmt.Printf("Sample InRecString : %+v\n", string(sampleInRec))
	fmt.Printf("Sample Record: %+v\n", string(sampleRecord))
	fmt.Printf("aRec     : %+v\n", aRec)
	fmt.Printf("aRecStringHex: %+v\n", aRecStringHex)
	fmt.Printf("aRecBytes: %+v\n", aRecBytes)

	// aRec := "0161047465737407636f756e747279000001000100000e10000480000001"
	// aRecBytes, err := hex.DecodeString(aRec)
	dnsresolver, err := NewDNSResolver(tconfig.client, domain)
	require.Nil(t, err, "Failed to create resolver for domain: %s", domain)
	// opts, err := generateTxOpts(registrant, registrantKey, "0")
	require.Nil(t, err, "Failed to generate transaction options")
	fmt.Println("Existing records before Clear")
	existingRec, err := dnsresolver.Record("a.test.country.", 1)
	// fmt.Printf("err: %+v\n", err)
	fmt.Printf("existingRec: %+v\n", existingRec)
	require.Nil(t, err, "Failed to get existing records")

	fmt.Println("About to clear records")
	// opts, err = generateTxOpts(registrant, registrantKey, "0")
	fmt.Printf("opts: %+v\n", opts)
	tx, err := dnsresolver.ClearRecords(opts)
	// fmt.Printf("err: %+v\n", err)
	// fmt.Printf("tx: %+v\n", tx)
	require.Nil(t, err, "Failed to clear records")
	waitForTransaction(tx.Hash())

	fmt.Println("Existing records after Clear")
	existingRec, err = dnsresolver.Record("a.test.country.", 1)
	// fmt.Printf("err: %+v\n", err)
	fmt.Printf("existingRec: %+v\n", existingRec)
	require.Nil(t, err, "Failed to get existing records")
	fmt.Println("About to Set Records")
	// opts, err = generateTxOpts(registrant, registrantKey, "0.1 Ether")
	// fmt.Printf("opts: %+v\n", opts)
	opts, err = generateTxOpts(registrant, registrantKey, "0")
	tx, err = dnsresolver.SetRecords(opts, aRecBytes)
	fmt.Println("Set DNS Records")
	// fmt.Printf("err: %+v\n", err)
	// fmt.Printf("tx: %+v\n", tx)
	require.Nil(t, err, "Failed to update A record")
	waitForTransaction(tx.Hash())

	// Retrieve A Record
}

func TestARecordReadOnly(t *testing.T) {
	domain := "test.country."
	name := "a.test.country."
	fmt.Println("In TestARecordReadOnly")
	fmt.Printf("name: %s domain %s\n", name, domain)
	ethDomain := strings.TrimSuffix(domain, ".")
	fmt.Printf("ethDomain: %+v\n", ethDomain)
	dnsResolver, err := NewDNSResolver(tconfig.client, ethDomain)
	require.Nil(t, err, "Failed to retrieve dnsResolver")
	arec, err := dnsResolver.Record(name, dns.TypeA)
	require.Nil(t, err, "Failed to retrieve arec")
	fmt.Printf("arec: %+v\n", arec)
}

func TestARecord(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	opts, err := generateTxOpts(registrant, registrantKey, "0")
	domain := unregisteredDomain()
	name, err := NewName(tclient, domain)
	require.Nil(t, err, "Failed to create name")

	// Register stage 1
	require.Nil(t, err, "Failed to generate transaction options")
	tx, secret, err := name.RegisterStageOne(registrant, tconfig.duration, opts)
	require.Nil(t, err, "Failed to send stage one transaction")
	// // Wait until mined
	waitForTransaction(tx.Hash())

	// Wait until ready to submit stage 2
	interval, err := name.RegistrationInterval()
	require.Nil(t, err, "Failed to obtain registration interval")
	time.Sleep(interval)
	// Sleep for 10 more seconds to ensure we are over the interval
	time.Sleep(10 * time.Second)

	// Register stage 2
	opts, err = generateTxOpts(registrant, registrantKey, "1200 Ether")
	require.Nil(t, err, "Failed to generate transaction options")
	tx, err = name.RegisterStageTwo(registrant, tconfig.duration, secret, opts)
	require.Nil(t, err, "Failed to send stage two transaction")
	// // Wait until mined
	waitForTransaction(tx.Hash())

	// Confirm registered
	registeredRegistrant, err := name.Registrant()
	require.Nil(t, err, "Failed to obtain registrant")
	assert.Equal(t, tconfig.NameWrapper, registeredRegistrant, "failed to register name")

	// Write A Record
	// Sample Mapping
	// a.test.country. 3600 IN A 1.2.3.4
	/*
	   name: a.test.country.
	   type: A
	   class: IN
	   ttl: 3600
	   address: 1.2.3.4
	*/
	// arec, err := dns.NewRR("a.testjw.country. 3600 IN A 1.2.3.4")
	// r := new(dns.A)
	// r.Hdr = dns.RR_Header{Name: "a.test.country.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
	// r.A = net.IPv4(1, 2, 3, 4)
	// arec := r
	// arec, err := dns.NewRR("a.test.country 1.2.3.4")
	anameinfo := "a." + domain + ". 3600 IN A 1.2.3.4"
	arec, err := dns.NewRR(anameinfo)
	require.Nil(t, err, "Failed to create a record  for domain")
	fmt.Printf("arec                                                   : %+v\n", arec)
	fmt.Printf("arec.Header()                                          : %+v\n", arec.Header())
	fmt.Printf("arec.String()                                          : %+v\n", arec.String())
	fmt.Printf("[]byte(arec.String())                                  : %+v\n", []byte(arec.String()))
	fmt.Printf("hex.EncodeToString([]byte(arec.String()))              : %+v\n", hex.EncodeToString([]byte(arec.String())))
	fmt.Printf("0x + hex.EncodeToString([]byte(arec.String()))         : %+v\n", ("0x" + hex.EncodeToString([]byte(arec.String()))))
	fmt.Printf("[]byte(0x + hex.EncodeToString([]byte(arec.String()))) : %+v\n", []byte("0x"+hex.EncodeToString([]byte(arec.String()))))
	// fmt.Printf("arec.unpack(): %+v\n", arec.unpack())

	dnsresolver, err := NewDNSResolver(tconfig.client, "testjw.country")
	require.Nil(t, err, "Failed to create resolver for domain: %s", domain)
	opts, err = generateTxOpts(registrant, registrantKey, "0")
	require.Nil(t, err, "Failed to generate transaction options")
	fmt.Println("About to Set Records")
	tx, err = dnsresolver.SetRecords(opts, []byte("0x"+hex.EncodeToString([]byte(arec.String()))))
	fmt.Println("Set DNS Records")
	fmt.Printf("err: %+v\n", err)
	fmt.Printf("tx: %+v\n", tx)
	require.Nil(t, err, "Failed to update A record")
	waitForTransaction(tx.Hash())

	// Retrieve A Record
}
