package routes

import (
    "encoding/json"
    "net/http"
    "aws-waf-ip-sets-management/backend/config"
    "aws-waf-ip-sets-management/backend/utils"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/wafv2"
    "log"
    "strings"
)

type RemoveIPAddressRequest struct {
    Id       string   `json:"id"`
    Name     string   `json:"name"`
    Addresses []string `json:"addresses"`
    LockToken string  `json:"lockToken"`
}

func RemoveIPAddress(w http.ResponseWriter, r *http.Request) {
    var req RemoveIPAddressRequest
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

    existingAddresses := make(map[string]bool)
    for _, addr := range ipSet.IPSet.Addresses {
        existingAddresses[aws.StringValue(addr)] = true
    }

    var invalidAddresses []string
    for _, addr := range req.Addresses {
        if !existingAddresses[addr] {
            invalidAddresses = append(invalidAddresses, addr)
        }
    }

    if len(invalidAddresses) > 0 {
        utils.RespondWithError(w, http.StatusBadRequest, "The following IP addresses do not exist in the IP set: "+strings.Join(invalidAddresses, ", "))
        log.Println("The following IP addresses do not exist in the IP set:", strings.Join(invalidAddresses, ", "))
        return
    }

    var updatedAddresses []string
    for _, addr := range ipSet.IPSet.Addresses {
        if !contains(req.Addresses, aws.StringValue(addr)) {
            updatedAddresses = append(updatedAddresses, aws.StringValue(addr))
        }
    }

    params := &wafv2.UpdateIPSetInput{
        Id:       aws.String(req.Id),
        Name:     aws.String(req.Name),
        Scope:    aws.String("REGIONAL"),
        LockToken: aws.String(req.LockToken),
        Addresses: aws.StringSlice(updatedAddresses),
    }

    resp, err := config.WAFv2.UpdateIPSet(params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error updating IP set:", err)
        return
    }

    addressesJSON, _ := json.Marshal(req.Addresses)
    _, err = config.DB.Exec("INSERT INTO actions (name, type, action) VALUES (?, ?, ?)", req.Name, "remove", string(addressesJSON))
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error inserting action into database:", err)
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, resp)
    log.Println("Successfully removed IP addresses from IP set:", req.Name)
}

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
