package data

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func Init() {

}

func Test() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("cassandra-seed", "cassandra-node-1", "cassandra-node-2") //replace PublicIP with the IP addresses used by your cluster.
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = time.Second * 10
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 3}
	//cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	//defer session.Close()

	// create keyspaces
	err = session.Query("create keyspace if not exists example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").Exec()
	if err != nil {
		log.Println(err)
		return nil
	}

	fmt.Println("finished creating keyspace")

	// create table
	err = session.Query("CREATE TABLE IF NOT EXISTS example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));").Exec()
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println("made table")

	// err = session.Query("create index on example.tweet(timeline);").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println("index table")

	// insert some practice data
	err = session.Query(`INSERT INTO example.tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec()
	err = session.Query(`INSERT INTO example.tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec()
	err = session.Query(`INSERT INTO example.tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec()
	err = session.Query(`INSERT INTO example.tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec()

	if err != nil {
		log.Println(err)
		return nil
	}

	fmt.Println("inserted stuff")
	// // Return stuff
	var id gocql.UUID
	var text string
	fmt.Println("querying")
	iter := session.Query("SELECT id, text FROM example.tweet;").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return session

}
