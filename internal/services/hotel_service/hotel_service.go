package hotel_service

type HotelService interface {
	GetHotels(hotelDataFilePath string, ids []string, destinations []int) ([]Hotel, error)
	UpdateHotelsFromSuppliers(suppliersFilePath string, hotelDataFilePath string) ([]string, error)
}

type Location struct {
	Lat     float64 `json:"lat,omitempty"`
	Long    float64 `json:"lng,omitempty"`
	Address string  `json:"address,omitempty"`
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
}

type Amenities struct {
	General []string `json:"general,omitempty"`
	Room    []string `json:"room,omitempty"`
}

type Image struct {
	Link        string `json:"link,omitempty"`
	Description string `json:"description,omitempty"`
}

type Hotel struct {
	ID               string             `json:"id,omitempty"`
	DestinationID    int                `json:"destination_id,omitempty"`
	HotelName        string             `json:"hotel_name,omitempty"`
	Location         Location           `json:"location,omitempty"`
	Description      []string           `json:"description,omitempty"`
	Amenities        Amenities          `json:"amenities,omitempty"`
	Images           map[string][]Image `json:"images,omitempty"`
	BookingCondition []string           `json:"booking_condition,omitempty"`
}
