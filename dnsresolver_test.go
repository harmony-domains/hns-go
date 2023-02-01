package onens

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestARecord(t *testing.T) {
	registrant := tconfig.testAccounts.aliceAddress
	registrantKey := tconfig.testAccounts.alicePrivateKey
	// domain := unregisteredDomain()
	domain := "a.test.country."
	// name, err := NewName(tclient, domain)
	// require.Nil(t, err, "Failed to create name")

	// // Register stage 1
	// opts, err := generateTxOpts(registrant, registrantKey, "0")
	// require.Nil(t, err, "Failed to generate transaction options")
	// tx, secret, err := name.RegisterStageOne(registrant, tconfig.duration, opts)
	// require.Nil(t, err, "Failed to send stage one transaction")
	// // // Wait until mined
	// waitForTransaction(tx.Hash())

	// // Wait until ready to submit stage 2
	// interval, err := name.RegistrationInterval()
	// require.Nil(t, err, "Failed to obtain registration interval")
	// time.Sleep(interval)
	// // Sleep for 10 more seconds to ensure we are over the interval
	// time.Sleep(10 * time.Second)

	// // Register stage 2
	// opts, err = generateTxOpts(registrant, registrantKey, "1200 Ether")
	// require.Nil(t, err, "Failed to generate transaction options")
	// tx, err = name.RegisterStageTwo(registrant, tconfig.duration, secret, opts)
	// require.Nil(t, err, "Failed to send stage two transaction")
	// // // Wait until mined
	// waitForTransaction(tx.Hash())

	// // Confirm registered
	// registeredRegistrant, err := name.Registrant()
	// require.Nil(t, err, "Failed to obtain registrant")
	// assert.Equal(t, tconfig.NameWrapper, registeredRegistrant, "failed to register name")

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
	arec, err := dns.NewRR("a.test.country. 3600 IN A 1.2.3.4")
	require.Nil(t, err, "Failed to create a record  for domain")
	// arec, err := dns.NewRR(domain + " 3600 IN A 1.2.3.4")
	fmt.Printf("arec                                                   : %+v\n", arec)
	fmt.Printf("arec.Header()                                          : %+v\n", arec.Header())
	fmt.Printf("arec.String()                                          : %+v\n", arec.String())
	fmt.Printf("[]byte(arec.String())                                  : %+v\n", []byte(arec.String()))
	fmt.Printf("hex.EncodeToString([]byte(arec.String()))              : %+v\n", hex.EncodeToString([]byte(arec.String())))
	fmt.Printf("0x + hex.EncodeToString([]byte(arec.String()))         : %+v\n", ("0x" + hex.EncodeToString([]byte(arec.String()))))
	fmt.Printf("[]byte(0x + hex.EncodeToString([]byte(arec.String()))) : %+v\n", []byte("0x"+hex.EncodeToString([]byte(arec.String()))))
	// fmt.Printf("arec.unpack(): %+v\n", arec.unpack())

	dnsresolver, err := NewDNSResolver(tconfig.client, domain)
	require.Nil(t, err, "Failed to create resolver for domain: %s", domain)
	opts, err := generateTxOpts(registrant, registrantKey, "0")
	require.Nil(t, err, "Failed to generate transaction options")
	fmt.Println("About to Set Records")
	tx, err := dnsresolver.SetRecords(opts, []byte("0x"+hex.EncodeToString([]byte(arec.String()))))
	waitForTransaction(tx.Hash())
	// // func (r *DNSResolver) SetRecords(opts *bind.TransactOpts, []byte) (*types.Transaction, error)

	// Retrieve A Record
}
