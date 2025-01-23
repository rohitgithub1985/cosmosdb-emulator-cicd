package main

import (
	"context"
	"demo/user"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stretchr/testify/require"
)

var container *azcosmos.ContainerClient


func TestMain(m *testing.M) {

	ctx := context.Background()

	connectionString := os.Getenv("COSMOSDB_CONNECTION_STRING")

    client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: os.Getenv("COSMOSDB_DATABASE_NAME")}, nil)
    if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	log.Println("Database created")


	db, err := client.NewDatabase(os.Getenv("COSMOSDB_DATABASE_NAME"))
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	_, err = db.CreateContainer(ctx, azcosmos.ContainerProperties{ID: os.Getenv("COSMOSDB_CONTAINER_NAME"), PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{Paths: []string{"/id"}}}, nil)

	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}

	log.Println("Container created")

	container, err = client.NewContainer(os.Getenv("COSMOSDB_DATABASE_NAME"), os.Getenv("COSMOSDB_CONTAINER_NAME"))
	if err != nil {
		log.Fatalf("failed to get container: %v", err)
	}

	m.Run()

	_, err = db.Delete(ctx, nil)
	if err != nil {
		log.Fatalf("failed to delete database: %v", err)
	}

	log.Println("Database deleted")

}

func TestCreate(t *testing.T) {
    ctx := context.Background()

	var newUser = user.User{
		ID: "42",
		Email: "user42@test.com",
		Active: true,
	}

	err := user.Create(ctx, container, newUser)

    require.NoError(t, err)
}

func TestReadItem(t *testing.T) {
    ctx := context.Background()
	userID := "42"

    userItem, err := user.Read(ctx, container, userID)

    require.NoError(t, err)
	require.Equal(t, userID, userItem.ID)
}

