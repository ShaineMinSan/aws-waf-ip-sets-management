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

type CreateIPSetRequest struct {
    Name        string   `json:"name"`
    Addresses   []string `json:"addresses"`
    Description string   `json:"description"`
}

func CreateIPSet(w http.ResponseWriter, r *http.Request) {
    var req CreateIPSetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        log.Println("Error decoding request body:", err)
        return
    }

    params := &wafv2.CreateIPSetInput{
        Name:          aws.String(req.Name),
        Scope:         aws.String("REGIONAL"),
        Addresses:     aws.StringSlice(req.Addresses),
        IPAddressVersion: aws.String("IPV4"),
        Description:   aws.String(req.Description),
    }

    resp, err := config.WAFv2.CreateIPSet(params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error creating IP set:", err)
        return
    }

    addressesJSON, _ := json.Marshal(req.Addresses)
    _, err = config.DB.Exec("INSERT INTO actions (name, type, action) VALUES (?, ?, ?)", req.Name, "create", string(addressesJSON))
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error inserting action into database:", err)
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, resp)
    log.Println("Successfully created IP set:", req.Name)
}
