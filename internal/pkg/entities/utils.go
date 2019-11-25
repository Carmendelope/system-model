/*
 * Copyright 2019 Nalej
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

package entities

import (
	"github.com/satori/go.uuid"
	"math/rand"
)

// GenerateUUID generates a new UUID.
func GenerateUUID() string {
	return uuid.NewV4().String()
}

// GenerateInt64 generates a random int64 variable.
func GenerateInt64() int64 {
	return rand.Int63()
}
