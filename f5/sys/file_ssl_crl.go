// Copyright e-Xpert Solutions SA. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sys

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/e-XpertSolutions/f5-rest-client/f5"
)

// FileSSLCRLConfigList holds a list of FileSSLCRL configuration.
type FileSSLCRLConfigList struct {
	Items    []FileSSLCRLConfig `json:"items"`
	Kind     string             `json:"kind"`
	SelfLink string             `json:"selflink"`
}

// FileSSLCRLConfig holds the configuration of a single FileSSLCRL.
type FileSSLCRLConfig struct {
}

// FileSSLCRLEndpoint represents the REST resource for managing FileSSLCRL.
const FileSSLCRLEndpoint = "/file/ssl-crl"

// FileSSLCRLResource provides an API to manage FileSSLCRL configurations.
type FileSSLCRLResource struct {
	c f5.Client
}

// ListAll  lists all the FileSSLCRL configurations.
func (r *FileSSLCRLResource) ListAll() (*FileSSLCRLConfigList, error) {
	var list FileSSLCRLConfigList
	if err := r.c.ReadQuery(BasePath+FileSSLCRLEndpoint, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get a single FileSSLCRL configuration identified by id.
func (r *FileSSLCRLResource) Get(id string) (*FileSSLCRLConfig, error) {
	var item FileSSLCRLConfig
	if err := r.c.ReadQuery(BasePath+FileSSLCRLEndpoint, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Create a new FileSSLCRL configuration.
func (r *FileSSLCRLResource) Create(name, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to gather information about '%s': %v", path, err)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to read file from path: %v", err)
	}
	defer f.Close()

	req, err := r.c.MakeUploadRequest(f5.UploadRESTPath+"/"+filepath.Base(path), f, info.Size())
	if err != nil {
		return fmt.Errorf("failed to create upload request: %v", err)
	}
	resp, err := r.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file '%s': %v", path, err)
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	log.Print("DEBUG resp=", string(bs))

	data := map[string]string{
		"name":        name,
		"source-path": "file://localhost/var/config/rest/downloads/" + filepath.Base(path),
	}
	if err := r.c.ModQuery("POST", BasePath+FileSSLCRLEndpoint, data); err != nil {
		return fmt.Errorf("failed to create FileSSLCRL configuration: %v", err)
	}

	return nil
}

// Edit a FileSSLCRL configuration identified by id.
func (r *FileSSLCRLResource) Edit(id, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to gather information about '%s': %v", path, err)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to read file from path: %v", err)
	}
	defer f.Close()

	req, err := r.c.MakeUploadRequest(f5.UploadRESTPath+"/"+filepath.Base(path), f, info.Size())
	if err != nil {
		return fmt.Errorf("failed to create upload request: %v", err)
	}
	resp, err := r.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file '%s': %v", path, err)
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	log.Print("DEBUG resp=", string(bs))

	data := map[string]string{
		"source-path": "file://localhost/var/config/rest/downloads/" + filepath.Base(path),
	}
	if err := r.c.ModQuery("PUT", BasePath+FileSSLCRLEndpoint+"/"+id, data); err != nil {
		return fmt.Errorf("failed to create FileSSLCRL configuration: %v", err)
	}

	return nil
}

// Delete a single FileSSLCRL configuration identified by id.
func (r *FileSSLCRLResource) Delete(id string) error {
	if err := r.c.ModQuery("DELETE", BasePath+FileSSLCRLEndpoint+"/"+id, nil); err != nil {
		return err
	}
	return nil
}
