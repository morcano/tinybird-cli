package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	tokensapi "tinybird-cli/api/tokens"
)

var scope string

var PutToken = &cobra.Command{
	Use:   "put",
	Short: "Put a token",
	Long: `Modifies an Auth token. More than one scope can be sent per request, all of them will be added as Auth token scopes. 
	Everytime an Auth token scope is modified, it overrides the existing one(s).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		adminTok := cmd.Flag("admin-token").Value.String()
		client := tokensapi.NewClient(adminTok)

		fmt.Println("Getting tokens...")
		tokens, err := client.Get()
		if err != nil {
			log.Fatalf("Failed to get tokens: %v", err)
		}

		bulkUpdate(client, tokens.FilterByName("AccountToken").Tokens)

		return nil
	},
}

// bulkUpdate sends a PUT request to the TinyBird API using the provided auth tokens.
// It adds the specified scopes to the request URL and sets the necessary headers.
// If the request fails or returns a non-200 status code, it returns an error.
func bulkUpdate(client *tokensapi.Client, tokens []tokensapi.Token) {
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
		go func(t tokensapi.Token) {
			defer wg.Done()
			payload := tokensapi.PutTokenParams{
				Token: t.Token,
			}
			// TODO: Remove the hardcoded sql_filter
			scopeList := strings.Split(scope, ",")
			for _, scope := range scopeList {
				payload.Scope = append(payload.Scope, scope+":accountId="+t.GetAccountIdFromName())
			}
			err := client.Put(payload)
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
}
