package domains

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	database "systems-management-api/core/database"
)

// UserService service which provides methos to access and modify database data
type DomainService struct{}

// Retrieves all domains instances
func (service *DomainService) all() (*[]Domain, error) {
	db := database.DB()
	collection := db.D.Collection("domain")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return nil, err
	} else {
		domains := []Domain{}
		for cursor.Next(context.TODO()) {
			var domain Domain
			cursor.Decode(&domain)
			domains = append(domains, domain)
		}
		return &domains, nil
	}
}

// Retrieves a domain instance given its ID
func (service *DomainService) GetById(id string) (*Domain, error) {
	db := database.DB()
	domain := Domain{}

	if err := db.GetById("domain", id, &domain); err != nil {
		return nil, err
	} else {
		return &domain, nil
	}
}

// Saves the domain model to database
// Returns boolean result and error
func (service *DomainService) Save(domain *Domain) (bool, error) {
	db := database.DB()
	collection := db.D.Collection("domain")

	if domain.ID.IsZero() {
		// insert
		res, err := collection.InsertOne(context.TODO(), domain)

		if err != nil {
			zap.S().Error("Error inserting domain: ", err)
			return false, err
		} else {
			zap.S().Info(fmt.Sprintf("Domain %s inserted succesfully", domain.Name))
			// update user ID
			domain.ID = res.InsertedID.(primitive.ObjectID)
			return true, nil
		}
	} else {
		// update
		filter := bson.M{"_id": domain.ID}
		_, err := collection.ReplaceOne(context.TODO(), filter, domain)

		if err != nil {
			zap.S().Error("Error inserting domain: ", err)
			return false, err
		} else {
			zap.S().Info(fmt.Sprintf("Domain %s updated succesfully", domain.Name))
			// update user ID
			return true, nil
		}
	}
}

// Deletes the domain model from databse
// Returns boolean result and error
func (service *DomainService) Delete(domain *Domain) (bool, error) {
	db := database.DB()
	collection := db.D.Collection("domain")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": domain.ID})

	if err != nil {
		zap.S().Error("Error deleting domain: ", err)
		return false, err
	} else {
		zap.S().Info(fmt.Sprintf("Domain %s deleted succesfully", domain.Name))
		// update user ID
		return true, nil
	}
}
