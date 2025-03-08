package models

type FeedbackDto struct {
	Name                string `json:"name" validate:"required"`
	Email               string `json:"email" validate:"required,email"`
	CityOrRegion        string `json:"cityOrRegion"`
	Organization        string `json:"organization"`
	Position            string `json:"position"`
	GratitudeExpression string `json:"gratitudeExpression"`
	Message             string `json:"message" validate:"required"`
}
