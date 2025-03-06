package hotel_service

import (
	"ascenda-loyalty-assignment/pkg/logging"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	Responses map[string]*http.Response
	Errors    map[string]error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if err, ok := m.Errors[req.URL.String()]; ok {
		return nil, err
	}
	if resp, ok := m.Responses[req.URL.String()]; ok {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func TestGetAllHotels(t *testing.T) {
	testCases := []struct {
		description  string
		dataFilePath func() string
		expectedErr  bool
		ids          []string
		destinations []int
		expectedData []Hotel
	}{
		{
			description:  "fail to get data from invalid data file path",
			ids:          []string{"iJhz"},
			destinations: []int{5432},
			dataFilePath: func() string {
				return "invalid"
			},
			expectedErr: true,
		},
		{
			description: "return empty hotel list when not parsing any ids or destinations",
			dataFilePath: func() string {
				return "invalid"
			},
			expectedErr:  false,
			expectedData: []Hotel{},
		},
		{
			description: "successfully get the correct data from data file path",
			dataFilePath: func() string {
				wd, _ := os.Getwd()
				filePath := filepath.Join(wd, "test_data", "test_hotels.json")
				return filePath
			},
			ids:          []string{"iJhz"},
			destinations: []int{5432},
			expectedErr:  false,
			expectedData: []Hotel{
				{
					ID:            "iJhz",
					DestinationID: 5432,
					HotelName:     "Beach Villas Singapore",
					Location: Location{
						Lat:     1.264751,
						Long:    103.824006,
						Address: "8 Sentosa Gateway, Beach Villas, 098269",
						Country: "Singapore",
					},
					Description: []string{
						"This 5 star hotel is located on the coastline of Singapore.",
						"Located at the western tip of Resorts World Sentosa, guests at the Beach Villas are guaranteed privacy while they enjoy spectacular views of glittering waters. Guests will find themselves in paradise with this series of exquisite tropical sanctuaries, making it the perfect setting for an idyllic retreat. Within each villa, guests will discover living areas and bedrooms that open out to mini gardens, private timber sundecks and verandahs elegantly framing either lush greenery or an expanse of sea. Guests are assured of a superior slumber with goose feather pillows and luxe mattresses paired with 400 thread count Egyptian cotton bed linen, tastefully paired with a full complement of luxurious in-room amenities and bathrooms boasting rain showers and free-standing tubs coupled with an exclusive array of ESPA amenities and toiletries. Guests also get to enjoy complimentary day access to the facilities at Asia’s flagship spa – the world-renowned ESPA.",
						"Surrounded by tropical gardens, these upscale villas in elegant Colonial-style buildings are part of the Resorts World Sentosa complex and a 2-minute walk from the Waterfront train station. Featuring sundecks and pool, garden or sea views, the plush 1- to 3-bedroom villas offer free Wi-Fi and flat-screens, as well as free-standing baths, minibars, and tea and coffeemaking facilities. Upgraded villas add private pools, fridges and microwaves; some have wine cellars. A 4-bedroom unit offers a kitchen and a living room. There's 24-hour room and butler service. Amenities include posh restaurant, plus an outdoor pool, a hot tub, and free parking.",
					},
					Amenities: Amenities{
						General: []string{
							"Pool", "BusinessCenter", "WiFi", "DryCleaning", "Breakfast", "Aircon", "Tv", "Coffee machine", "Kettle",
							"Hair dryer", "Iron", "Tub", "outdoor pool", "indoor pool", "business center", "childcare",
						},
						Room: []string{"tv", "coffee machine", "kettle", "hair dryer", "iron"},
					},
					Images: map[string][]Image{
						"amenities": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/0.jpg", Description: "RWS"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/6.jpg", Description: "Sentosa Gateway"},
						},
						"rooms": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/2.jpg", Description: "Double room"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/4.jpg", Description: "Bathroom"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/3.jpg", Description: "Double room"},
						},
						"site": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/1.jpg", Description: "Front"},
						},
					},
					BookingCondition: []string{
						"All children are welcome. One child under 12 years stays free of charge when using existing beds. One child under 2 years stays free of charge in a child's cot/crib. One child under 4 years stays free of charge when using existing beds. One older child or adult is charged SGD 82.39 per person per night in an extra bed. The maximum number of children's cots/cribs in a room is 1. There is no capacity for extra beds in the room.",
						"Pets are not allowed.",
						"WiFi is available in all areas and is free of charge.",
						"Free private parking is possible on site (reservation is not needed).",
						"Guests are required to show a photo identification and credit card upon check-in. Please note that all Special Requests are subject to availability and additional charges may apply. Payment before arrival via bank transfer is required. The property will contact you after you book to provide instructions. Please note that the full amount of the reservation is due before arrival. Resorts World Sentosa will send a confirmation with detailed payment information. After full payment is taken, the property's details, including the address and where to collect keys, will be emailed to you. Bag checks will be conducted prior to entry to Adventure Cove Waterpark. === Upon check-in, guests will be provided with complimentary Sentosa Pass (monorail) to enjoy unlimited transportation between Sentosa Island and Harbour Front (VivoCity). === Prepayment for non refundable bookings will be charged by RWS Call Centre. === All guests can enjoy complimentary parking during their stay, limited to one exit from the hotel per day. === Room reservation charges will be charged upon check-in. Credit card provided upon reservation is for guarantee purpose. === For reservations made with inclusive breakfast, please note that breakfast is applicable only for number of adults paid in the room rate. Any children or additional adults are charged separately for breakfast and are to paid directly to the hotel.",
					},
				},
				{
					ID:            "SjyX",
					DestinationID: 5432,
					HotelName:     "InterContinental",
					Location: Location{
						Address: "1 Nanson Rd, Singapore 238909",
						Country: "Singapore",
					},
					Description: []string{
						"Enjoy sophisticated waterfront living at the new InterContinental® Singapore Robertson Quay, luxury's preferred address nestled in the heart of Robertson Quay along the Singapore River, with the CBD just five minutes drive away. Magnifying the comforts of home, each of our 225 studios and suites features a host of thoughtful amenities that combine modernity with elegance, whilst maintaining functional practicality. The hotel also features a chic, luxurious Club InterContinental Lounge.",
						"InterContinental Singapore Robertson Quay is luxury's preferred address offering stylishly cosmopolitan riverside living for discerning travelers to Singapore. Prominently situated along the Singapore River, the 225-room inspiring luxury hotel is easily accessible to the Marina Bay Financial District, Central Business District, Orchard Road and Singapore Changi International Airport, all located a short drive away. The hotel features the latest in Club InterContinental design and service experience, and five dining options including Publico, an Italian landmark dining and entertainment destination by the waterfront.",
					},
					Amenities: Amenities{
						General: []string{
							"Pool", "WiFi", "Aircon", "BusinessCenter", "BathTub", "Breakfast", "DryCleaning", "Bar",
							"outdoor pool", "business center", "childcare", "parking", "bar", "dry cleaning", "wifi",
							"breakfast", "concierge",
						},
						Room: []string{"aircon", "minibar", "tv", "bathtub", "hair dryer"},
					},
					Images: map[string][]Image{
						"rooms": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i93_m.jpg", Description: "Double room"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i94_m.jpg", Description: "Bathroom"},
						},
						"site": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i1_m.jpg", Description: "Restaurant"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i2_m.jpg", Description: "Hotel Exterior"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i5_m.jpg", Description: "Entrance"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i24_m.jpg", Description: "Bar"},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			logger := logging.LogrusLogger()
			ctx := context.Background()
			hotelService := NewHotelService(logger, nil, ctx)
			hotels, err := hotelService.GetHotels(tc.dataFilePath(), tc.ids, tc.destinations)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.ElementsMatch(t, tc.expectedData, hotels)
			}
		})
	}
}

func TestUpdateHotelsFromSuppliers(t *testing.T) {
	testCases := []struct {
		description            string
		suppliersFilePath      func() string
		hotelDataFilePath      func() string
		httpClient             func() HTTPClient
		ids                    []string
		destinations           []int
		expectedErr            bool
		expectedData           []Hotel
		expectedFetchedSources []string
	}{
		{
			description: "fail to read suppliers file",
			suppliersFilePath: func() string {
				return "invalid_path"
			},
			hotelDataFilePath: func() string {
				return "valid_path"
			},
			httpClient: func() HTTPClient {
				mockClient := &MockHTTPClient{}
				return mockClient
			},
			expectedErr: true,
		},
		{
			description: "fail to parse suppliers data",
			suppliersFilePath: func() string {
				wd, _ := os.Getwd()
				filePath := filepath.Join(wd, "test_data", "invalid_suppliers.json")
				return filePath
			},
			hotelDataFilePath: func() string {
				return "valid_path"
			},
			httpClient: func() HTTPClient {
				mockClient := &MockHTTPClient{}
				return mockClient
			},
			expectedErr: true,
		},
		{
			description: "successfully update hotels from suppliers",
			suppliersFilePath: func() string {
				wd, _ := os.Getwd()
				filePath := filepath.Join(wd, "test_data", "test_suppliers.json")
				return filePath
			},
			hotelDataFilePath: func() string {
				wd, _ := os.Getwd()
				filePath := filepath.Join(wd, "test_data", "test_sanitize_data.json")
				return filePath
			},
			httpClient: func() HTTPClient {
				mockSupplierData := `[
		{
			"Id": "SjyX",
			"DestinationId": 5432,
			"Name": "InterContinental Singapore Robertson Quay",
			"Address": " 1 Nanson Road",
			"City": "Singapore",
			"Country": "SG",
			"PostalCode": "238909",
			"Description": "Enjoy sophisticated waterfront living at the new InterContinental® Singapore Robertson Quay.",
			"Facilities": [
				"Pool",
				"WiFi ",
				"Aircon",
				"BusinessCenter",
				"BathTub",
				"Breakfast",
				"DryCleaning",
				"Bar"
			]
		},
		{
			"Id": "f8c9",
			"DestinationId": 1122,
			"Name": "Hilton Shinjuku Tokyo",
			"Address": "160-0023, SHINJUKU-KU, 6-6-2 NISHI-SHINJUKU, JAPAN",
			"City": "Tokyo",
			"Country": "JP",
			"PostalCode": "160-0023",
			"Description": "Hilton Tokyo is located in Shinjuku, the heart of Tokyo's business, shopping and entertainment district, and is an ideal place to experience modern Japan. A complimentary shuttle operates between the hotel and Shinjuku station and the Tokyo Metro subway is connected to the hotel. Relax in one of the modern Japanese-style rooms and admire stunning city views. The hotel offers WiFi and internet access throughout all rooms and public space.",
			"Facilities": [
				"Pool",
				"WiFi ",
				"BusinessCenter",
				"DryCleaning",
				" Breakfast",
				"Bar",
				"BathTub"
			]
		},
		{
			"id": "f8c9",
			"destination": 1122,
			"name": "Hilton Tokyo Shinjuku",
			"lat": 35.6926,
			"lng": 139.690965,
			"images": {
				"rooms": [
					{
						"url": "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i10_m.jpg",
						"description": "Suite"
					},
					{
						"url": "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i11_m.jpg",
						"description": "Suite - Living room"
					}
				],
				"amenities": [
					{
						"url": "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i57_m.jpg",
						"description": "Bar"
					}
				]
			}
		},
		{
			"hotel_id": "iJhz",
			"destination_id": 5432,
			"hotel_name": "Beach Villas Singapore",
			"location": {
				"address": "8 Sentosa Gateway, Beach Villas, 098269",
				"country": "Singapore"
			},
			"amenities": {
				"general": [
					"outdoor pool",
					"indoor pool",
					"business center",
					"childcare"
				],
				"room": [
					"tv",
					"coffee machine",
					"kettle",
					"hair dryer",
					"iron"
				]
			},
			"images": {
				"rooms": [
					{
						"link": "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/2.jpg",
						"caption": "Double room"
					},
					{
						"link": "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/3.jpg",
						"caption": "Double room"
					}
				],
				"site": [
					{
						"link": "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/1.jpg",
						"caption": "Front"
					}
				]
			},
			"booking_conditions": [
				"Pets are not allowed.",
				"WiFi is available in all areas and is free of charge."
			]
		}
	]`
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(mockSupplierData)),
				}

				mockClient := &MockHTTPClient{
					Responses: map[string]*http.Response{
						"https://example.com/hotels": mockResponse,
					},
					Errors: map[string]error{},
				}

				return mockClient
			},
			expectedErr:            false,
			expectedFetchedSources: []string{"https://example.com/hotels"},
			expectedData: []Hotel{
				{
					ID:            "SjyX",
					DestinationID: 5432,
					HotelName:     "InterContinental Singapore Robertson Quay",
					Location: Location{
						Address: "1 Nanson Road",
						City:    "Singapore",
						Country: "SG",
					},
					Description: []string{
						"Enjoy sophisticated waterfront living at the new InterContinental® Singapore Robertson Quay.",
					},
					Amenities: Amenities{
						General: []string{
							"Pool", "WiFi", "Aircon", "BusinessCenter", "BathTub", "Breakfast", "DryCleaning", "Bar",
						},
					},
				},
				{
					ID:            "f8c9",
					DestinationID: 1122,
					HotelName:     "Hilton Tokyo Shinjuku",
					Location: Location{
						Lat:     35.6926,
						Long:    139.690965,
						Address: "160-0023, SHINJUKU-KU, 6-6-2 NISHI-SHINJUKU, JAPAN",
						City:    "Tokyo",
						Country: "JP",
					},
					Description: []string{
						"Hilton Tokyo is located in Shinjuku, the heart of Tokyo's business, shopping and entertainment district, and is an ideal place to experience modern Japan. A complimentary shuttle operates between the hotel and Shinjuku station and the Tokyo Metro subway is connected to the hotel. Relax in one of the modern Japanese-style rooms and admire stunning city views. The hotel offers WiFi and internet access throughout all rooms and public space.",
					},
					Amenities: Amenities{
						General: []string{
							"Pool", "WiFi", "BusinessCenter", "DryCleaning", "Breakfast", "Bar", "BathTub",
						},
					},
					Images: map[string][]Image{
						"amenities": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i57_m.jpg", Description: "Bar"},
						},
						"rooms": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i10_m.jpg", Description: "Suite"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/YwAr/i11_m.jpg", Description: "Suite - Living room"},
						},
					},
				},
				{
					ID:            "iJhz",
					DestinationID: 5432,
					HotelName:     "Beach Villas Singapore",
					Location: Location{
						Address: "8 Sentosa Gateway, Beach Villas, 098269",
						Country: "Singapore",
					},
					Amenities: Amenities{
						General: []string{
							"outdoor pool", "indoor pool", "business center", "childcare",
						},
						Room: []string{
							"tv", "coffee machine", "kettle", "hair dryer", "iron",
						},
					},
					Images: map[string][]Image{
						"rooms": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/2.jpg", Description: "Double room"},
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/3.jpg", Description: "Double room"},
						},
						"site": {
							{Link: "https://d2ey9sqrvkqdfs.cloudfront.net/0qZF/1.jpg", Description: "Front"},
						},
					},
					BookingCondition: []string{
						"Pets are not allowed.",
						"WiFi is available in all areas and is free of charge.",
					},
				},
			},
			ids: []string{"SjyX", "f8c9", "iJhz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			logger := logging.LogrusLogger()
			ctx := context.Background()
			hotelService := NewHotelService(logger, tc.httpClient(), ctx)

			fetchedSources, err := hotelService.UpdateHotelsFromSuppliers(tc.suppliersFilePath(), tc.hotelDataFilePath())
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				newData, err := hotelService.GetHotels(tc.hotelDataFilePath(), tc.ids, tc.destinations)
				assert.Nil(t, err)
				assert.ElementsMatch(t, tc.expectedData, newData)
				assert.ElementsMatch(t, tc.expectedFetchedSources, fetchedSources)
			}
		})
	}
}
