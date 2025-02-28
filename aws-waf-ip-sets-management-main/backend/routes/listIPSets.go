package routes

import (
    "net/http"
    "aws-waf-ip-sets-management/backend/config"
    "aws-waf-ip-sets-management/backend/utils"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/wafv2"
    "log"
)

func ListIPSets(w http.ResponseWriter, r *http.Request) {
    params := &wafv2.ListIPSetsInput{
        Scope: aws.String("REGIONAL"),
    }

    result, err := config.WAFv2.ListIPSets(params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error listing IP sets:", err)
        return
    }

    var ipSetsDetails []map[string]interface{}
    for _, ipSet := range result.IPSets {
        ipSetData, err := config.WAFv2.GetIPSet(&wafv2.GetIPSetInput{
            Id:    ipSet.Id,
            Name:  ipSet.Name,
            Scope: aws.String("REGIONAL"),
        })
        if err != nil {
            utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
            log.Println("Error getting IP set details:", err)
            return
        }

        ipSetsDetails = append(ipSetsDetails, map[string]interface{}{
            "Name":         aws.StringValue(ipSetData.IPSet.Name),
            "Id":           aws.StringValue(ipSetData.IPSet.Id),
            "Addresses":    aws.StringValueSlice(ipSetData.IPSet.Addresses),
            "AddressCount": len(ipSetData.IPSet.Addresses),
            "LockToken":    aws.StringValue(ipSetData.LockToken),
        })
    }

    utils.RespondWithJSON(w, http.StatusOK, ipSetsDetails)
    log.Println("Successfully listed IP sets")
}
