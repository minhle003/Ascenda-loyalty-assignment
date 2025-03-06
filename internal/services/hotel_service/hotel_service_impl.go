package hotel_service

import (
	"ascenda-loyalty-assignment/pkg/logging"
	"ascenda-loyalty-assignment/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type hotelServiceImpl struct {
	logger     logging.Logger
	httpClient HTTPClient
	ctx        context.Context
}

func NewHotelService(logger logging.Logger, httpClient HTTPClient, ctx context.Context) HotelService {
	return &hotelServiceImpl{
		logger:     logger,
		httpClient: httpClient,
		ctx:        ctx,
	}
}

func (h *hotelServiceImpl) GetHotels(hotelDataFilePath string, ids []string, destinations []int) ([]Hotel, error) {
	if len(ids) == 0 && len(destinations) == 0 {
		return []Hotel{}, nil
	}

	hotels, err := h.getHotelDataFromDataFile(hotelDataFilePath)
	if err != nil {
		return []Hotel{}, err
	}

	filteredHotels := make([]Hotel, 0)
	addedHotelIds := make(map[string]bool)
	h.filterHotelIds(hotels, &filteredHotels, addedHotelIds, ids)
	h.filterDestinationIds(hotels, &filteredHotels, addedHotelIds, destinations)

	return filteredHotels, nil
}

func (h *hotelServiceImpl) UpdateHotelsFromSuppliers(suppliersFilePath string, hotelDataFilePath string) ([]string, error) {
	currentHotelData, err := h.getHotelDataFromDataFile(hotelDataFilePath)
	if err != nil {
		return []string{}, err
	}

	isFileEmpty, err := utils.IsFileEmpty(suppliersFilePath)
	if isFileEmpty || err != nil {
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to read file %s", suppliersFilePath), err)
		} else {
			h.logger.Error("There is no suppliers in data file")
		}
		return []string{}, fmt.Errorf("failed to get suppliers data")
	}
	data, err := utils.ReadJSONFile(suppliersFilePath)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to read file %s", suppliersFilePath), err)
		return []string{}, fmt.Errorf("failed to get suppliers data")
	}
	suppliers, err := h.unmarshalSuppliers(data)
	if err != nil {
		h.logger.Error("fail to parse json suppliers data", err)
		return []string{}, fmt.Errorf("failed to get suppliers data")
	}

	if len(suppliers) == 0 {
		h.logger.Error("There is no suppliers in data file")
		return []string{}, fmt.Errorf("failed to get suppliers data")
	}

	hotelsDataFromSuppliers, fetchedDataSources, err := h.fetchDataFromSuppliers(suppliers)
	if err != nil {
		h.logger.Error("Fail to get data from data sources", err)
		return []string{}, fmt.Errorf("unable to update new hotel data")
	}
	h.sanitizeHotelData(hotelsDataFromSuppliers, currentHotelData)

	err = utils.WriteJSONFile(hotelDataFilePath, currentHotelData)
	if err != nil {
		h.logger.Error("Fail to write to hotel json data file", err)
		return []string{}, fmt.Errorf("unable to update new hotel data")
	}

	return fetchedDataSources, nil
}

func (h *hotelServiceImpl) getHotelDataFromDataFile(hotelDataFilePath string) (map[string]Hotel, error) {
	isFileEmpty, err := utils.IsFileEmpty(hotelDataFilePath)
	if isFileEmpty || err != nil {
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to read file %s", hotelDataFilePath), err)
			return map[string]Hotel{}, fmt.Errorf("failed to get hotel data")
		}
		return map[string]Hotel{}, nil
	}

	data, err := utils.ReadJSONFile(hotelDataFilePath)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to read file %s", hotelDataFilePath), err)
		return map[string]Hotel{}, fmt.Errorf("failed to get hotel data")
	}

	hotels, err := h.unmarshalHotels(data)
	if err != nil {
		h.logger.Error("fail to parse json hotels data", err)
		return map[string]Hotel{}, fmt.Errorf("failed to get hotel data")
	}

	return hotels, nil
}

func (h *hotelServiceImpl) filterHotelIds(
	hotels map[string]Hotel,
	filteredHotels *[]Hotel,
	addedHotelIds map[string]bool,
	ids []string,
) {
	if len(ids) == 0 {
		return
	}
	idsMap := map[string]bool{}
	for _, id := range ids {
		idsMap[id] = true
	}
	for id, hotel := range hotels {
		if idsMap[id] && !addedHotelIds[id] {
			hotel.ID = id
			*filteredHotels = append(*filteredHotels, hotel)
			addedHotelIds[id] = true
		}
	}
}

func (h *hotelServiceImpl) filterDestinationIds(
	hotels map[string]Hotel,
	filteredHotels *[]Hotel,
	addedHotelIds map[string]bool,
	destinations []int,
) {
	if len(destinations) == 0 {
		return
	}
	destinationIdsMap := map[int]bool{}
	for _, id := range destinations {
		destinationIdsMap[id] = true
	}

	for id, hotel := range hotels {
		if destinationIdsMap[hotel.DestinationID] && !addedHotelIds[id] {
			hotel.ID = id
			*filteredHotels = append(*filteredHotels, hotel)
			addedHotelIds[id] = true
		}
	}
}

func (h *hotelServiceImpl) unmarshalHotels(data []byte) (map[string]Hotel, error) {
	var hotels map[string]Hotel
	err := json.Unmarshal(data, &hotels)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

func (h *hotelServiceImpl) unmarshalSuppliers(data []byte) ([]string, error) {
	var suppliers []string
	err := json.Unmarshal(data, &suppliers)
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (h *hotelServiceImpl) fetchDataFromSuppliers(suppliers []string) ([]map[string]interface{}, []string, error) {
	var fetchedHotelsData []map[string]interface{}
	var fetchedDataSources []string
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(suppliers))

	routineCtx, cancel := context.WithCancel(h.ctx)
	defer cancel()

	errChan := make(chan error, len(suppliers))

	for _, url := range suppliers {
		go func(url string) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(routineCtx, http.MethodGet, url, nil)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("error creating request for %s: %w", url, err):
				case <-routineCtx.Done():
					return
				}
				return
			}

			resp, err := h.httpClient.Do(req)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("error fetching data from %s: %w", url, err):
				case <-routineCtx.Done():
					return
				}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				select {
				case errChan <- fmt.Errorf("supplier %s returned status code: %d", url, resp.StatusCode):
				case <-routineCtx.Done():
					return
				}
				return
			}

			var hotels []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&hotels)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("error decoding JSON from %s: %w", url, err):
				case <-routineCtx.Done():
					return
				}
				return
			}

			mu.Lock()
			fetchedHotelsData = append(fetchedHotelsData, hotels...)
			fetchedDataSources = append(fetchedDataSources, url)
			h.logger.Info("Successfully fetching data from supplier ", url)
			mu.Unlock()
		}(url)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		h.logger.Error("Error occurred while fetching data from suppliers", strings.Join(errs, "\n"))
		if len(errs) >= len(suppliers) {
			return fetchedHotelsData, fetchedDataSources, fmt.Errorf(strings.Join(errs, "\n"))
		}
	}

	return fetchedHotelsData, fetchedDataSources, nil
}

func (h *hotelServiceImpl) sanitizeHotelData(updatedData []map[string]interface{}, currentHotelData map[string]Hotel) {
	for _, hotel := range updatedData {
		id := h.getHotelIdFromUpdatedData(hotel)
		if id == "" {
			continue
		}
		var newHotelData Hotel
		if _, exists := currentHotelData[id]; !exists {
			newHotelData = Hotel{
				ID: id,
			}
		} else {
			newHotelData = currentHotelData[id]
		}
		destinationId := h.getDestinationIdFromUpdatedData(hotel)
		if destinationId == -1 {
			continue
		}
		newHotelData.DestinationID = destinationId

		newHotelData.HotelName = h.getHotelNameFromUpdatedData(hotel)

		newLocation := h.getLocationFromUpdatedData(hotel)
		if newLocation.Lat != 0.0 {
			newHotelData.Location.Lat = newLocation.Lat
		}
		if newLocation.Long != 0.0 {
			newHotelData.Location.Long = newLocation.Long
		}
		if newLocation.Address != "" {
			newHotelData.Location.Address = newLocation.Address
		}
		if newLocation.City != "" {
			newHotelData.Location.City = newLocation.City
		}
		if newLocation.Country != "" {
			newHotelData.Location.Country = newLocation.Country
		}

		newDescription := h.getHotelDescriptionFromUpdatedData(hotel)
		if newDescription != "" && !utils.SliceContains(newHotelData.Description, newDescription) {
			newHotelData.Description = append(newHotelData.Description, newDescription)
		}

		newConditions := h.getHotelBookingConditionFromUpdatedData(hotel)
		for _, condition := range newConditions {
			if !utils.SliceContains(newHotelData.BookingCondition, condition) {
				newHotelData.BookingCondition = append(newHotelData.BookingCondition, condition)
			}
		}

		newAmenities := h.getHotelAmenitiesFromUpdatedData(hotel)
		for _, genAmenity := range newAmenities.General {
			if !utils.SliceContains(newHotelData.Amenities.General, genAmenity) {
				newHotelData.Amenities.General = append(newHotelData.Amenities.General, genAmenity)
			}
		}
		for _, roomAmenity := range newAmenities.Room {
			if !utils.SliceContains(newHotelData.Amenities.Room, roomAmenity) {
				newHotelData.Amenities.Room = append(newHotelData.Amenities.Room, roomAmenity)
			}
		}

		newImages := h.getHotelImagesFromUpdatedData(hotel)
		for imageCategory, images := range newImages {
			if newHotelData.Images == nil {
				newHotelData.Images = make(map[string][]Image)
			}
			curImages, exists := newHotelData.Images[imageCategory]
			if !exists {
				newHotelData.Images[imageCategory] = images
				continue
			}
			for _, image := range images {
				exist := false
				for _, curImage := range curImages {
					if curImage.Link == image.Link {
						exist = true
						break
					}
				}
				if !exist {
					newHotelData.Images[imageCategory] = append(newHotelData.Images[imageCategory], image)
				}
			}
		}

		currentHotelData[id] = newHotelData
	}
}

func (h *hotelServiceImpl) getHotelIdFromUpdatedData(hotel map[string]interface{}) string {
	idKeys := []string{"id", "Id", "hotel_id"} // add another key here  for future data sources
	var id interface{}

	for _, idKey := range idKeys {
		if val, ok := hotel[idKey]; ok {
			id = val
			break
		}
	}
	if id == nil {
		h.logger.Warn("Data is invalid, missing hotelId", hotel)
		return ""
	}

	return fmt.Sprintf("%v", id)
}

func (h *hotelServiceImpl) getDestinationIdFromUpdatedData(hotel map[string]interface{}) int {
	destinationIdKeys := []string{"destination_id", "DestinationId", "destination"} // add another key here for future data sources
	var destinationId interface{}

	for _, idKey := range destinationIdKeys {
		if val, ok := hotel[idKey]; ok {
			destinationId = val
			break
		}
	}
	if destinationId == nil {
		h.logger.Warn("Data is invalid, missing destination id", hotel)
		return -1
	}

	destinationIdFloat, ok := destinationId.(float64)
	if !ok {
		h.logger.Warn("Data is invalid, destination id is not in integer format", hotel)
		return -1
	}
	if destinationIdFloat != float64(int(destinationIdFloat)) {
		h.logger.Warn("Data is invalid, destination id is not in integer format", hotel)
	}
	return int(destinationIdFloat)
}

func (h *hotelServiceImpl) getHotelNameFromUpdatedData(hotel map[string]interface{}) string {
	nameKeys := []string{"Name", "name", "hotel_name"} // add another key here for future data sources
	var name interface{}

	for _, nameKey := range nameKeys {
		if val, ok := hotel[nameKey]; ok {
			name = val
			break
		}
	}

	return utils.ConvertInterfaceToString(name)
}

func (h *hotelServiceImpl) getLocationFromUpdatedData(hotel map[string]interface{}) Location {
	locationKeys := []string{"location", "Location"} // add another key here  for future data sources
	latKeys := []string{"lat", "Latitude"}           // add another key here  for future data sources
	lngKeys := []string{"Longitude", "lng"}          // add another key here  for future data sources
	cityKeys := []string{"City", "city"}             // add another key here  for future data sources
	addressKeys := []string{"Address", "address"}    // add another key here  for future data sources
	countryKeys := []string{"country", "Country"}    // add another key here  for future data sources

	hotelLocationDataContainer := hotel
	for _, locKey := range locationKeys {
		if loc, ok := hotel[locKey]; ok {
			hotelLocationDataContainer = loc.(map[string]interface{})
			break
		}
	}

	var newLocation Location

	for _, latKey := range latKeys {
		if val, ok := hotelLocationDataContainer[latKey]; ok {
			if latVal, ok := val.(float64); ok {
				newLocation.Lat = latVal
			} else {
				h.logger.Warn("Latitude data type not supported", val)
			}
			break
		}
	}

	for _, lngKey := range lngKeys {
		if val, ok := hotelLocationDataContainer[lngKey]; ok {
			if lngVal, ok := val.(float64); ok {
				newLocation.Long = lngVal
			} else {
				h.logger.Warn("Longitude data type not supported", val)
			}
			break
		}
	}

	for _, addKey := range addressKeys {
		if address, ok := hotelLocationDataContainer[addKey]; ok {
			newLocation.Address = utils.ConvertInterfaceToString(address)
			break
		}
	}

	for _, cityKey := range cityKeys {
		if city, ok := hotelLocationDataContainer[cityKey]; ok {
			newLocation.City = utils.ConvertInterfaceToString(city)
			break
		}
	}

	for _, countryKey := range countryKeys {
		if country, ok := hotelLocationDataContainer[countryKey]; ok {
			newLocation.Country = utils.ConvertInterfaceToString(country)
			break
		}
	}

	return newLocation

}

func (h *hotelServiceImpl) getHotelDescriptionFromUpdatedData(hotel map[string]interface{}) string {
	descKeys := []string{"description", "Description", "info", "details"} // add another key here  for future data sources
	var description interface{}

	for _, descKey := range descKeys {
		if val, ok := hotel[descKey]; ok {
			description = val
			break
		}
	}
	if description == nil {
		return ""
	}

	return utils.ConvertInterfaceToString(description)
}

func (h *hotelServiceImpl) getHotelBookingConditionFromUpdatedData(hotel map[string]interface{}) []string {
	condKeys := []string{"booking_conditions"} // add another key here  for future data sources
	var bookingConditions []string

	for _, condKey := range condKeys {
		if conditions, ok := hotel[condKey]; ok {
			if conditionsVal, ok := conditions.([]interface{}); ok {
				for _, conditions := range conditionsVal {
					bookingConditions = append(bookingConditions, utils.ConvertInterfaceToString(conditions))
				}
			} else {
				h.logger.Warn("Booking Condition data type not supported", conditions)
			}
			break
		}
	}
	return bookingConditions
}

func (h *hotelServiceImpl) getHotelAmenitiesFromUpdatedData(hotel map[string]interface{}) Amenities {
	amenityKey := []string{"amenities", "Facilities"} // add another key here  for future data sources
	var hotelAmenities Amenities

	for _, aKey := range amenityKey {
		if amenities, ok := hotel[aKey]; ok {
			if amenitiesVal, ok := amenities.([]interface{}); ok {
				for _, amenity := range amenitiesVal {
					hotelAmenities.General = append(hotelAmenities.General, utils.ConvertInterfaceToString(amenity))
				}
				continue
			} else {
				h.logger.Warn("Amenities data type not supported", amenities)
			}
			if amenitiesVal, ok := amenities.(map[string]interface{}); ok {
				if generalAmenities, exists := amenitiesVal["general"]; exists {
					if listGeneralAmenities, ok := generalAmenities.([]interface{}); ok {
						for _, amenity := range listGeneralAmenities {
							hotelAmenities.General = append(hotelAmenities.General, utils.ConvertInterfaceToString(amenity))
						}
					} else {
						h.logger.Warn("Amenities data type not supported", generalAmenities)
					}
				}
				if roomAmenities, exists := amenitiesVal["room"]; exists {
					if listRoomAmenities, ok := roomAmenities.([]interface{}); ok {
						for _, amenity := range listRoomAmenities {
							hotelAmenities.Room = append(hotelAmenities.Room, utils.ConvertInterfaceToString(amenity))
						}
					} else {
						h.logger.Warn("Amenities data type not supported", roomAmenities)
					}
				}
			} else {
				h.logger.Warn("Amenities data type not supported", amenities)
			}
			break
		}
	}

	return hotelAmenities
}

func (h *hotelServiceImpl) getHotelImagesFromUpdatedData(hotel map[string]interface{}) map[string][]Image {
	imageKeys := []string{"images", "Images"}           // add another key here for future data sources
	imageLinkKeys := []string{"url", "link"}            // add another key here for future data sources
	imageDescKeys := []string{"caption", "description"} // add another key here for future data sources

	hotelImages := make(map[string][]Image)
	for _, iKey := range imageKeys {
		if images, ok := hotel[iKey]; ok {
			if imagesVal, ok := images.(map[string]interface{}); ok {
				for imageType, imagesOfType := range imagesVal {
					hotelImages[imageType] = []Image{}
					if listImages, ok := imagesOfType.([]interface{}); ok {
						for _, image := range listImages {
							if imageMap, ok := image.(map[string]interface{}); ok {
								newImage := Image{}
								for _, imageLinkKey := range imageLinkKeys {
									if link, ok := imageMap[imageLinkKey]; ok {
										newImage.Link = utils.ConvertInterfaceToString(link)
										break
									}
								}
								for _, imageDescKey := range imageDescKeys {
									if desc, ok := imageMap[imageDescKey]; ok {
										newImage.Description = utils.ConvertInterfaceToString(desc)
										break
									}
								}
								hotelImages[imageType] = append(hotelImages[imageType], newImage)
							} else {
								h.logger.Warn("Images data type not supported", image)
							}
						}
					} else {
						h.logger.Warn("Images data type not supported", imagesOfType)
					}
				}
			} else {
				h.logger.Warn("Images data type not supported", images)
			}
			break
		}
	}
	return hotelImages
}
