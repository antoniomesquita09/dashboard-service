package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MemoryModel struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Used      float64            `json:"used,omitempty" validate:"required"`
	Committed float64            `json:"committed,omitempty" validate:"required"`
	Total     float64            `json:"total,omitempty" validate:"required"`
	DateTime  string             `json:"date_time,omitempty" validate:"required"`
}
