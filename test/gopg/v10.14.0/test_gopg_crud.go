// Copyright (c) 2025 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/alibaba/loongsuite-go-agent/test/verifier"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

var db *pg.DB

type User struct {
	ID   string
	Name string
	Age  uint8
}

func TestCreateTable() {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)`); err != nil {
		log.Printf("%v", err)
	}
}

func TestInsert() {
	user := User{Name: "opentelemetry", Age: 18}
	if _, err := db.Model(&user).Insert(); err != nil {
		log.Printf("%v", err)
	}
}

func TestQuery() {
	var users []User
	if err := db.Model(&users).Select(); err != nil {
		log.Printf("%v", err)
	}
}

func TestUpdate() {
	if _, err := db.Model(&User{ID: "1", Age: 10}).WherePK().Update("name", "hello"); err != nil {
		log.Printf("%v", err)
	}
}

func TestDelete() {
	if _, err := db.Model(&User{ID: "1"}).WherePK().Delete(); err != nil {
		log.Printf("%v", err)
	}
}

func TestDropTable() {
	if err := db.Model(&User{}).DropTable(&orm.DropTableOptions{}); err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	db = pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:" + os.Getenv("POSTGRES_PORT"),
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	TestCreateTable()
	TestInsert()
	TestQuery()
	TestUpdate()
	TestDelete()
	TestDropTable()
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "postgresql", "postgresql", "127.0.0.1", "CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)", "", "", nil)
		verifier.VerifyDbAttributes(stubs[1][0], "INSERT", "postgresql", "127.0.0.1", "INSERT INTO \"users\" (\"id\", \"name\", \"age\") VALUES (DEFAULT, 'opentelemetry', 18)", "INSERT", "", nil)
		verifier.VerifyDbAttributes(stubs[2][0], "SELECT", "postgresql", "127.0.0.1", "SELECT \"user\".\"id\", \"user\".\"name\", \"user\".\"age\" FROM \"users\" AS \"user\"", "SELECT", "", nil)
		verifier.VerifyDbAttributes(stubs[3][0], "UPDATE", "postgresql", "127.0.0.1", "UPDATE \"users\" AS \"user\" SET \"name\" = NULL, \"age\" = 10 WHERE \"user\".\"id\" = '1'", "UPDATE", "", nil)
		verifier.VerifyDbAttributes(stubs[4][0], "DELETE", "postgresql", "127.0.0.1", "DELETE FROM \"users\" AS \"user\" WHERE \"user\".\"id\" = '1'", "DELETE", "", nil)
		verifier.VerifyDbAttributes(stubs[5][0], "DROP TABLE", "postgresql", "127.0.0.1", "DROP TABLE \"users\"", "DROP TABLE", "", nil)
	}, 1)
}
