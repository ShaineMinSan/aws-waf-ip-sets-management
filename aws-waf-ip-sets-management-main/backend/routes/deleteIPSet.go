package routes

import (
    "encoding/json"
    "net/http"
    "aws-waf-ip-sets-management/backend/config"
    "aws-waf-ip-sets-management/backend/utils"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/wafv2"
    "log"
)

type DeleteIPSetRequest struct {
    Id       string `json:"id"`
    Name     string `json:"name"`
    LockToken string `json:"lockToken"`
}

func DeleteIPSet(w http.ResponseWriter, r *http.Request) {
    var req DeleteIPSetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        log.Println("Error decoding request body:", err)
        return
    }

    ipSet, err := config.WAFv2.GetIPSet(&wafv2.GetIPSetInput{
        Id:    aws.String(req.Id),
        Name:  aws.String(req.Name),
        Scope: aws.String("REGIONAL"),
    })
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error getting IP set:", err)
        return
    }

    params := &wafv2.DeleteIPSetInput{
        Id:       aws.String(req.Id),
        Name:     aws.String(req.Name),
        Scope:    aws.String("REGIONAL"),
        LockToken: aws.String(req.LockToken),
    }

    resp, err := config.WAFv2.DeleteIPSet(params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error deleting IP set:", err)
        return
    }

    addressesJSON, _ := json.Marshal(ipSet.IPSet.Addresses)
    _, err = config.DB.Exec("INSERT INTO actions (name, type, action) VALUES (?, ?, ?)", req.Name, "delete", string(addressesJSON))
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error inserting action into database:", err)
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, resp)
    log.Println("Successfully deleted IP set:", req.Name)
}
