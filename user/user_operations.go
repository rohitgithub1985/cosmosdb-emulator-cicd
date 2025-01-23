package user

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type User struct {
    ID string `json:"id"`
    Email string `json:"email"`
    Active bool `json:"active"`
}

func Create(ctx context.Context, container *azcosmos.ContainerClient, user User) error {

    userBytes, err := json.Marshal(user)
    if err != nil {
        return err
    }

    //log.Println("Creating user", string(userBytes))

    _, err = container.CreateItem(ctx, azcosmos.NewPartitionKeyString(user.ID), userBytes, nil)
    return err
}

func Read(ctx context.Context, container *azcosmos.ContainerClient, id string) (User, error) {
    itemResponse, err := container.ReadItem(ctx, azcosmos.NewPartitionKeyString(id), id, nil)
    if err != nil {
        return User{}, err
    }
    var user User
    err = json.Unmarshal(itemResponse.Value, &user)
    return user, err
}