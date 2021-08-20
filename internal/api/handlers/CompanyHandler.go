package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yosa12978/halo/internal/pkg/dto"
	cache "github.com/yosa12978/halo/internal/pkg/redis"
	"github.com/yosa12978/halo/internal/pkg/repositories"
	"github.com/yosa12978/halo/pkg/helpers"
)

type ICompanyHandler interface {
	CreateCompany(w http.ResponseWriter, r *http.Request)
	GetCompanies(w http.ResponseWriter, r *http.Request)
	GetCompany(w http.ResponseWriter, r *http.Request)
	SearchCompany(w http.ResponseWriter, r *http.Request)
	UpdateCompany(w http.ResponseWriter, r *http.Request)
	DeleteCompany(w http.ResponseWriter, r *http.Request)
}

type CompanyHandler struct{}

func NewCompanyHandler() ICompanyHandler {
	return &CompanyHandler{}
}

func (ch *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company dto.CreateCompany
	json.NewDecoder(r.Body).Decode(&company)

	if err := validator.New().Struct(company); err != nil {
		helpers.RespondStatusCode(w, 400, err.Error())
		return
	}
	compRepo := repositories.NewCompanyRepository()
	err := compRepo.CreateCompany(company)
	if err != nil {
		helpers.RespondStatusCode(w, 400, err.Error())
		return
	}
	helpers.RespondStatusCode(w, 201, "Created")
}

func (ch *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	rcache := cache.GetClient()
	companiesC, err := rcache.Get("companies").Result()
	if err == nil {
		helpers.RespondJson(w, 200, companiesC)
		return
	}
	compRepo := repositories.NewCompanyRepository()
	companies := compRepo.GetCompanies()
	rcache.Set("companies", companies, 1*time.Minute)
	helpers.RespondJson(w, 200, companies)
}

func (ch *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	rcache := cache.GetClient()
	companyC, err := rcache.Get(fmt.Sprintf("company_%s", id)).Result()
	if err == nil {
		helpers.RespondJson(w, 200, companyC)
		return
	}
	compRepo := repositories.NewCompanyRepository()
	company, err := compRepo.GetCompany(id)
	if err != nil {
		helpers.RespondStatusCode(w, 404, err.Error())
		return
	}
	rcache.Set(fmt.Sprintf("company_%s", id), company, 0)
	helpers.RespondJson(w, 200, company)
}

func (ch *CompanyHandler) SearchCompany(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	compRepo := repositories.NewCompanyRepository()
	company := compRepo.SearchCompany(query)
	helpers.RespondJson(w, 200, company)
}

func (ch *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var ucompany dto.UpdateCompany
	json.NewDecoder(r.Body).Decode(&ucompany)
}

func (ch *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {

}
