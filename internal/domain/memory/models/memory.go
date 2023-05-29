package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MemoryModel struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Used      float64            `json:"name,omitempty" validate:"required"`
	Committed float64            `json:"location,omitempty" validate:"required"`
	Total     float64            `json:"title,omitempty" validate:"required"`
}
