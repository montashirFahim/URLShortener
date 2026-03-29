package v1

// func RatingHandler(container *cache.Container) http.Handler {
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		userID := URLParamParser("id", r)
// 		userType := cnst.User(r.URL.Query().Get("type"))

// 		if userType == "" {
// 			userType = cnst.UserUser
// 		}

// 		rating := action.Rating{}
// 		err := coach.Captain(
// 			rating.ValidReqID(userID),
// 			//rating.ValidReqUserType(userType),
// 			rating.FetchRating(userID, userType, container),
// 		).Play()

// 		if err != nil {
// 			ServeErr(w, err)
// 			return
// 		}

// 		Serve(w, http.StatusOK, rating.CachedData.ToRating())
// 	}

// 	return http.HandlerFunc(handler)
// }
