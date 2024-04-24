package tokens

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"tinybird-cli/utils"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var scope string
var decryptKey string

var PutToken = &cobra.Command{
	Use:   "put",
	Short: "Put a token",
	Long: `Modifies an Auth token. More than one scope can be sent per request, all of them will be added as Auth token scopes. 
	Everytime an Auth token scope is modified, it overrides the existing one(s).`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateFile(cmd.Flag("file").Value.String())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := cmd.Flag("file").Value.String()
		adminTok := cmd.Flag("admin-token").Value.String()
		tokens, err := utils.GetItems(filePath, decryptKey, utils.DecryptAES256ECB)
		if err != nil {
			return err
		}
		process(tokens, adminTok)
		return nil
	},
}

// sendPutRequest sends a PUT request to the TinyBird API using the provided user and admin tokens.
// It adds the specified scopes to the request URL and sets the necessary headers.
// If the request fails or returns a non-200 status code, it returns an error.
func sendPutRequest(userToken string, adminToken string) error {
	baseURL := fmt.Sprintf("https://api.tinybird.co/v0/tokens/%s", userToken)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	scopeList := strings.Split(scope, ",")
	query := parsedURL.Query()
	for _, scope := range scopeList {
		query.Add("scope", scope)
	}
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPut, parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status")
	}

	return nil
}

func process(tokens []string, adminToken string) {
	p := mpb.New(mpb.WithWidth(64))
	total := len(tokens)
	var successCount, failedCount int32

	bar := p.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name("Running... "),
			decor.Percentage(decor.WC{W: 5})),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 4}),
			decor.Any(func(s decor.Statistics) string {
				return fmt.Sprintf("Success: %d, Failed: %d", atomic.LoadInt32(&successCount), atomic.LoadInt32(&failedCount))
			}, decor.WC{W: 25}),
		),
	)

	var wg sync.WaitGroup
	wg.Add(total)

	for _, token := range tokens {
		go func(t string) {
			defer wg.Done()
			err := sendPutRequest(t, adminToken)
			if err != nil {
				atomic.AddInt32(&failedCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
			bar.Increment()
		}(token)
	}

	wg.Wait()
	p.Wait()
}

func init() {
	PutToken.Flags().StringVarP(&scope, "scope", "s", "", "New scope(s) will override existing ones")
	PutToken.Flags().StringVarP(&decryptKey, "decrypt-key", "d", "", "Key to decrypt tokens")
}
