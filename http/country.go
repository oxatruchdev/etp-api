package http

func (s *Server) registerCountryRoutes() {
}

// func (s *Server) apiGetCountries(w http.ResponseWriter, r *http.Request) error {
// 	// Parsing filter parameters from request (use a helper method if needed)
// 	var filter etp.CountryFilter
// 	if err := Bind(r, &filter); err != nil {
// 		return Error(w, etp.Errorf(etp.EINVALID, "invalid body"))
// 	}
//
// 	// Fetch countries from the service layer
// 	countries, count, err := s.CountryService.GetCountries(r.Context(), filter)
// 	if err != nil {
// 		return Error(w, err)
// 	}
//
// 	// Return JSON response
// 	return JSON(w, http.StatusOK, map[string]interface{}{
// 		"countries": countries,
// 		"count":     count,
// 	})
// }
//
// func (s *Server) apiGetCountry(w http.ResponseWriter, r *http.Request) error {
// 	// Parse country ID from URL path
// 	idStr := Param(r, "id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return Error(w, etp.Errorf(etp.EINVALID, "invalid id"))
// 	}
//
// 	// Fetch the country by ID from the service layer
// 	country, err := s.CountryService.GetCountryById(r.Context(), id)
// 	if err != nil {
// 		return Error(w, err)
// 	}
//
// 	// Return country information in JSON
// 	return JSON(w, http.StatusOK, country)
// }
//
// func (s *Server) apiCreateCountry(w http.ResponseWriter, r *http.Request) error {
// 	// Parse the new country data from request body
// 	var country etp.Country
// 	if err := Bind(r, &country); err != nil {
// 		return Error(w, etp.Errorf(etp.EINVALID, "invalid body"))
// 	}
//
// 	// Call the service layer to create the new country
// 	if err := s.CountryService.CreateCountry(r.Context(), &country); err != nil {
// 		return Error(w, err)
// 	}
//
// 	// Return 201 Created status without a body
// 	return NoContent(w, http.StatusCreated)
// }
//
// func (s *Server) apiUpdateCountry(w http.ResponseWriter, r *http.Request) {
// 	// Parse country ID from URL path
// 	idStr := Param(r, "id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return Error(w, etp.Errorf(etp.EINVALID, "invalid id"))
// 	}
//
// 	// Parse the update data from request body
// 	var updCountry etp.CountryUpdate
// 	if err := Bind(r, &updCountry); err != nil {
// 		Error(w, r, etp.Errorf(etp.EINVALID, "invalid body"))
// 		return
// 	}
//
// 	// Call the service layer to update the country
// 	updatedCountry, err := s.CountryService.UpdateCountry(r.Context(), id, &updCountry)
// 	if err != nil {
// 		return Error(w, err)
// 	}
//
// 	// Return the updated country information in JSON
// 	return JSON(w, http.StatusOK, updatedCountry)
// }
