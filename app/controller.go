package app

import (
	"ProductService/models"
	"ProductService/services"
	"ProductService/utils"
	enum "ProductService/utils/enums"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
)

type Controller interface {
	HandleProduct(w http.ResponseWriter, r *http.Request)
}

type ProductController struct {
	Proc       services.ProductMsgProc
	HttpClient *http.Client
}

func ProductHandler(p services.ProductMsgProc) *ProductController {
	httpCLient, e := utils.GetHttpClient()
	if e != nil {
		log.Println("Error in creating http client", e)
		return nil
	}
	return &ProductController{
		Proc:       p,
		HttpClient: httpCLient,
	}
}

func (c *ProductController) HandleProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entered HandleProduct")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Printf("END POINT: %v", r.RequestURI)

	// Getting details from request body
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in reading body", err)
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: err.Error(),
			ResponseBody:        nil,
		}
		data, _ := json.Marshal(msg)
		w.Write(data)
		return
	}

	//if it is not a GET request
	if r.Method == http.MethodPost && len(jsonData) != 0 {
		validate := validator.New()
		res := validate.Var(string(jsonData), "json")
		if res != nil {
			log.Println("Error in parsing body", res)
			msg := models.Result{
				ResponseCode:        enum.FailureCode400,
				ResponseStatus:      enum.FailureMessage400,
				ResponseDescription: res.Error(),
				ResponseBody:        nil,
			}
			data, _ := json.Marshal(msg)
			w.Write(data)
			return
		}
	}

	if r.Method == http.MethodPut && len(jsonData) == 0 {
		err := errors.New("unable to fetch details from req body")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: err.Error(),
			ResponseBody:        nil,
		}
		data, _ := json.Marshal(msg)
		w.Write(data)
		return
	}

	log.Println("Request received: ", string(jsonData))
	format, err := c.Proc.Decode(jsonData)
	if err != nil {
		log.Println("Json data decode failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: err.Error(),
			ResponseBody:        nil,
		}
		data, _ := json.Marshal(msg)
		w.Write(data)
		return
	}

	_, err = json.Marshal(format)
	if err != nil {
		log.Println("Json marshal of request body failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: err.Error(),
			ResponseBody:        nil,
		}
		data, _ := json.Marshal(msg)
		w.Write(data)
		return
	}

	e := c.Proc.Validate(format)
	if e != nil {
		log.Println("Json validation failed error in json structure, fields missing")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: e.Error(),
			ResponseBody:        nil,
		}
		data, _ := json.Marshal(msg)
		w.Write(data)
		return
	}

	msg, err := c.Proc.ProcessMsg(format, r)
	if err != nil {
		log.Println("Error in ProcessMsg", err)
		//w.WriteHeader(http.StatusInternalServerError)
		data, statusCode, er := c.Proc.Encode(msg)
		if er != nil {
			log.Println("Error in Encode", er)
			msg = models.Result{
				ResponseCode:        enum.FailureCode400,
				ResponseStatus:      enum.FailureMessage400,
				ResponseDescription: err.Error(),
				ResponseBody:        nil,
			}
			data, _ = json.Marshal(msg)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
		}
		w.WriteHeader(statusCode)
		w.Write(data)
		return
	}
	//w.WriteHeader(http.StatusOK)
	data, statusCode, err := c.Proc.Encode(msg)
	if err != nil {
		log.Println("Error in Encode", err)
		msg = models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: err.Error(),
			ResponseBody:        nil,
		}
		w.Write(data)
	}
	log.Printf("Response is: %v", string(data))
	log.Printf("End Handle HandleProduct")
	w.WriteHeader(statusCode)
	w.Write(data)
}
