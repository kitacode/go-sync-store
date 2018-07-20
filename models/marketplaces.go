package models

type Marketplaces struct {
    Marketplaces        Marketplace
}

type Marketplace struct {
    Id                  string      `bson:"_id,omitempty"`
    Name                string
}
