package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/jitsucom/eventnative/resources"
	"strings"
)

type Token struct {
	Id           string   `mapstructure:"id" json:"id,omitempty"`
	ClientSecret string   `mapstructure:"client_secret" json:"client_secret,omitempty"`
	ServerSecret string   `mapstructure:"server_secret" json:"server_secret,omitempty"`
	Origins      []string `mapstructure:"origins" json:"origins,omitempty"`
}

type TokensPayload struct {
	Tokens []Token `json:"tokens,omitempty"`
}

type TokensHolder struct {
	//origins by client token
	clientTokensOrigins map[string][]string
	//origins by server token
	serverTokensOrigins map[string][]string

	//all token ids
	ids []string
	//token by: client_secret/server_secret/id
	all map[string]Token
}

func (th *TokensHolder) IsEmpty() bool {
	return th == nil || len(th.ids) == 0
}

//parse tokens from json bytes
func parseFromBytes(b []byte) (*TokensHolder, error) {
	payload := &TokensPayload{}
	err := json.Unmarshal(b, payload)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling tokens. Payload must be json with 'tokens' key: %v", err)
	}

	return reformat(payload.Tokens), nil
}

func fromStrings(clientSecrets, serverSecrets []string) *TokensHolder {
	var tokens []Token
	for _, clientSecret := range clientSecrets {
		tokens = append(tokens, Token{ClientSecret: clientSecret})
	}
	for _, serverSecret := range serverSecrets {
		tokens = append(tokens, Token{ServerSecret: serverSecret})
	}
	return reformat(tokens)
}

func reformat(tokens []Token) *TokensHolder {
	clientTokensOrigins := map[string][]string{}
	serverTokensOrigins := map[string][]string{}
	all := map[string]Token{}
	var ids []string

	for _, tokenObj := range tokens {
		if tokenObj.Id == "" {
			//hash from client,server secret will be id
			tokenObj.Id = resources.GetHash([]byte(tokenObj.ClientSecret + tokenObj.ServerSecret))
		}

		all[tokenObj.Id] = tokenObj
		ids = append(ids, tokenObj.Id)

		trimmedClientToken := strings.TrimSpace(tokenObj.ClientSecret)
		if trimmedClientToken != "" {
			clientTokensOrigins[trimmedClientToken] = tokenObj.Origins
			all[trimmedClientToken] = tokenObj
		}

		trimmedServerToken := strings.TrimSpace(tokenObj.ServerSecret)
		if trimmedServerToken != "" {
			serverTokensOrigins[trimmedServerToken] = tokenObj.Origins
			all[trimmedServerToken] = tokenObj
		}
	}

	return &TokensHolder{
		clientTokensOrigins: clientTokensOrigins,
		serverTokensOrigins: serverTokensOrigins,
		ids:                 ids,
		all:                 all,
	}
}
