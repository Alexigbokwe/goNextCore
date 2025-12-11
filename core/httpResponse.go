package core

import "log"

type HttpResponseType[T interface{}] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    T      `json:"data,omitempty"`
}

func HttpSuccess(message string, code int) HttpResponseType[interface{}] {
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  true,
	}
}

func HttpSuccessWithData(message string, code int, data interface{}) HttpResponseType[interface{}] {
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  true,
		Data:    data,
	}
}

func HttpError(message string, code int) HttpResponseType[interface{}] {
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  false,
	}
}

func HttpErrorWithData(message string, code int, data interface{}) HttpResponseType[interface{}] {
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  false,
		Data:    data,
	}
}

func HttpErrorWithLog(message string, code int, err error) HttpResponseType[interface{}] {
	if err != nil {
		log.Printf("Error: %s: %v", message, err)
	}
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  false,
	}
}

func HttpErrorWithDataAndLog(message string, code int, data interface{}, err error) HttpResponseType[interface{}] {
	if err != nil {
		log.Printf("Error: %s: %v", message, err)
	}
	return HttpResponseType[interface{}]{
		Code:    code,
		Message: message,
		Status:  false,
		Data:    data,
	}
}
