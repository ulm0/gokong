package gokong

import (
	"encoding/json"
	"fmt"
)

type AclClient struct {
	config *Config
}

type AclRequest struct {
	Group *string `json:"group"`
}

type Acl struct {
	Id         *string `json:"id"`
	CreatedAt  *int    `json:"created_at"`
	Group      *string `json:"group"`
	ConsumerId *string `json:"consumer_id,omitempty"`
}

type Acls struct {
	Results []*Acl `json:"data,omitempty"`
	Total   int    `json:"total,omitempty"`
	Next    string `json:"next,omitempty"`
	Offset  string `json:"offset,omitempty"`
}

type AclFilter struct {
	Id         *string `json:"id"`
	ConsumerId *string `json:"consumer_id,omitempty"`
}

const AclsPath = "/acls/"

func (aclClient *AclClient) GetConsumerPerAcl(aclId string) (*Consumer, error) {
	r, body, errs := newGet(aclClient.config, aclClient.config.HostAddress+AclsPath+aclId+"/consumer").End()

	if errs != nil {
		return nil, fmt.Errorf("could not get acl consumer, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	consumer := &Consumer{}
	err := json.Unmarshal([]byte(body), consumer)
	if err != nil {
		return nil, fmt.Errorf("could not parse acl get response, error: %v", err)
	}

	if consumer.Id == "" {
		return nil, nil
	}

	return consumer, nil
}

func (aclClient *AclClient) GetAclsPerConsumer(consumerId string) (*Acls, error) {
	r, body, errs := newGet(aclClient.config, aclClient.config.HostAddress+ConsumersPath+consumerId+AclsPath).End()

	if errs != nil {
		return nil, fmt.Errorf("could not get acls for consumer, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	acls := &Acls{}
	err := json.Unmarshal([]byte(body), acls)
	if err != nil {
		return nil, fmt.Errorf("could not parse acl get response, error: %v", err)
	}

	// if acls. == nil {
	// 	return nil, nil
	// }

	return acls, nil
}

func (aclClient *AclClient) List() (*Acls, error) {
	return aclClient.ListFiltered(nil)
}

func (aclClient *AclClient) ListFiltered(filter *AclFilter) (*Acls, error) {
	address, err := addQueryString(aclClient.config.HostAddress+AclsPath, filter)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for acls filter, error: %v", err)
	}

	r, body, errs := newGet(aclClient.config, address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get apis, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	acls := &Acls{}
	err = json.Unmarshal([]byte(body), acls)
	if err != nil {
		return nil, fmt.Errorf("could not parse acls list response, error: %v", err)
	}

	return acls, nil
}

func (aclClient *AclClient) Create(consumerId string, newAcl *AclRequest) (*Acl, error) {

	r, body, errs := newPost(aclClient.config, aclClient.config.HostAddress+ConsumersPath+consumerId+AclsPath).Send(newAcl).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new acl, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	createdAcl := &Acl{}
	err := json.Unmarshal([]byte(body), createdAcl)
	if err != nil {
		return nil, fmt.Errorf("could not parse acl creation response, error: %v %s", err, body)
	}

	if createdAcl.Id == nil {
		return nil, fmt.Errorf("could not create acl, error: %v", body)
	}

	return createdAcl, nil
}

func (aclClient *AclClient) DeleteByName(name string) error {
	return aclClient.DeleteById(name)
}

// func (aclClient *AclClient) DeleteById(id string) error {

// 	r, body, errs := newDelete(aclClient.config, aclClient.config.HostAddress+ApisPath+id).End()
// 	if errs != nil {
// 		return fmt.Errorf("could not delete api, result: %v error: %v", r, errs)
// 	}

// 	if r.StatusCode == 401 || r.StatusCode == 403 {
// 		return fmt.Errorf("not authorised, message from kong: %s", body)
// 	}

// 	return nil
// }

// func (aclClient *AclClient) UpdateByName(name string, apiRequest *AclRequest) (*Acl, error) {
// 	return aclClient.UpdateById(name, apiRequest)
// }

// func (aclClient *AclClient) UpdateById(id string, apiRequest *AclRequest) (*Acl, error) {

// 	r, body, errs := newPatch(aclClient.config, aclClient.config.HostAddress+ApisPath+id).Send(apiRequest).End()
// 	if errs != nil {
// 		return nil, fmt.Errorf("could not update api, error: %v", errs)
// 	}

// 	if r.StatusCode == 401 || r.StatusCode == 403 {
// 		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
// 	}

// 	updatedApi := &Api{}
// 	err := json.Unmarshal([]byte(body), updatedApi)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not parse api update response, error: %v", err)
// 	}

// 	if updatedApi.Id == nil {
// 		return nil, fmt.Errorf("could not update api, error: %v", body)
// 	}

// 	return updatedApi, nil
// }

// func (a *AclRequest) MarshalJSON() ([]byte, error) {

// 	uris := a.Uris
// 	if uris == nil {
// 		uris = make([]*string, 0)
// 	}

// 	hosts := a.Hosts
// 	if hosts == nil {
// 		hosts = make([]*string, 0)
// 	}

// 	methods := a.Methods
// 	if methods == nil {
// 		methods = make([]*string, 0)
// 	}

// 	type Alias ApiRequest
// 	return json.Marshal(&struct {
// 		Uris    []*string `json:"uris"`
// 		Hosts   []*string `json:"hosts"`
// 		Methods []*string `json:"methods"`
// 		*Alias
// 	}{
// 		Uris:    uris,
// 		Hosts:   hosts,
// 		Methods: methods,
// 		Alias:   (*Alias)(a),
// 	})
// }

// func (a *Acl) UnmarshalJSON(data []byte) error {

// 	fixedJson := strings.Replace(string(data), `"hosts":{}`, `"hosts":[]`, -1)
// 	fixedJson = strings.Replace(fixedJson, `"uris":{}`, `"uris":[]`, -1)
// 	fixedJson = strings.Replace(fixedJson, `"methods":{}`, `"methods":[]`, -1)

// 	type Alias Api
// 	aux := &struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(a),
// 	}

// 	return json.Unmarshal([]byte(fixedJson), &aux)
// }
