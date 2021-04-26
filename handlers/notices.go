package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jamolh/notice-board/db"
	"github.com/jamolh/notice-board/helpers"
	"github.com/jamolh/notice-board/models"
	"github.com/julienschmidt/httprouter"
)

// @Summary Create new notice based on parametrs
// @Description Method to create new notice
// @Accept json
// @Produce json
// @Param CreateNoticeRequest body models.Notice true "Create Notice"
// @Success 200 {string} models.Response{payload=id} "Success"
// @Failure 208 {object} models.ErrorResponse "Already exists"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Not found"
// @Failure 500 {object} models.ErrorResponse "Server internal error"
// @Router /v1/notices [POST]
func CreateNoticeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		response = models.Response{
			ReqID: uuid.NewString(),
			Code:  200,
		}
		request models.Notice
		err     error
	)
	defer response.Send(w, r)

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Code = http.StatusBadRequest
		log.Println("CreateNoticeHandler parcing request failed:", err)
		return
	}

	if err = validateCreateNoticeRequest(request); err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return
	}

	exists, err := db.CheckNoticeExistsByTitle(r.Context(), request.Title)
	switch {
	case err != nil:
		response.Code = http.StatusInternalServerError
		log.Println("CreateNoticeHandler, db CheckNoticeExistsByTitle error:", err)
		return
	case exists:
		response.Code = http.StatusAlreadyReported
		response.Message = "Notice by this title exists"
		return
	}

	err = db.CreateNotice(r.Context(), &request)
	if err != nil {
		response.Code = http.StatusInternalServerError
		log.Println("CreateNoticeHandler, db CreateNotice error:", err)
		return
	}

	response.Payload = request.ID
}

// @Summary Method for getting notice by ID
// @Description Method for getting a specific notice by id
// @Accept json
// @Produce json
// @Param id path string true "Notice ID"
// @Param fields query string false "Get notice with all images"
// @Success 200 {object} models.GetNoticeRequestDto "Success"
// @Failure 208 {object} models.ErrorResponse "Already exists"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Not found"
// @Failure 500 {object} models.ErrorResponse "Server internal error"
// @Router /v1/notices/{id} [GET]
func GetNoticeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		request = &models.GetNoticeRequestDto{
			ID: p.ByName("id"),
		}
		response = models.Response{
			ReqID: uuid.NewString(),
			Code:  200,
		}
	)

	defer response.Send(w, r)

	if !helpers.IsValidUUID(request.ID) {
		response.Code = http.StatusBadRequest
		return
	}

	if filter := r.URL.Query().Get("fields"); filter == "" {
		request.GetAllImages = false
	}

	notice, err := db.GetNoticeByID(r.Context(), request)
	if err != nil {
		log.Println("GetNoticeHandler, db GetNotice error:", err)

		if db.IsNotFound(err) {
			response.Code = http.StatusNotFound
			response.Message = "Notice with this id not exists"
			return
		}
		response.Code = http.StatusInternalServerError
		return
	}

	response.Payload = notice
}

// @Summary Method to take all notices
// @Description get all notices
// @Accept json
// @Produce json
// @Param sort_field query string false "Sort by field created_at or price" Enums(created_at, price)
// @Param sort_type query string false "Sort by ascending or descending" Enums(asc, desc)
// @Success 200 {array} models.Notice "Success"
// @Failure 208 {object} models.ErrorResponse "Already exists"
// @Failure 500 {object} models.ErrorResponse "Server internal error"
// @Router /v1/notices/ [GET]
func GetNoticesHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		response = models.Response{
			ReqID: uuid.NewString(),
			Code:  200,
		}
	)

	defer response.Send(w, r)

	request := getNoticesFilter(r)
	notices, err := db.GetNotices(r.Context(), request)
	if err != nil {
		log.Println("GetNoticesHandler, db GetNotices error:", err)

		if db.IsNotFound(err) {
			response.Code = http.StatusNotFound
			response.Message = "Notices not found"
		}
		response.Code = http.StatusInternalServerError
		return
	}
	response.Payload = notices
}

func validateCreateNoticeRequest(request models.Notice) error {
	if helpers.RemoveNonLetter(request.Title) == "" {
		return errors.New("Wrong title for notice")
	}

	if len(helpers.RemoveNonLetter(request.Description)) < 5 {
		return errors.New("Description is too short")
	}

	return nil
}

func getNoticesFilter(r *http.Request) models.GetNoticesRequestDto {
	return models.GetNoticesRequestDto{
		Field: r.URL.Query().Get("sort_field"),
		Order: r.URL.Query().Get("sort_type"),
	}
}
