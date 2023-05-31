package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CPUModel struct from gateway-service metric proxy
type CPUModel struct {
	ID         primitive.ObjectID `json:"id,omitempty"`
	Percentage float64            `json:"percentage"`
}
