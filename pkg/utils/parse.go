package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pantho/golang/internal/conv/hash"
	errs "pantho/golang/pkg"
)

const (
	headerUserID   = "pathao-user-id"
	headerUserType = "pathao-user-type"
)

// headUsrID parses the user ID from the 'pathao-user-id' request header.
// It returns an error if the header is missing or contains an invalid value.
func headUsrID(header http.Header) (int32, error) {
	val := header.Get(headerUserID)
	if val == "" {
		return 0, errs.NewBadReq("header 'pathao-user-id' is missing", "user id is required")
	}

	userID, err := strconv.Atoi(val)
	if err != nil {
		logMsg := fmt.Sprintf("invalid user id in header: %s", val)
		return 0, errs.NewBadReq(logMsg, "invalid user id")
	}
	return int32(userID), nil
}

func paramUserID(c *gin.Context) (int32, error) {
	userId := c.Param("userId")
	if userId == "" {
		return 0, errs.NewBadReq("user id is required", "user id is required")
	}
	usrId, err := strconv.Atoi(userId)
	if err != nil {
		return 0, errs.NewBadReq(err.Error(), "invalid user id")
	}
	return int32(usrId), nil
}

func paramSlug(c *gin.Context) (string, error) {
	slug := c.Param("slug")
	if slug == "" {
		return "", errs.NewBadReq("slug is required", "slug is required")
	}
	return slug, nil
}

func paramParcelID(c *gin.Context) (uint64, error) {
	strParcelId := c.Param("parcelId")
	if strParcelId == "" {
		return 0, errs.NewBadReq("parcel id is required", "parcel id is required")
	}

	var parcelId uint64
	if id, err := strconv.ParseUint(strParcelId, 10, 64); err == nil {
		parcelId = id
	} else {
		parcelId = hash.Decode(strParcelId)
	}

	if parcelId == 0 {
		return 0, errs.NewBadReq("parcel id must be greater than 0", "invalid parcel id")
	}
	return parcelId, nil
}

// headUsrType parses the user type from the 'pathao-user-type' request header.
// It returns an error if the header is missing or contains an unsupported value.
func headUsrType(header http.Header) (string, error) {
	userType := header.Get(headerUserType)
	switch userType {
	case "user", "driver":
		return userType, nil
	default:
		logMsg := fmt.Sprintf("invalid user type in header: %s", userType)
		return "", errs.NewBadReq(logMsg, "invalid user type")
	}
}

func headCityID(header http.Header) (int32, error) {
	id, err := strconv.Atoi(header.Get("City-Id"))
	if err != nil {
		return 0, errs.NewBadReq(err.Error(), "invalid City-Id in the header")
	}
	return int32(id), err
}

func headCountryID(header http.Header) (int32, error) {
	id, err := strconv.Atoi(header.Get("Country-Id"))
	if err != nil {
		return 0, errs.NewBadReq(err.Error(), "invalid Country-Id in the header")
	}
	return int32(id), err
}

// parseLang parses the language from the 'lang' query parameter, defaulting to "en".
func parseLang(c *gin.Context) string {
	lang := c.Query("lang")
	switch lang {
	case "bn", "ne", "en":
		return lang
	default:
		return "en"
	}
}

func parseLocale(c *gin.Context) string {
	lang := c.Query("localization")
	switch lang {
	case "bn", "ne", "en":
		return lang
	default:
		return "en"
	}
}
