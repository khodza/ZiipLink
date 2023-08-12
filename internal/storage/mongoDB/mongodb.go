package mongodb

import (
	"context"
	"errors"
	"fmt"
	"zipinit/internal/lib/random"
	"zipinit/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client          *mongo.Client
	database        *mongo.Database
	collection      *mongo.Collection
	aliasCollection *mongo.Collection
}

func NewStorage(connectionString, dbName string) (*Storage, error) {
	const op = "storage.mongodb.New"
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	database := client.Database(dbName)
	collection := database.Collection("urls")
	aliasCollection := database.Collection("aliases")

	// Create unique indexes on 'alias'
	_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"alias": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create indexes: %w", op, err)
	}
	// Create unique indexes on 'alias'
	_, err = aliasCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"alias": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create indexes: %w", op, err)
	}

	return &Storage{
		client:          client,
		database:        database,
		collection:      collection,
		aliasCollection: aliasCollection,
	}, nil
}

func (s *Storage) SaveAndGenerateRandomStrings(maxAliases int) error {
	const op = "storage.mongodb.GenerateRandomStrings"

	const aliasLength = 6
	aliasCount := 0

	for aliasCount < maxAliases {
		randomString := random.NewRandomString(aliasLength)
		fmt.Println("randomString: ", randomString)
		_, err := s.aliasCollection.InsertOne(context.Background(), bson.M{
			"alias": randomString,
			"used":  false,
		})
		if err != nil {
			fmt.Println("err: ", err)
			if mongo.IsDuplicateKeyError(err) {
				continue // Generate a new random string and try again
			}
			return fmt.Errorf("%s: %w", op, err)
		}
		aliasCount++
	}

	return nil
}

func (s *Storage) SaveUrl(urlToSave, providedAlias string) (string, error) {
	const op = "storage.mongodb.SaveUrl"
	_, err := s.collection.InsertOne(context.Background(), bson.M{
		"url":   urlToSave,
		"alias": providedAlias,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", storage.ErrUrlExists
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return providedAlias, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.mongodb.GetUrl"
	var result struct {
		URL string `bson:"url"`
	}
	err := s.collection.FindOne(context.Background(), bson.M{
		"alias": alias,
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result.URL, nil
}

func (s *Storage) GetAnyUnusedAlias() (string, error) {
	const op = "storage.mongodb.GetUnusedAlias"
	var result struct {
		Alias string `bson:"alias"`
	}
	err := s.aliasCollection.FindOne(context.Background(), bson.M{
		"used": false,
	}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", storage.ErrAliasNotFound // No unused documents found
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return result.Alias, nil
}

func (s *Storage) MakeAliasUsed(alias string) error {
	const op = "storage.mongodb.MarkAliasAsUsed"

	filter := bson.M{"alias": alias}
	update := bson.M{"$set": bson.M{"used": true}}

	result, err := s.aliasCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: Alias not found", op)
	}

	return nil
}

func (s *Storage) SaveNewAlias(alias string) error {
	const op = "storage.mongodb.SaveNewAlias"
	// Insert the provided alias into the alias collection
	_, err := s.aliasCollection.InsertOne(context.Background(), bson.M{
		"alias": alias,
		"used":  true,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return storage.ErrAliasExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
