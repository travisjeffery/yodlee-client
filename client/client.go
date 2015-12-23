package yodlee

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

// Client is the thing you use to talk to Yodlee's API.
type Client struct {
	Login        string
	Password     string
	SessionToken string
}

// New creates a `Client`.
func New(login, password string) *Client {
	return &Client{
		Login:    login,
		Password: password,
	}
}

// Authenticate authenticates your client with Yodlee.
func (c *Client) Authenticate() []error {
	token, errs := c.GetCobSessionToken()
	if errs != nil {
		return errs
	}
	c.SessionToken = token
	return nil
}

// GetCobSessionToken autenticates a cobrand.
func (c *Client) GetCobSessionToken() (string, []error) {
	var j struct {
		CobrandConversationCredentials struct {
			SessionToken string
		}
	}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/authenticate/coblogin", struct {
		CobrandLogin    string `json:"cobrandLogin"`
		CobrandPassword string `json:"cobrandPassword"`
	}{
		c.Login,
		c.Password,
	}, &j)
	if errs != nil {
		return "", errs
	}
	return j.CobrandConversationCredentials.SessionToken, nil
}

// GetUserSessionToken authenticates a user's login and password.
func (c *Client) GetUserSessionToken(login, password string) (string, []error) {
	var j struct {
		UserContext struct {
			ConversationCredentials struct {
				SessionToken string
			}
		}
	}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/authenticate/login", struct {
		Login           string `json:"login"`
		Password        string `json:"password"`
		CobSessionToken string `json:"cobSessionToken"`
	}{
		login,
		password,
		c.SessionToken,
	}, &j)
	if errs != nil {
		return "", errs
	}
	return j.UserContext.ConversationCredentials.SessionToken, nil
}

// GetAccountsOutput represents an account.
type GetAccountsOutput struct {
	Created                string `json:"created"`
	CredentialsChangedTime int    `json:"credentialsChangedTime"`
	Disabled               bool   `json:"disabled"`
	IsAgentError           bool   `json:"isAgentError"`
	IsCustom               bool   `json:"isCustom"`
	IsSiteError            bool   `json:"isSiteError"`
	IsUARError             bool   `json:"isUARError"`
	RetryCount             int    `json:"retryCount"`
	SiteAccountID          int    `json:"siteAccountId"`
	SiteInfo               struct {
		BaseURL             string `json:"baseUrl"`
		ContentServiceInfos []struct {
			ContainerInfo struct {
				AssetType     int    `json:"assetType"`
				ContainerName string `json:"containerName"`
			} `json:"containerInfo"`
			ContentServiceID int `json:"contentServiceId"`
			SiteID           int `json:"siteId"`
		} `json:"contentServiceInfos"`
		DefaultDisplayName    string `json:"defaultDisplayName"`
		DefaultOrgDisplayName string `json:"defaultOrgDisplayName"`
		EnabledContainers     []struct {
			AssetType     int    `json:"assetType"`
			ContainerName string `json:"containerName"`
		} `json:"enabledContainers"`
		HdLogoLastModified   int           `json:"hdLogoLastModified"`
		IsAlreadyAddedByUser bool          `json:"isAlreadyAddedByUser"`
		IsCustom             bool          `json:"isCustom"`
		IsHdLogoAvailable    bool          `json:"isHdLogoAvailable"`
		IsHeld               bool          `json:"isHeld"`
		IsOauthEnabled       bool          `json:"isOauthEnabled"`
		LoginForms           []interface{} `json:"loginForms"`
		LoginURL             string        `json:"loginUrl"`
		OrgID                int           `json:"orgId"`
		Popularity           int           `json:"popularity"`
		SiteID               int           `json:"siteId"`
		SiteSearchVisibility bool          `json:"siteSearchVisibility"`
	} `json:"siteInfo"`
	SiteRefreshInfo struct {
		Code               int  `json:"code"`
		IsMFAInputRequired bool `json:"isMFAInputRequired"`
		ItemRefreshInfo    []struct {
			ErrorCode         int `json:"errorCode"`
			ItemSuggestedFlow struct {
				SuggestedFlow   string `json:"suggestedFlow"`
				SuggestedFlowID int    `json:"suggestedFlowId"`
			} `json:"itemSuggestedFlow"`
			MemItemID  int `json:"memItemId"`
			RetryCount int `json:"retryCount"`
		} `json:"itemRefreshInfo"`
		NextUpdate      int `json:"nextUpdate"`
		NoOfRetry       int `json:"noOfRetry"`
		SiteRefreshMode struct {
			RefreshMode   string `json:"refreshMode"`
			RefreshModeID int    `json:"refreshModeId"`
		} `json:"siteRefreshMode"`
		SiteRefreshStatus struct {
			SiteRefreshStatus   string `json:"siteRefreshStatus"`
			SiteRefreshStatusID int    `json:"siteRefreshStatusId"`
		} `json:"siteRefreshStatus"`
		UpdateInitTime int `json:"updateInitTime"`
	} `json:"siteRefreshInfo"`
}

// GetAccounts gets the accounts for the given user session token.
func (c *Client) GetAccounts(token string) ([]*GetAccountsOutput, []error) {
	var output []*GetAccountsOutput
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/jsonsdk/SiteAccountManagement/getSiteAccounts", struct {
		CobSessiontoken  string `json:"cobSessionToken"`
		UserSessionToken string `json:"userSessionToken"`
	}{
		c.SessionToken,
		token,
	}, &output)
	return output, errs
}

// GetTransactionInput represents the arguments used to fetch transactions.
type GetTransactionInput struct {
	ContainerType    string `json:"transactionSearchRequest.containerType"`
	HigherFetchLimit string `json:"transactionSearchRequest.higherFetchLimit"`
	LowerFetchLimit  string `json:"transactionSearchRequest.lowerFetchLimit"`
	IgnoreUserInput  string `json:"transactionSearchRequest.ignoreUserInput"`
	EndNumber        int    `json:"transactionSearchRequest.resultRange.endNumber"`
	StartNumber      int    `json:"transactionSearchRequest.resultRange.startNumber"`
	CurrencyCode     string `json:"transactionSearchRequest.searchFilter.currencyCode"`
}

// NewGetTransactionInput creates a `GetTransactionInput` with defaults set.
func NewGetTransactionInput() *GetTransactionInput {
	return &GetTransactionInput{
		ContainerType:    "All",
		HigherFetchLimit: "500",
		LowerFetchLimit:  "1",
		IgnoreUserInput:  "true",
		EndNumber:        5,
		StartNumber:      1,
		CurrencyCode:     "USD",
	}
}

// GetTransactionsOutput represents the results of getting the transactions for a user.
type GetTransactionsOutput struct {
	CountOfAllTransaction      int `json:"countOfAllTransaction"`
	CountOfProjectedTxns       int `json:"countOfProjectedTxns"`
	CreditTotalOfProjectedTxns struct {
		Amount       int    `json:"amount"`
		CurrencyCode string `json:"currencyCode"`
	} `json:"creditTotalOfProjectedTxns"`
	CreditTotalOfTxns struct {
		Amount       float64 `json:"amount"`
		CurrencyCode string  `json:"currencyCode"`
	} `json:"creditTotalOfTxns"`
	DebitTotalOfProjectedTxns struct {
		Amount       int    `json:"amount"`
		CurrencyCode string `json:"currencyCode"`
	} `json:"debitTotalOfProjectedTxns"`
	DebitTotalOfTxns struct {
		Amount       float64 `json:"amount"`
		CurrencyCode string  `json:"currencyCode"`
	} `json:"debitTotalOfTxns"`
	NumberOfHits     int `json:"numberOfHits"`
	SearchIdentifier struct {
		Identifier string `json:"identifier"`
	} `json:"searchIdentifier"`
	SearchResult struct {
		Transactions []struct {
			AccessLevelRequired int `json:"accessLevelRequired"`
			Account             struct {
				AccountBalance struct {
					Amount       float64 `json:"amount"`
					CurrencyCode string  `json:"currencyCode"`
				} `json:"accountBalance"`
				AccountDisplayName struct {
					DefaultNormalAccountName string `json:"defaultNormalAccountName"`
				} `json:"accountDisplayName"`
				AccountName         string `json:"accountName"`
				AccountNumber       string `json:"accountNumber"`
				DecryptionStatus    bool   `json:"decryptionStatus"`
				IsAccountName       int    `json:"isAccountName"`
				ItemAccountID       int    `json:"itemAccountId"`
				ItemAccountStatusID int    `json:"itemAccountStatusId"`
				SiteName            string `json:"siteName"`
				SumInfoID           int    `json:"sumInfoId"`
			} `json:"account"`
			Amount struct {
				Amount       int    `json:"amount"`
				CurrencyCode string `json:"currencyCode"`
			} `json:"amount"`
			CategorisationSourceID int    `json:"categorisationSourceId"`
			CategorizationKeyword  string `json:"categorizationKeyword"`
			Category               struct {
				CategoryID            int    `json:"categoryId"`
				CategoryName          string `json:"categoryName"`
				CategoryTypeID        int    `json:"categoryTypeId"`
				IsBusiness            bool   `json:"isBusiness"`
				LocalizedCategoryName string `json:"localizedCategoryName"`
			} `json:"category"`
			CheckNumber         struct{} `json:"checkNumber"`
			ClassUpdationSource string   `json:"classUpdationSource"`
			Description         struct {
				Description          string `json:"description"`
				IsOlbUserDescription bool   `json:"isOlbUserDescription"`
				MerchantName         string `json:"merchantName"`
				SimpleDescription    string `json:"simpleDescription"`
				TransactionTypeDesc  string `json:"transactionTypeDesc"`
				ViewPref             bool   `json:"viewPref"`
			} `json:"description"`
			InvestmentTransactionView struct {
				HoldingType struct {
					HoldingTypeID int `json:"holdingTypeId"`
				} `json:"holdingType"`
				LotHandling struct {
					LotHandlingID int `json:"lotHandlingId"`
				} `json:"lotHandling"`
				NetCost int `json:"netCost"`
			} `json:"investmentTransactionView"`
			IsBusiness                   bool     `json:"isBusiness"`
			IsClosingTxn                 int      `json:"isClosingTxn"`
			IsMedical                    bool     `json:"isMedical"`
			IsPersonal                   bool     `json:"isPersonal"`
			IsReimbursable               bool     `json:"isReimbursable"`
			IsTaxable                    bool     `json:"isTaxable"`
			LocalizedTransactionBaseType string   `json:"localizedTransactionBaseType"`
			LocalizedTransactionType     string   `json:"localizedTransactionType"`
			Memo                         struct{} `json:"memo"`
			PostDate                     string   `json:"postDate"`
			Price                        struct {
				Amount       int    `json:"amount"`
				CurrencyCode string `json:"currencyCode"`
			} `json:"price"`
			RunningBalance int `json:"runningBalance"`
			Status         struct {
				Description          string `json:"description"`
				LocalizedDescription string `json:"localizedDescription"`
				StatusID             int    `json:"statusId"`
			} `json:"status"`
			TransactionBaseType         string `json:"transactionBaseType"`
			TransactionBaseTypeID       int    `json:"transactionBaseTypeId"`
			TransactionPostingOrder     int    `json:"transactionPostingOrder"`
			TransactionSearchResultType string `json:"transactionSearchResultType"`
			TransactionType             string `json:"transactionType"`
			TransactionTypeID           int    `json:"transactionTypeId"`
			ViewKey                     struct {
				ContainerType          string `json:"containerType"`
				IsParentMatch          bool   `json:"isParentMatch"`
				IsSystemGeneratedSplit bool   `json:"isSystemGeneratedSplit"`
				RowNumber              int    `json:"rowNumber"`
				TransactionCount       int    `json:"transactionCount"`
				TransactionID          int    `json:"transactionId"`
			} `json:"viewKey"`
		} `json:"transactions"`
	} `json:"searchResult"`
}

// GetTransactions gets transactions for the user and input.
func (c *Client) GetTransactions(token string, input *GetTransactionInput) (*GetTransactionsOutput, []error) {
	output := &GetTransactionsOutput{}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/jsonsdk/TransactionSearchService/executeUserSearchRequest", struct {
		*GetTransactionInput
		CobSessionToken  string `json:"cobSessionToken"`
		UserSessionToken string `json:"userSessionToken"`
	}{
		input,
		c.SessionToken,
		token,
	}, output)
	if errs != nil {
		return nil, errs
	}
	return output, nil
}

// request is a helper for making requests to Yodlee and formatting their responses.
func request(url string, content interface{}, data interface{}) []error {
	req := gorequest.New()
	_, body, errs := req.Post(url).
		Type("form").
		Send(content).
		End()
	if errs != nil {
		return errs
	}
	if err := json.Unmarshal([]byte(body), data); err != nil {
		return []error{err}
	}
	return nil
}
