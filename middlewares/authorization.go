package middlewares

import (
	"net/http"
	"strings"

	"github.com/oceano-dev/microservices-go-common/httputil"

	"github.com/gin-gonic/gin"
)

func Authorization(claimName string, claimValue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		getClaims, permissionOk := c.Get("claims")
		if permissionOk {
			claims := getClaims.([]interface{})
			permissionOk = validateClaims(claims, claimName, claimValue)
			// permissionOk = verifyClaimsPermission(claims, claimName, claimValue)
		}

		if !permissionOk {
			httputil.NewResponseAbort(c, http.StatusForbidden, "you do not have permission")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"you do not have permission"},
			// })
			return
		}

		c.Next()
	}
}

func validateClaims(userClaims []interface{}, claimType string, claimValue string) bool {
	_claimType := strings.TrimSpace(claimType)
	_claimValue := strings.TrimSpace(claimValue)
	if len(_claimType) == 0 || len(_claimValue) == 0 {
		return false
	}

	_claimValueSplitted := strings.Split(_claimValue, ",")

	for _, interator := range userClaims {
		result := interator.(map[string]interface{})
		if result["type"] == _claimType {
			userValueClaimSplitted := strings.Split(result["value"].(string), ",")

			return arrayContainsArray(_claimValueSplitted, userValueClaimSplitted)
		}
	}

	return false
}

func arrayContainsArray(array1 []string, array2 []string) bool {
	for _, val1 := range array1 {
		found := false
		for _, val2 := range array2 {
			if val1 == val2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// func verifyClaimsPermission(claims []interface{}, claimType string, claimValue string) bool {
// 	sClaimType := strings.TrimSpace(claimType)
// 	sClaimValue := strings.TrimSpace(claimValue)
// 	if len(sClaimType) == 0 || len(sClaimValue) == 0 {
// 		return false
// 	}
// 	for _, interator := range claims {
// 		values := interator.(map[string]interface{})
// 		if values["type"] == sClaimType && strings.Contains(values["value"].(string), sortClaimValue(sClaimValue)) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func sortClaimValue(sClaimsValue string) string {
// 	list := strings.Split(sClaimsValue, ",")

// 	sort.Slice(list, func(i, j int) bool {
// 		return list[i] < list[j]
// 	})

// 	claimsList := strings.Join(list, ",")

// 	return claimsList
// }
