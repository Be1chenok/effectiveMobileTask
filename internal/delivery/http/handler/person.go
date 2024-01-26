package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Be1chenok/effectiveMobileTask/internal/domain"
	"github.com/gorilla/mux"
)

const (
	defaultPage = 1
	defaultSize = 5
)

func (h Handler) FindPersons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	urlParams := getUrlParams(r)
	searchParams := domain.PersonSearchParams{
		Gender:      urlParams.Gender,
		Nationality: urlParams.Nationality,
		Offset:      (urlParams.Page - 1) * urlParams.Size,
		Limit:       urlParams.Size,
	}

	persons, err := h.service.Person.Find(ctx, &searchParams)
	if err != nil {
		if errors.Is(err, domain.ErrNothingFound) {
			writeJsonErrorResponse(w, http.StatusBadRequest, domain.ErrNothingFound)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	writeJsonResponse(w, http.StatusOK, persons)
}

func (h Handler) FindPersonById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	personId, err := strconv.Atoi(id)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrIdIsNotNumber)
		return
	}

	person, err := h.service.FindById(ctx, personId)
	if err != nil {
		if errors.Is(err, domain.ErrNothingFound) {
			writeJsonErrorResponse(w, http.StatusBadRequest, domain.ErrNothingFound)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	writeJsonResponse(w, http.StatusOK, person)
}

func (h Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	defer r.Body.Close()

	var input FullName
	if err := json.Unmarshal(body, &input); err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if input.Name == "" || input.Surname == "" {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	person, err := h.service.Person.Add(ctx, &domain.Person{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	})
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	writeJsonResponse(w, http.StatusOK, person)
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	personId, err := strconv.Atoi(id)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrIdIsNotNumber)
		return
	}

	if err := h.service.Person.DeleteById(ctx, personId); err != nil {
		if errors.Is(err, domain.ErrNothingWasDeleted) {
			writeJsonErrorResponse(w, http.StatusBadRequest, domain.ErrNothingWasDeleted)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	writeJsonResponse(w, http.StatusNoContent, nil)
}

func (h Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	defer r.Body.Close()

	var person domain.Person
	if err := json.Unmarshal(body, &person); err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	if person.ID <= 0 ||
		person.Name == "" ||
		person.Surname == "" ||
		person.Age <= 0 ||
		person.Gender == "" ||
		person.Nationality == "" {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.service.Person.UpdateById(ctx, &person); err != nil {
		if errors.Is(err, domain.ErrNothingUpdated) {
			writeJsonErrorResponse(w, http.StatusBadRequest, domain.ErrNothingUpdated)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	writeJsonResponse(w, http.StatusNoContent, nil)
}

func getUrlParams(r *http.Request) UrlParams {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil || size <= 0 {
		size = defaultSize
	}

	return UrlParams{
		Gender:      r.URL.Query().Get("gender"),
		Nationality: r.URL.Query().Get("nationality"),
		Page:        page,
		Size:        size,
	}
}
