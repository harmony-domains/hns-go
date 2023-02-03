package onens

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestARegistration(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	domain := unregisteredDomain()
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

func TestARecordCreateImport(t *testing.T) {
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
	arecO := []byte{0x1, 0x61, 0x4, 0x74, 0x65, 0x73, 0x74, 0x7, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0xe, 0x10, 0x0, 0x4, 0x80, 0x0, 0x0, 0x1}
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
	assert.Equal(t, existingRec, arecO, "Failed to clear existing records")

}

func TestARecordCreateRR(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	domain := "test.country"
	name := "a.test.country."

	// record imported from ens-deployer testing
	aRecTest := "0161047465737407636f756e747279000001000100000e10000480000001"
	arecTestBytes := []byte{0x1, 0x61, 0x4, 0x74, 0x65, 0x73, 0x74, 0x7, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0xe, 0x10, 0x0, 0x4, 0x80, 0x0, 0x0, 0x1}

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
	// TODO format our own records using https://pkg.go.dev/github.com/miekg/dns
	// blog https://miek.nl/2014/september/21/idn-and-private-rr-in-go-dns/
	// js code to review and potentially migrate https://www.npmjs.com/package/dns-js
	aRec := new(dns.A)
	aRec.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
	aRec.A = []byte{128, 0, 0, 1}
	u := new(dns.RFC3597)
	u.ToRFC3597(aRec)
	aRecStringHex := hex.EncodeToString([]byte(aRec.String()))
	aRecBytes, err := hex.DecodeString(aRecStringHex)
	assert.Equal(t, aRecStringHex, aRecTest, "Failed to create correct A record")
	assert.Equal(t, aRecBytes, arecTestBytes, "Failed to clear corect A record Bytes")

	opts, err = generateTxOpts(registrant, registrantKey, "0")
	tx, err = dnsresolver.SetRecords(opts, aRecBytes)
	require.Nil(t, err, "Failed to update A record")
	waitForTransaction(tx.Hash())

	// Check A Record Created
	existingRec, err = dnsresolver.Record("a.test.country.", 1)
	require.Nil(t, err, "Failed to get records")
	assert.Equal(t, existingRec, arecTestBytes, "Failed to clear existing records")

}
