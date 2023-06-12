package middlewares

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/JohnSalazar/microservices-go-common/httputil"
	"github.com/JohnSalazar/microservices-go-common/security"

	"github.com/JohnSalazar/microservices-go-common/helpers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Authentication struct {
	logger        *logrus.Logger
	managerTokens *security.ManagerTokens
}

func NewAuthentication(
	logger *logrus.Logger,
	managerTokens *security.ManagerTokens,
) *Authentication {
	return &Authentication{
		logger:        logger,
		managerTokens: managerTokens,
	}
}

func (auth *Authentication) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.managerTokens.ReadHeadAccessToken(c)
		if err != nil {
			auth.logger.Error(err.Error())
			httputil.NewResponseAbort(c, http.StatusUnauthorized, "token is not valid")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"invalid cookie"},
			// })
			return
		}

		id := claims.Sub
		if !helpers.IsValidID(id) {
			auth.logger.Error("ID is not valid")
			httputil.NewResponseAbort(c, http.StatusUnauthorized, "ID is not valid")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"ID is not valid"},
			// })
			return
		}
		c.Set("user", id)

		var claimsList []interface{}
		data, _ := json.Marshal(claims.Claims)
		json.Unmarshal(data, &claimsList)

		claimsList = auth.sortClaims(claimsList)

		c.Set("claims", claimsList)

		c.Next()
	}
}

func (auth *Authentication) sortClaims(claimsList []interface{}) []interface{} {
	for key, value := range claimsList {
		values := value.(map[string]interface{})
		list := strings.Split(values["value"].(string), ",")
		sort.Slice(list, func(i, j int) bool {
			return list[i] < list[j]
		})

		values["value"] = strings.Join(list, ",")
		claimsList[key] = values
	}

	return claimsList
}
