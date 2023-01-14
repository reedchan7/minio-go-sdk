//go:build example
// +build example

/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2018 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/reedchan7/minio-go-sdk/v7"
	"github.com/reedchan7/minio-go-sdk/v7/pkg/credentials"
	"github.com/minio/sio"
	"golang.org/x/crypto/argon2"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-testfile, my-bucketname and
	// my-objectname are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := minio.New("s3.amazonaws.com", &minio.Options{
		Creds:  credentials.NewStaticV4("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	obj, err := s3Client.GetObject(context.Background(), "my-bucketname", "my-objectname", minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	localFile, err := os.Create("my-testfile")
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()

	password := []byte("myfavoritepassword")                    // Change as per your needs.
	salt := []byte(path.Join("my-bucketname", "my-objectname")) // Change as per your needs.
	_, err = sio.Decrypt(localFile, obj, sio.Config{
		Key: argon2.IDKey(password, salt, 1, 64*1024, 4, 32), // generate a 256 bit long key.
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully decrypted 'my-objectname' to local file 'my-testfile'")
}
