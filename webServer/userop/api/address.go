package api

import (
	"context"
	"net/http"
	"strconv"
	"webServer/userop/forms"
	"webServer/userop/global"
	"webServer/userop/proto"

	"github.com/gin-gonic/gin"
)

type address struct{}

func NewAddress() *address {
	return &address{}
}

func (a *address) List(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	rsp, err := global.GrpcAddressClient.QueryAddressList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.AddressRequest{
			UserId: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}

func (a *address) Add(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var addressDetailForm forms.AddressDetailForm
	if errs := c.ShouldBindJSON(&addressDetailForm); errs != nil {
		printValidateErrorTips(c, errs)
		return
	}

	IsDefault := false
	if addressDetailForm.IsDefault != nil {
		IsDefault = *addressDetailForm.IsDefault
	}
	_, err = global.GrpcAddressClient.AddAddress(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.AddressRequest{
			UserId:       int32(claims.ID),
			Province:     addressDetailForm.Province,
			City:         addressDetailForm.City,
			District:     addressDetailForm.District,
			Address:      addressDetailForm.Address,
			SignerName:   addressDetailForm.SignerName,
			SignerMobile: addressDetailForm.SignerMobile,
			IsDefault:    IsDefault,
		},
	)

	if err != nil {
		printGrpcError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func (a *address) Update(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))

	var addressDetailForm forms.AddressDetailForm
	if errs := c.ShouldBindJSON(&addressDetailForm); errs != nil {
		printValidateErrorTips(c, errs)
		return
	}

	args := &proto.AddressRequest{
		Id:           int32(id),
		UserId:       int32(claims.ID),
		Province:     addressDetailForm.Province,
		City:         addressDetailForm.City,
		District:     addressDetailForm.District,
		Address:      addressDetailForm.Address,
		SignerName:   addressDetailForm.SignerName,
		SignerMobile: addressDetailForm.SignerMobile,
	}
	if addressDetailForm.IsDefault != nil {
		args.IsDefault = *addressDetailForm.IsDefault
	}

	_, err = global.GrpcAddressClient.UpdateAddress(
		context.WithValue(context.Background(), "ginContext", c),
		args,
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (a *address) Delete(c *gin.Context) {
	var addressIdForm forms.AddressIdForm
	if errs := c.ShouldBindUri(&addressIdForm); errs != nil {
		printValidateErrorTips(c, errs)
		return
	}

	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	_, err = global.GrpcAddressClient.DeleteAddress(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.AddressRequest{
			Id:     int32(addressIdForm.Id),
			UserId: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
