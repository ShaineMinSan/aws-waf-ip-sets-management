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

type AddIPAddressRequest struct {
    Id       string   `json:"id"`
    Name     string   `json:"name"`
    Addresses []string `json:"addresses"`
    LockToken string  `json:"lockToken"`
}

func AddIPAddress(w http.ResponseWriter, r *http.Request) {
    var req AddIPAddressRequest
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

    var newAddresses []string
    for _, addr := range req.Addresses {
        if !existingAddresses[addr] {
            newAddresses = append(newAddresses, addr)
        }
    }

    if len(newAddresses) == 0 {
        utils.RespondWithError(w, http.StatusBadRequest, "All provided IP addresses already exist in the IP set.")
        log.Println("All provided IP addresses already exist in the IP set.")
        return
    }

    updatedAddresses := append(ipSet.IPSet.Addresses, aws.StringSlice(newAddresses)...)

    params := &wafv2.UpdateIPSetInput{
        Id:       aws.String(req.Id),
        Name:     aws.String(req.Name),
        Scope:    aws.String("REGIONAL"),
        LockToken: aws.String(req.LockToken),
        Addresses: updatedAddresses,
    }

    resp, err := config.WAFv2.UpdateIPSet(params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error updating IP set:", err)
        return
    }

    addressesJSON, _ := json.Marshal(newAddresses)
    _, err = config.DB.Exec("INSERT INTO actions (name, type, action) VALUES (?, ?, ?)", req.Name, "add", string(addressesJSON))
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        log.Println("Error inserting action into database:", err)
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, resp)
    log.Println("Successfully added IP addresses to IP set:", req.Name)
}
