// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package c2

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/db"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	// ListenerBucketName - stores listeners
	ListenerBucketName = "listeners"
	// ListenerNamespace - Listener namespace
	ListenerNamespace = "listener"
	storageLog        = log.ServerLogger("listeners", "persistence")
)

// ListenerConfig - Stores the configuration for a listener
type ListenerConfig struct {
	// ID
	ID int32 `json:"id"`
	// All
	LHost       string `json:"lhost,omitempty"`
	LPort       uint16 `json:"lport,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// DNS
	DNSDomains     []string `json:"dns_domains,omitempty"`
	EnableCanaries bool     `json:"enable_canaries,omitempty"`

	// HTTPS
	HTTPDomain  string `json:"http_domain,omitempty"`
	Secure      bool   `json:"secure,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	Key         string `json:"key,omitempty"`
	LetsEncrypt bool   `json:"lets_encrypt,omitempty"`
	Website     string `json:"website,omitempty"`

	// Stagers
	ImplantStage string `json:"implant_stage,omitempty"`
}

// PersistMTLS - mTLS listener
func PersistMTLS(job *core.Job, lhost string) error {

	listener := &ListenerConfig{
		ID:          rand.Int31(),
		LHost:       lhost,
		LPort:       job.Port,
		Protocol:    job.Protocol,
		Name:        job.Name,
		Description: job.Description,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	key := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(key, rawListener)
}

// PersistDNS - DNS
func PersistDNS(job *core.Job, enableCanaries bool, dnsDomains []string) error {

	listener := &ListenerConfig{
		ID:             rand.Int31(),
		LPort:          job.Port,
		Protocol:       job.Protocol,
		Name:           job.Name,
		Description:    job.Description,
		DNSDomains:     dnsDomains,
		EnableCanaries: enableCanaries,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	key := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(key, rawListener)
}

// PersistHTTPS - HTTPS
func PersistHTTPS(job *core.Job, lhost string, cert string, key string, secure bool, domain string, website string, letsEncrypt bool) error {

	listener := &ListenerConfig{
		ID:          rand.Int31(),
		LHost:       lhost,
		LPort:       job.Port,
		Protocol:    job.Protocol,
		Name:        job.Name,
		Description: job.Description,
		Secure:      secure,
		Certificate: cert,
		Key:         key,
		HTTPDomain:  domain,
		Website:     website,
		LetsEncrypt: letsEncrypt,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	bucketKey := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(bucketKey, rawListener)
}

// PersistTCPStager - TCP stager
func PersistTCPStager(job *core.Job, lhost string, implantStage string) error {

	listener := &ListenerConfig{
		ID:           rand.Int31(),
		LHost:        lhost,
		LPort:        job.Port,
		Protocol:     job.Protocol,
		Name:         job.Name,
		Description:  job.Description,
		ImplantStage: implantStage,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	key := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(key, rawListener)
}

// PersistHTTPStager - HTTP
func PersistHTTPStager(job *core.Job, lhost string, implantStage string) error {

	listener := &ListenerConfig{
		ID:           rand.Int31(),
		LHost:        lhost,
		LPort:        job.Port,
		Protocol:     job.Protocol,
		Name:         job.Name,
		Description:  job.Description,
		ImplantStage: implantStage,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	key := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(key, rawListener)
}

// PersistHTTPSStager - HTTPS
func PersistHTTPSStager(job *core.Job, lhost string, implantStage string) error {

	listener := &ListenerConfig{
		ID:           rand.Int31(),
		LHost:        lhost,
		LPort:        job.Port,
		Protocol:     job.Protocol,
		Name:         job.Name,
		Description:  job.Description,
		ImplantStage: implantStage,
	}

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}
	rawListener, err := json.Marshal(listener)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved persistent listener (name: %s proto: %s, address: %s:%d) ",
		listener.Name, listener.Protocol, listener.LHost, listener.LPort)

	key := fmt.Sprintf("%s.%s_%s-%d_%d", ListenerNamespace, listener.Protocol, listener.LHost, listener.LPort, listener.ID)
	return bucket.Set(key, rawListener)
}
